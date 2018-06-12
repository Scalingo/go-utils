package nsqlbproducer

import (
	"strings"

	"github.com/Scalingo/go-utils/env"
	"github.com/Scalingo/go-utils/nsqproducer"
)

func FromEnv() (*NsqLBProducer, error) {
	E := env.InitMapFromEnv(map[string]string{
		"NSQD_TLS":        "false",
		"NSQD_TLS_CACERT": "",
		"NSQD_TLS_CERT":   "",
		"NSQD_TLS_KEY":    "",

		"NSQLOOKUPD_URLS": "localhost:4161",

		"NSQD_HOSTS": "localhost:4150",

		"NSQ_MAX_IN_FLIGHT": "20",
	})

	var hosts []Host
	for _, host := range strings.Split(E["NSQD_HOSTS"], ",") {
		splittedHost := strings.Split(host, ":")
		hosts = append(hosts, Host{
			Host: splittedHost[0],
			Port: splittedHost[1],
		})
	}

	nsqConfig, err := nsqproducer.NsqConfigFromEnv(E)
	if err != nil {
		return nil, err
	}

	return New(LBProducerOpts{
		Hosts:     hosts,
		NsqConfig: nsqConfig,
	})
}
