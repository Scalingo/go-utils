package graceful

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cloudflare/tableflip"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
)

type Service struct {
	httpServers []*http.Server
	graceful    *tableflip.Upgrader
	mx          sync.Mutex
	wg          *sync.WaitGroup
	// waitDuration is the duration which is waited for all connections to stop
	// in order to graceful shutdown the server. If some connections are still up
	// after this timer they'll be cut aggressively.
	waitDuration time.Duration
	// reloadWaitDuration is the duration the old process is waiting for
	// connection to close when a graceful restart has been ordered. The new
	// process is already working as expecting.
	reloadWaitDuration time.Duration
	// numServers is the number of servers to register before being ready
	numServers int
	// pidFile tracks the pid of the last child among the chain of graceful restart
	// Required for daemon manager to track the service
	pidFile string
}

type Option func(*Service)

func NewService(opts ...Option) *Service {
	s := &Service{
		httpServers:        make([]*http.Server, 0),
		wg:                 &sync.WaitGroup{},
		waitDuration:       time.Minute,
		reloadWaitDuration: 30 * time.Minute,
		numServers:         1,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithWaitDuration(d time.Duration) Option {
	return Option(func(s *Service) {
		s.waitDuration = d
	})
}

func WithReloadWaitDuration(d time.Duration) Option {
	return Option(func(s *Service) {
		s.reloadWaitDuration = d
	})
}

func WithPIDFile(path string) Option {
	return Option(func(s *Service) {
		s.pidFile = path
	})
}

func WithNumServers(n int) Option {
	return Option(func(s *Service) {
		s.numServers = n
	})
}

func (s *Service) getTableflipUpgrader(ctx context.Context) (*tableflip.Upgrader, error) {
	var err error
	s.mx.Lock()
	if s.graceful == nil {
		s.graceful, err = tableflip.New(tableflip.Options{
			UpgradeTimeout: s.reloadWaitDuration,
			PIDFile:        s.pidFile,
		})
		if err != nil {
			return nil, errors.Wrap(ctx, err, "creating tableflip upgrader")
		}
	}
	s.mx.Unlock()
	return s.graceful, nil
}

func (s *Service) ListenAndServeTLS(ctx context.Context, proto string, addr string, handler http.Handler, tlsConfig *tls.Config) error {
	httpServer := &http.Server{
		Addr:      addr,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}
	return s.listenAndServe(ctx, proto, addr, httpServer)
}

func (s *Service) ListenAndServe(ctx context.Context, proto string, addr string, handler http.Handler) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return s.listenAndServe(ctx, proto, addr, httpServer)
}

func (s *Service) listenAndServe(ctx context.Context, _ string, addr string, server *http.Server) error {
	log := logger.Get(ctx)

	s.mx.Lock()
	curServerCount := len(s.httpServers)
	s.httpServers = append(s.httpServers, server)
	if curServerCount == 0 {
		err := s.prepare(ctx)
		if err != nil {
			// purposefully do not wrap error here, as it is wrapped in prepare
			return err
		}
	}
	s.mx.Unlock()

	// Use tableflip to handle graceful restart requests
	upg, err := s.getTableflipUpgrader(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get upgrader")
	}

	// Listen must be called before Ready
	ln, err := upg.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(ctx, err, "upgrader listen")
	}

	if server.TLSConfig != nil {
		ln = tls.NewListener(ln, server.TLSConfig)
	}

	go func() {
		err := server.Serve(ln)
		if !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("http server serve")
		}
	}()

	if curServerCount+1 == s.numServers {
		err := s.finalize(ctx)
		if err != nil {
			// purposefully do not wrap error here, as it is wrapped in finalize
			return err
		}
	}

	return nil
}

// prepare is called before the first server is started.
func (s *Service) prepare(ctx context.Context) error {
	if s.pidFile != "" {
		err := os.Remove(s.pidFile)
		if err != nil && !os.IsNotExist(err) {
			return errors.Wrap(ctx, err, "fail to remove PID file")
		}
	}

	// setup the signal handling
	go s.setupSignals(ctx)

	return nil
}

// finalize is called when all servers are started.
func (s *Service) finalize(ctx context.Context) error {
	log := logger.Get(ctx)

	upg, err := s.getTableflipUpgrader(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get upgrader")
	}
	defer upg.Stop()

	log.Info("ready")
	if err := upg.Ready(); err != nil {
		return errors.Wrapf(ctx, err, "upgrader notify ready")
	}
	<-upg.Exit()
	log.Info("upgrader finished")

	// Normally the server should be always gracefully stopped and entering the
	// above condition when server is closed If by any mean the serve stops
	// without error, we're stopping the server ourselves here.  This code is a
	// security to free resource but should be unreachable
	ctx, cancel := context.WithTimeout(ctx, s.waitDuration)
	defer cancel()
	err = s.shutdown(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to shutdown service")
	}

	// Wait for connections to drain.
	s.mx.Lock()
	errChan := make(chan error, len(s.httpServers))
	for i, httpServer := range s.httpServers {
		err = httpServer.Shutdown(ctx)
		if err != nil {
			errChan <- errors.Wrapf(ctx, err, "server shutdown %d", i)
		}
	}
	s.mx.Unlock()
	close(errChan)
	var shutdownErr error
	for err := range errChan {
		if shutdownErr == nil {
			shutdownErr = err
		} else {
			shutdownErr = errors.Wrap(ctx, shutdownErr, err.Error())
		}
	}
	if shutdownErr != nil {
		return shutdownErr
	}

	return nil
}

// IncConnCount has to be used when connections are hijacked because in
// this case http.Server doesn't track these connection anymore, but you
// may not want to cut them abrutely.
func (s *Service) IncConnCount(ctx context.Context) {
	log := logger.Get(ctx)
	log.Debug("inc conn count")
	s.wg.Add(1)
}

// DecConnCount is the same as IncConnCount, but you need to call it when
// the hijacked connection is stopped
func (s *Service) DecConnCount(ctx context.Context) {
	log := logger.Get(ctx)
	log.Debug("dec conn count")
	s.wg.Done()
}

// shutdown stops the HTTP listener and then wait for any active hijacked
// connection to stop http.Server#Shutdown is graceful but the documentation
// specifies hijacked connections and websockets have to be handled by the
// developer.
func (s *Service) shutdown(ctx context.Context) error {
	log := logger.Get(ctx)

	errChan := make(chan error, len(s.httpServers))
	var wg sync.WaitGroup

	for i, httpServer := range s.httpServers {
		wg.Add(1)
		go func(i int, httpServer *http.Server) {
			defer wg.Done()
			log := logger.Get(ctx)
			if len(s.httpServers) > 1 {
				log = log.WithField("index", i)
			}
			log.Info("shutting down http server")
			err := httpServer.Shutdown(ctx)
			if err != nil {
				log.WithError(err).Error("fail to shutdown http server")
				errChan <- errors.Wrapf(ctx, err, "fail to shutdown http server %d", i)
			} else {
				log.Info("http server is stopped")
			}
		}(i, httpServer)
	}

	wg.Wait()
	close(errChan)

	var shutdownErr error
	for err := range errChan {
		if shutdownErr == nil {
			shutdownErr = err
		} else {
			shutdownErr = errors.Wrap(ctx, shutdownErr, err.Error())
		}
	}

	if shutdownErr != nil {
		return shutdownErr
	}

	log.Info("wait hijacked connections")
	err := s.waitHijackedConnections(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to wait hijacked connections")
	}
	log.Info("no more connection running")

	return nil
}

func (s *Service) waitHijackedConnections(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
