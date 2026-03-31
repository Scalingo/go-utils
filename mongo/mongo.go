package mongo

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/logger"
)

var (
	DefaultDatabaseName string
	sessionOnce         = sync.Once{}
	_session            *mgo.Session
)

// Session tries to open a session against a MongoDB database in a loop.
// The argument is NOT a context but a logger for the following reason. We create an initial connection and clone it. It means that only the initial session will have a context. Once this context expires the initial session is closed, and cloning it is not possible anymore, leading to errors in all cloned sessions.
// The linter `contextcheck` may complain that this function does not take a context in argument. In such case, add the following comment before the line calling this function: "//nolint: contextcheck"
func Session(log logrus.FieldLogger) *mgo.Session {
	sessionOnce.Do(func() {
		log := log.WithField("process", "mongo-init")
		ctx := logger.ToCtx(context.Background(), log)

		err := errors.New(ctx, "")
		for err != nil {
			_session, err = BuildSession(ctx, os.Getenv("MONGO_URL"))
			if err != nil {
				retryDelay := 10 * time.Second
				log.WithError(err).Errorf("Failed to create a MongoDB session, retry in %v", retryDelay)
				time.Sleep(retryDelay)
			}
		}
	})
	return _session
}

func BuildSession(ctx context.Context, rawURL string) (*mgo.Session, error) {
	log := logger.Get(ctx)
	if rawURL == "" {
		rawURL = "mongodb://localhost:27017/" + DefaultDatabaseName
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "not a valid MONGO_URL")
	}

	withTLS := false
	if u.Query().Get("ssl") == "true" {
		withTLS = true
		rawURL = strings.Replace(rawURL, "?ssl=true", "?", 1)
		rawURL = strings.Replace(rawURL, "&ssl=true", "", 1)
	}

	timeout := 10 * time.Second
	queryTimeout := u.Query().Get("timeout")
	if queryTimeout != "" {
		timeout, err = time.ParseDuration(queryTimeout)
		if err != nil {
			return nil, errors.New(ctx, "invalid duration in timeout parameter")
		}
		rawURL = strings.Replace(rawURL, "?timeout="+queryTimeout, "?", 1)
		rawURL = strings.Replace(rawURL, "&timeout="+queryTimeout, "", 1)
	}

	info, err := mgo.ParseURL(rawURL)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "parse MongoDB connection URL")
	}
	info.Timeout = timeout
	if withTLS {
		tlsConfig := &tls.Config{}
		tlsConfig.InsecureSkipVerify = true
		info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	log.WithField("mongodb_host", u.Host).Info("Initialize the MongoDB connection")
	s, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "connect to MongoDB")
	}
	return s, nil
}
