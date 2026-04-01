module github.com/Scalingo/go-utils/gomock_generator

go 1.25.0

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger

require (
	github.com/Scalingo/go-utils/errors/v3 v3.2.0
	github.com/Scalingo/go-utils/logger v1.12.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.4
	github.com/urfave/cli/v3 v3.8.0
)

require (
	golang.org/x/sys v0.42.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)
