package graceful

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cloudflare/tableflip"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-utils/logger"
)

type Service struct {
	httpServer *http.Server
	graceful   *tableflip.Upgrader
	wg         *sync.WaitGroup
	stopped    chan error
	// waitDuration is the duration which is waited for all connections to stop
	// in order to graceful shutdown the server. If some connections are still up
	// after this timer they'll be cut aggressively.
	waitDuration time.Duration
	// reloadWaitDuration is the duration the old process is waiting for
	// connection to close when a graceful restart has been ordered. The new
	// process is already working as expecting.
	reloadWaitDuration time.Duration
	// pidFile tracks the pid of the last child among the chain of graceful restart
	// Required for daemon manager to track the service
	pidFile string
}

type Option func(*Service)

func NewService(opts ...Option) *Service {
	s := &Service{
		wg:                 &sync.WaitGroup{},
		stopped:            make(chan error),
		waitDuration:       time.Minute,
		reloadWaitDuration: 30 * time.Minute,
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

func (s *Service) listenAndServe(ctx context.Context, proto string, addr string, server *http.Server) error {
	log := logger.Get(ctx)

	if s.pidFile != "" {
		pid := os.Getpid()
		err := os.WriteFile(s.pidFile, []byte(fmt.Sprintf("%d\n", pid)), 0600)
		if err != nil {
			return errgo.Notef(err, "fail to write PID file")
		}
	}

	// Use tableflip to handle graceful restart requests
	upg, err := tableflip.New(tableflip.Options{
		UpgradeTimeout: s.reloadWaitDuration,
		PIDFile:        s.pidFile,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer upg.Stop()
	s.graceful = upg

	// setup the signal handling
	go s.setupSignals(ctx)

	// Listen must be called before Ready
	ln, err := upg.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Can't listen:", err)
	}

	if server.TLSConfig != nil {
		ln = tls.NewListener(ln, server.TLSConfig)
	}

	s.httpServer = server

	go func() {
		err := server.Serve(ln)
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println("HTTP server:", err)
		}
	}()

	log.Printf("ready")
	if err := upg.Ready(); err != nil {
		panic(err)
	}
	<-upg.Exit()

	// Normally the server should be always gracefully stopped and entering the
	// above condition when server is closed If by any mean the serve stops
	// without error, we're stopping the server ourselves here.  This code is a
	// security to free resource but should be unreachable
	ctx, cancel := context.WithTimeout(ctx, s.waitDuration)
	defer cancel()
	log.Println("shutting down")
	err = s.shutdown(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to shutdown server")
	}

	// Wait for connections to drain.
	err = server.Shutdown(context.Background())
	if err != nil {
		log.Println("server shutdown:", err)
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
	log.Info("shutting down http server")
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to shutdown http server")
	}
	log.Info("http server is stopped")

	log.Info("wait hijacked connections")
	err = s.waitHijackedConnections(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to wait hijacked connections")
	}
	log.Info("no more connection running")

	return nil
}

func (s *Service) waitStopped() error {
	return <-s.stopped
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
