module github.com/Scalingo/go-utils/nsqproducer

go 1.24.0

require (
	github.com/Scalingo/go-utils/env v1.2.1
	github.com/Scalingo/go-utils/logger v1.9.1
	github.com/gofrs/uuid/v5 v5.3.2
	github.com/golang/mock v1.6.0
	github.com/nsqio/go-nsq v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	gopkg.in/errgo.v1 v1.0.1
)

require (
	github.com/golang/snappy v1.0.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

// In Dev you can uncomment the following line to use the local packages
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/env => ../env
