package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	etcdcron "github.com/go-utils/cron"
)

func main() {
	log.Println("starting")
	cron, err := etcdcron.New()
	if err != nil {
		log.Fatal("fail to create etcd-cron", err)
	}
	err = cron.AddJob(etcdcron.Job{
		Name:   "test",
		Rhythm: "*/4 * * * * *",
		Func: func(_ context.Context) error {
			// Use default logging of etcd-cron
			return errors.New("Horrible Error")
		},
	})
	if err != nil {
		log.Fatal("Fail to add the cron job", err)
	}

	err = cron.AddJob(etcdcron.Job{
		Name:   "test-v2",
		Rhythm: "*/10 * * * * *",
		Func: func(_ context.Context) error {
			log.Println("Every 10 seconds from", os.Getpid())
			return nil
		},
	})
	if err != nil {
		log.Fatal("Fail to add the cron job", err)
	}

	cron.Start(context.Background())
	time.Sleep(100 * time.Second)
	cron.Stop()
}
