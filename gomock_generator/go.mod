module github.com/Scalingo/go-utils/gomock_generator

go 1.16

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger

require (
	github.com/Scalingo/go-utils/logger v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli v1.22.9
	golang.org/x/sys v0.0.0-20211020174200-9d6173849985 // indirect
)
