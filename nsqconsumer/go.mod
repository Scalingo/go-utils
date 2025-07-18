module github.com/Scalingo/go-utils/nsqconsumer

go 1.24

require (
	github.com/Scalingo/go-utils/errors/v2 v2.5.1
	github.com/Scalingo/go-utils/logger v1.9.1
	github.com/Scalingo/go-utils/nsqproducer v1.3.1
	github.com/nsqio/go-nsq v1.1.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stvp/rollbar v0.5.1
)

require (
	github.com/Scalingo/go-utils/env v1.2.1 // indirect
	github.com/gofrs/uuid/v5 v5.3.2 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.34.0 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
)

// In Dev you can uncomment the following line to use the local packages
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/nsqproducer => ../nsqproducer
// replace github.com/Scalingo/go-utils/env => ../env
