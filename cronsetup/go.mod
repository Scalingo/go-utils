module github.com/Scalingo/go-utils/cronsetup

go 1.14

require (
	github.com/Scalingo/go-etcd-cron v1.2.1
	github.com/Scalingo/go-utils/logger v0.0.0-00010101000000-000000000000
	go.etcd.io/etcd/v3 v3.3.0-rc.0.0.20200826232710-c20cc05fc548
)

replace github.com/Scalingo/go-utils/logger => ../logger
