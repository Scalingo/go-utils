module github.com/Scalingo/go-utils/gomock_generator

go 1.14

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger

require (
	github.com/Scalingo/go-utils/logger v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/urfave/cli v1.22.4
	golang.org/x/sys v0.0.0-20201112073958-5cba982894dd // indirect
)
