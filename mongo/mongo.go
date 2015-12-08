package mongo

import (
	"errors"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	DefaultDatabaseName string
	sessionOnce         = sync.Once{}
	_session            *mgo.Session
)

func Session() *mgo.Session {
	sessionOnce.Do(func() {
		err := errors.New("")
		for err != nil {
			_session, err = buildSession()
			if err != nil {
				log.Println("init mongo: fail to create session", err, "wait 10sec")
				time.Sleep(10 * time.Second)
			}
		}
	})
	return _session
}

func buildSession() (*mgo.Session, error) {
	rawURL := os.Getenv("MONGO_URL")
	if rawURL == "" {
		rawURL = "mongodb://localhost:27017/" + DefaultDatabaseName
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.New("not a valid MONGO_URL")
	}
	log.Println("init mongo on", u.Host)
	s, err := mgo.Dial(rawURL)
	if err != nil {
		return nil, err
	}
	return s, nil
}
