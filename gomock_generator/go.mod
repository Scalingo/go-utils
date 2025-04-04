module github.com/Scalingo/go-utils/gomock_generator

go 1.23.3
toolchain go1.24.1

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger

require (
	github.com/Scalingo/go-utils/logger v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/urfave/cli v1.22.16
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)
