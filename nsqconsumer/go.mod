module github.com/Scalingo/go-utils/nsqconsumer

go 1.14

require (
	github.com/Scalingo/go-utils/logger v0.0.0-00010101000000-000000000000
	github.com/Scalingo/go-utils/nsqproducer v0.0.0-00010101000000-000000000000
	github.com/nsqio/go-nsq v1.0.8
	github.com/sirupsen/logrus v1.7.0
	github.com/stvp/rollbar v0.5.1
	gopkg.in/errgo.v1 v1.0.1
)

replace github.com/Scalingo/go-utils/logger => ../logger

replace github.com/Scalingo/go-utils/nsqproducer => ../nsqproducer

replace github.com/Scalingo/go-utils/env => ../env
