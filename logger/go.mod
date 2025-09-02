module github.com/Scalingo/go-utils/logger

go 1.24

retract v1.9.0 // Had accidentally been released rather than a lower version (1.6.0). Hence we retract it and release a new version 1.9.1.

require (
	github.com/Scalingo/logrus-rollbar v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/rollbar/rollbar-go v1.4.8
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/Scalingo/errgo-rollbar v0.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
