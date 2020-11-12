module github.com/Scalingo/go-utils/graceful

go 1.14

require (
	github.com/Scalingo/go-utils/logger v0.0.0-00010101000000-000000000000
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/stretchr/testify v1.2.2
	gopkg.in/errgo.v1 v1.0.1
)

replace github.com/Scalingo/go-utils/logger => ../logger
