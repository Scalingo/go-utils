package cronsetup

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	etcdclient "go.etcd.io/etcd/client/v3"

	etcdcron "github.com/Scalingo/go-etcd-cron"
	"github.com/Scalingo/go-utils/logger"

	"github.com/gofrs/uuid/v5"
)

type SetupOpts struct {
	EtcdConfig func() (etcdclient.Config, error)
	Jobs       []etcdcron.Job
}

func Setup(ctx context.Context, opts SetupOpts) (func(), error) {
	etcdConfig, err := opts.EtcdConfig()
	if err != nil {
		return nil, fmt.Errorf("fail to get etcd v3 config: %v", err)
	}

	etcdMutexBuilder, err := etcdcron.NewEtcdMutexBuilder(etcdConfig)
	if err != nil {
		return nil, fmt.Errorf("fail to get etcd mutex builder: %v", err)
	}

	funcCtx := func(ctx context.Context, j etcdcron.Job) context.Context {
		log := logger.Get(ctx)
		requestID, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("generate UUID v4")
		} else {
			ctx = context.WithValue(ctx, "request_id", requestID.String())
		}
		ctx, _ = logger.WithFieldsToCtx(ctx, logrus.Fields{
			"cron-job":   j.Name,
			"request_id": requestID.String(),
		})
		return ctx
	}

	errorHandler := func(ctx context.Context, j etcdcron.Job, err error) {
		log := logger.Get(ctx)
		log.WithError(err).Error("Error when running cron job")
	}

	c, err := etcdcron.New(
		etcdcron.WithEtcdErrorsHandler(errorHandler),
		etcdcron.WithErrorsHandler(errorHandler),
		etcdcron.WithEtcdMutexBuilder(etcdMutexBuilder),
		etcdcron.WithFuncCtx(funcCtx),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to create etcd cron: %v", err)
	}

	for _, job := range opts.Jobs {
		err := c.AddJob(job)
		if err != nil {
			return nil, errors.Wrap(err, "fail to add the cron job")
		}
	}

	log := logger.Get(ctx)
	log.Info("Starting etcd-cron")

	c.Start(ctx)
	return func() {
		log.Info("Stopping etcd-cron")
		c.Stop()
	}, nil
}
