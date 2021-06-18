module github.com/Scalingo/go-utils/cronsetup

go 1.16

require (
	github.com/Scalingo/go-etcd-cron v1.2.1
	github.com/Scalingo/go-utils/logger v1.0.0
	github.com/prometheus/common v0.15.0 // indirect
	go.etcd.io/etcd/v3 v3.3.0-rc.0.0.20200826232710-c20cc05fc548
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/sys v0.0.0-20201112073958-5cba982894dd // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20201111145450-ac7456db90a6 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
