module github.com/Scalingo/go-utils/cronsetup

go 1.14

require (
	github.com/Scalingo/go-etcd-cron v1.2.1
	github.com/Scalingo/go-utils/logger v1.0.0
	go.etcd.io/etcd/v3 v3.3.0-rc.0.0.20200826232710-c20cc05fc548
	golang.org/x/sys v0.0.0-20201112073958-5cba982894dd // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
