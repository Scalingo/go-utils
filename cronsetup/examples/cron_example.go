package main

import (
	"context"
	"errors"
	"os"
	"time"

	"go.etcd.io/etcd/clientv3"

	"github.com/Scalingo/go-utils/cronsetup"
	"github.com/Scalingo/go-utils/logger"
)

func main() {
	log := logger.Default()
	ctx := logger.ToCtx(context.Background(), log)

	log.Info("Starting cronsetup example")

	cancel, err := cronsetup.Setup(ctx, cronsetup.SetupOpts{
		EtcdConfig: func() (clientv3.Config, error) {
			return *cfg.EtcdConfig, nil
		},
		Jobs: []cronsetup.Job{
			cronsetup.Job{
				Name:   "test",
				Rhythm: "*/4 * * * * *",
				Func: func(_ context.Context) error {
					// Use default logging of etcd-cron
					return errors.New("horrible error")
				},
			},
			cronsetup.Job{
				Name:   "test-v2",
				Rhythm: "*/10 * * * * *",
				Func: func(_ context.Context) error {
					log.Println("Every 10 seconds from", os.Getpid())
					return nil
				},
			},
		},
	})
	if err != nil {
		return panic(err)
	}
	defer cancel()

	time.Sleep(100 * time.Second)
}
