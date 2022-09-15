module github.com/Scalingo/go-utils/nsqconsumer

go 1.17

require (
	github.com/Scalingo/go-utils/logger v1.2.0
	github.com/Scalingo/go-utils/nsqproducer v1.1.2
	github.com/nsqio/go-nsq v1.1.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stvp/rollbar v0.5.1
	gopkg.in/errgo.v1 v1.0.1
)

require (
	github.com/Scalingo/go-utils/env v1.1.0 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/golang/snappy v0.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

// In Dev you can uncomment the following line to use the local packages
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/nsqproducer => ../nsqproducer
// replace github.com/Scalingo/go-utils/env => ../env
