module github.com/Scalingo/go-utils/cronsetup

go 1.16

require (
	github.com/Scalingo/go-etcd-cron v1.2.2-0.20210728122616-b340a9263e36
	github.com/Scalingo/go-utils/logger v1.0.0
	go.etcd.io/etcd/client/v3 v3.5.0
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
