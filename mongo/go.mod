module github.com/Scalingo/go-utils/mongo

go 1.24.0

require (
	github.com/Scalingo/go-handlers v1.11.0
	github.com/Scalingo/go-utils/errors/v3 v3.2.0
	github.com/Scalingo/go-utils/logger v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.4
	github.com/stretchr/testify v1.11.1
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

require (
	github.com/Scalingo/go-utils/crypto v1.1.1 // indirect
	github.com/Scalingo/go-utils/errors/v2 v2.5.1 // indirect
	github.com/Scalingo/go-utils/security v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofrs/uuid/v5 v5.4.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/urfave/negroni/v3 v3.1.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.64.0 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Uncomment if you want to use the local version of these packages (for development purpose)
// replace github.com/Scalingo/go-utils/logger => ../logger
// replace github.com/Scalingo/go-utils/errors => ../errors
