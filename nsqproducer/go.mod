module github.com/Scalingo/go-utils/nsqproducer

go 1.16

require (
	github.com/Scalingo/go-utils/env v1.0.1
	github.com/Scalingo/go-utils/logger v1.1.0
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/golang/mock v1.6.0
	github.com/nsqio/go-nsq v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/sys v0.0.0-20211020174200-9d6173849985 // indirect
	gopkg.in/errgo.v1 v1.0.1
)

// In Dev you can uncomment the following line to use the local packages
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/env => ../env
