module github.com/Scalingo/go-utils/gomock_generator

go 1.14

replace github.com/Scalingo/go-utils/logger => ../logger

require (
	github.com/Scalingo/go-utils/logger v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/urfave/cli v1.22.4
)
