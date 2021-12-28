module github.com/Scalingo/go-utils/nsqlbproducer

go 1.16

require (
	github.com/Scalingo/go-utils/env v1.1.0
	github.com/Scalingo/go-utils/nsqproducer v1.1.1
	github.com/golang/mock v1.6.0
	github.com/golang/snappy v0.0.2 // indirect
	github.com/nsqio/go-nsq v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/errgo.v1 v1.0.1
)

// In Dev you can uncomment the following line to use the local packages
// replace github.com/Scalingo/go-utils/env => ../env
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/nsqproducer => ../nsqproducer
