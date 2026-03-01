module github.com/Scalingo/go-utils/graceful

go 1.25

require (
	github.com/Scalingo/go-utils/errors/v3 v3.2.0
	github.com/Scalingo/go-utils/logger v1.12.0
	github.com/cloudflare/tableflip v1.2.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	golang.org/x/sys v0.41.0 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
