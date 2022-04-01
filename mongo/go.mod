module github.com/Scalingo/go-utils/mongo

go 1.16

require (
	github.com/Scalingo/go-utils/errors v1.1.0
	github.com/Scalingo/go-utils/logger v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	golang.org/x/sys v0.0.0-20211020174200-9d6173849985 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

// Uncomment if you want to use the local version of these packages (for development purpose)
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/errors => ../errors
