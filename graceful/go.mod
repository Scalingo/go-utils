module github.com/Scalingo/go-utils/graceful

go 1.14

require (
	github.com/Scalingo/go-utils/logger v1.0.0 // indirect
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/stretchr/testify v1.2.2
	golang.org/x/sys v0.0.0-20201112073958-5cba982894dd // indirect
	gopkg.in/errgo.v1 v1.0.1
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
