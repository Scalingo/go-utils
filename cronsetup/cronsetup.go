package cronsetup

import (
	"context"

	"github.com/sirupsen/logrus"
	etcdclient "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/errors/v3"

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
		return nil, errors.Wrap(ctx, err, "get etcdv3 config")
	}

	etcdMutexBuilder, err := etcdcron.NewEtcdMutexBuilder(etcdConfig)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "get etcd mutex builder")
	}

	funcCtx := func(ctx context.Context, j etcdcron.Job) context.Context {
		log := logger.Get(ctx)
		requestID, ok := ctx.Value("request_id").(string)
		if !ok {
			requestUUID, err := uuid.NewV4()
			if err != nil {
				log.WithError(err).Error("Error generating UUID v4")
			} else {
				requestID = requestUUID.String()
				ctx = context.WithValue(ctx, "request_id", requestID) // nolint:revive
			}
		}
		ctx, _ = logger.WithFieldsToCtx(ctx, logrus.Fields{
			"cron-job":   j.Name,
			"request_id": requestID,
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
		return nil, errors.Wrap(ctx, err, "create etcd cron")
	}

	for _, job := range opts.Jobs {
		err := c.AddJob(job)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "add the cron job")
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
