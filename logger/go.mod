module github.com/Scalingo/go-utils/logger

go 1.25.0

retract v1.9.0 // Had accidentally been released rather than a lower version (1.6.0). Hence we retract it and release a new version 1.9.1.

require (
	github.com/Scalingo/logrus-rollbar v1.4.5
	github.com/rollbar/rollbar-go v1.4.8
	github.com/sirupsen/logrus v1.10.2
	github.com/stretchr/testify v1.12.1
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

require (
	github.com/Scalingo/errgo-rollbar v0.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
