# Changelog

## To be Released

## v1.5.3

* fix(nsqconsumer/errors) make unwrap compatible with all go-utils/errors versions

## v1.5.2

* fix(nsqconsumer.Error) stop unwrap two times when errors containing a nsqconsumer.Error type.

## v1.5.1

* chore(go): corrective bump - Go version regression from 1.24.3 to 1.24

## v1.5.0

* chore(go): upgrade to Go 1.24

## v1.4.1

* fix(disable-backoff): fix configuration of consumer to correctly handle DisableBackoff argument

## v1.4.0

* feat(disable-backoff): add option to completely disable backoff for a consumer

## v1.3.3

* refactor(consumer): match our current best practices around logging and errors

## v1.3.2

* feat: nsqconsumer.Error - `Unwrap` method to be compatible with `errors.Is/As()`
* feat: nsqconsumer.Error - `NoRetry` to get if the message should be retried or not to be consumed.

## v1.3.1

* fix: Use API for github.com/go-utils/logger instead of setting logger manually in context

## v1.3.0

* feat: add a configurable log level

## v1.2.0

* feat: Start unwraps errors to find noRetry field that can be wrapped in ErrCtx. Also use ErrCtx to enrich the logger
* build(deps): bump github.com/Scalingo/go-utils/logger from 1.1.1 to 1.2.0
* build(deps): bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0

## v1.1.1

* chore(go): use go 1.17
* build(deps): bump github.com/nsqio/go-nsq from 1.0.8 to 1.1.0

## v1.1.0

* Bump github.com/sirupsen/logrus from 1.7.0 to 1.8.1
* Bump go version to 1.16
* Bump github.com/Scalingo/go-utils/nsqproducer from 1.0.0 to 1.1.0
* Bump github.com/go-utils/logger from v1.0.0 to v1.1.0

## v1.0.0

* Initial breakdown of go-utils into subpackages
