module github.com/Scalingo/go-utils/nsqconsumer

go 1.20

require (
	github.com/Scalingo/go-utils/errors/v2 v2.2.0
	github.com/Scalingo/go-utils/logger v1.2.0
	github.com/Scalingo/go-utils/nsqproducer v1.1.2
	github.com/nsqio/go-nsq v1.1.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stvp/rollbar v0.5.1
	gopkg.in/errgo.v1 v1.0.1
)

require (
	github.com/Scalingo/go-utils/env v1.1.1 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.7.0 // indirect
)

// In Dev you can uncomment the following line to use the local packages
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/nsqproducer => ../nsqproducer
// replace github.com/Scalingo/go-utils/env => ../env
