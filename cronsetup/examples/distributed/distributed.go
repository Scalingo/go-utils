package main

import (
	"context"
	"errors"
	"os"
	"time"

	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/cronsetup"
	"github.com/Scalingo/go-utils/logger"
)

const (
	defaultEtcdEndpoint = "127.0.0.1:2379"
)

func main() {
	log := logger.Default()
	ctx := logger.ToCtx(context.Background(), log)

	log.Info("Starting cronsetup distributed mode example")

	cancel, err := cronsetup.Setup(ctx, cronsetup.SetupOpts{
		EtcdConfig: func() (etcdv3.Config, error) {
			return etcdv3.Config{
				Endpoints: []string{defaultEtcdEndpoint},
			}, nil
		},
		Jobs: []cronsetup.Job{
			{
				Name:   "test",
				Rhythm: "*/4 * * * * *",
				Func: func(_ context.Context) error {
					// Use default logging of cronsetup
					return errors.New("horrible error in cron job \"test\"")
				},
			},
			{
				Name:   "test-v2",
				Rhythm: "*/10 * * * * *",
				Func: func(ctx context.Context) error {
					log := logger.Get(ctx)
					log.Info("[test-v2] Every 10 seconds from ", os.Getpid())
					return nil
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	defer cancel()

	// Stop the example in 100 seconds
	time.Sleep(100 * time.Second)
}
