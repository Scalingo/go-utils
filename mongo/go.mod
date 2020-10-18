module github.com/Scalingo/go-utils/mongo

go 1.14

require (
	github.com/Scalingo/go-utils/errors v0.0.0-00010101000000-000000000000
	github.com/Scalingo/go-utils/logger v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.2.2
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

replace github.com/Scalingo/go-utils/logger => ../logger

replace github.com/Scalingo/go-utils/errors => ../errors
