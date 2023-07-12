# Changelog

## To be Released

* [BREAKING CHANGE]: feat(mongo/document/validation): Add distinction between internal and validation errors ([PR#552](https://github.com/Scalingo/go-utils/pull/552)).

## 1.3.2

feat(pagination): replace BadRequestError by the one from go-handlers to return HTTP 400.
feat(pagination): return HTTP 200 when the page is empty instead of BadRequestError

## v1.3.1

* fix: Close mongo session on CountUnscoped method

## v1.3.0

* feat(count): Add count method to document package

## v1.2.2

* docs: add doc for Session function [#397](https://github.com/Scalingo/go-utils/pull/397)
* build(deps): bump github.com/Scalingo/go-utils/logger from 1.1.1 to 1.2.0
* build(deps): bump github.com/Scalingo/go-utils/errors from 1.1.1 to 2.2.0
* build(deps): bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0
* build(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1

## v1.2.1

* chore(go): use go 1.17
* build(deps): bump github.com/Scalingo/go-utils/errors from 1.0.0 to 1.1.0
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1

## v1.2.0

* Bump go version to 1.16
* Bump github.com/go-utils/logger from v1.0.0 to v1.1.0

## v1.1.1

* Display the MongoDB object ID in the correct format on the logs [#187](https://github.com/Scalingo/go-utils/pull/187)
* Bump github.com/sirupsen/logrus from 1.7.0 to 1.8.1

## v1.1.0

* Add Scalingo pagination support to the package [#140](https://github.com/Scalingo/go-utils/pull/140)

## v1.0.0, v1.0.1

* Initial breakdown of go-utils into subpackages
