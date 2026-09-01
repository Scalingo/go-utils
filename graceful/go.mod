module github.com/Scalingo/go-utils/graceful

go 1.25.0

require (
	github.com/Scalingo/go-utils/errors/v3 v3.2.1
	github.com/Scalingo/go-utils/logger v1.12.2
	github.com/cloudflare/tableflip v1.2.3
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.10.2 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
	gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
