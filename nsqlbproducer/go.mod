module github.com/Scalingo/go-utils/nsqlbproducer

go 1.14

replace github.com/Scalingo/go-utils/env => ../env

replace github.com/Scalingo/go-utils/logger => ../logger

replace github.com/Scalingo/go-utils/nsqproducer => ../nsqproducer

require (
	github.com/Scalingo/go-utils/env v0.0.0-00010101000000-000000000000
	github.com/Scalingo/go-utils/nsqproducer v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.4.4
	github.com/nsqio/go-nsq v1.0.8
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	gopkg.in/errgo.v1 v1.0.1
)
