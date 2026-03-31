package nsqproducer

import (
	"context"
	"strconv"

	nsq "github.com/nsqio/go-nsq"

	"github.com/Scalingo/go-utils/env"
	"github.com/Scalingo/go-utils/errors/v3"
)

func FromEnv(ctx context.Context) (*NsqProducer, error) {
	envMap := env.InitMapFromEnv(map[string]string{
		"NSQD_TLS":        "false",
		"NSQD_TLS_CACERT": "",
		"NSQD_TLS_CERT":   "",
		"NSQD_TLS_KEY":    "",

		"NSQLOOKUPD_URLS": "localhost:4161",

		"NSQD_HOST": "localhost",
		"NSQD_PORT": "4150",

		"NSQ_MAX_IN_FLIGHT": "20",
	})

	nsqConfig, err := NsqConfigFromEnv(ctx, envMap)
	if err != nil {
		return nil, err
	}

	return New(ProducerOpts{
		Host: envMap["NSQD_HOST"],
		Port: envMap["NSQD_PORT"],

		NsqConfig: nsqConfig,
	})
}

func NsqConfigFromEnv(ctx context.Context, envMap map[string]string) (*nsq.Config, error) {
	nsqConfig := nsq.NewConfig()
	if envMap["NSQD_TLS"] == "true" {
		err := nsqConfig.Set("tls_v1", true)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "set tls_v1")
		}
		err = nsqConfig.Set("tls_root_ca_file", envMap["NSQD_TLS_CACERT"])
		if err != nil {
			return nil, errors.Wrap(ctx, err, "set tls_root_ca_file")
		}
		err = nsqConfig.Set("tls_cert", envMap["NSQD_TLS_CERT"])
		if err != nil {
			return nil, errors.Wrap(ctx, err, "set tls_cert")
		}
		err = nsqConfig.Set("tls_key", envMap["NSQD_TLS_KEY"])
		if err != nil {
			return nil, errors.Wrap(ctx, err, "set tls_key")
		}
	}

	maxInFlight, err := strconv.Atoi(envMap["NSQ_MAX_IN_FLIGHT"])
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "invalid max in flight: %s", envMap["NSQ_MAX_IN_FLIGHT"])
	}

	err = nsqConfig.Set("max_in_flight", maxInFlight)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "set max_in_flight")
	}

	return nsqConfig, nil
}
