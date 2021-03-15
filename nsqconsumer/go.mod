module github.com/Scalingo/go-utils/nsqconsumer

go 1.14

require (
	github.com/Scalingo/go-utils/logger v1.0.0
	github.com/Scalingo/go-utils/nsqproducer v1.0.0
	github.com/golang/snappy v0.0.2 // indirect
	github.com/nsqio/go-nsq v1.0.8
	github.com/sirupsen/logrus v1.8.1
	github.com/stvp/rollbar v0.5.1
	golang.org/x/sys v0.0.0-20201112073958-5cba982894dd // indirect
	gopkg.in/errgo.v1 v1.0.1
)

// In Dev you can uncomment the following line to use the local packages
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/nsqproducer => ../nsqproducer
// replace github.com/Scalingo/go-utils/env => ../env
