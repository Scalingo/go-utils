# Changelog

## To be Released

## v1.3.0

* refactor: only use `github.com/Scalingo/go-utils/errors/v3` for errors [BREAKING CHANGE]
* build(deps): bump `github.com/Scalingo/go-utils/nsqproducer` from 1.4.0 to 1.4.1
* build(deps): bump `github.com/Scalingo/go-utils/logger` from 1.12.0 to 1.12.1
* build(deps): bump `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/metric`, and `go.opentelemetry.io/otel/trace` from 1.40.0 to 1.43.0
* build(deps): bump `golang.org/x/sys` from 0.41.0 to 0.42.0

## v1.2.2

* refactor: replace `github.com/golang/mock` with `go.uber.org/mock`
* chore(go): upgrade to Go 1.25

## v1.2.1

* chore(go): corrective bump - Go version regression from 1.24.3 to 1.24

## v1.2.0

* chore(go): upgrade to Go 1.24
* build(deps): bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0
* build(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1

## v1.1.2

* chore(go): use go 1.17
* build(deps): bump github.com/Scalingo/go-utils/env from 1.0.1 to 1.1.0
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1

## v1.1.1

* Bump github.com/go-utils/logger from v1.0.0 to v1.1.0

## v1.1.0

* Add Publish timeout support: let a user configure how much time should we wait before going to the next producer
* Bump github.com/golang/mock from 1.4.4 to 1.5.0
* Bump github.com/sirupsen/logrus from 1.7.0 to 1.8.1
* Bump go version to 1.16
* bump github.com/Scalingo/go-utils/nsqproducer from 1.0.0 to 1.1.0

## v1.0.0

* Initial breakdown of go-utils into subpackages
