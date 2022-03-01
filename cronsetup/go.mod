module github.com/Scalingo/go-utils/cronsetup

go 1.16

require (
	github.com/Scalingo/go-etcd-cron v1.3.0
	github.com/Scalingo/go-utils/logger v1.1.0
	go.etcd.io/etcd/client/v3 v3.5.2
	golang.org/x/sys v0.0.0-20211020174200-9d6173849985 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
