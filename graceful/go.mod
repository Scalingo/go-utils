module github.com/Scalingo/go-utils/graceful

go 1.20

require (
	github.com/Scalingo/go-utils/logger v1.2.0
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/stretchr/testify v1.8.4
	gopkg.in/errgo.v1 v1.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/freeport v0.0.0-20150612182905-d4adf43b75b9 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
