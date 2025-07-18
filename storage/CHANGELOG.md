# Changelog

## To be Released

## v1.7.2

* build(deps): rollback github.com/aws/aws-sdk-go-v2 to v1.32.8

## v1.7.1

* chore(go): corrective bump - Go version regression from 1.24.3 to 1.24

## v1.7.0

* chore(go): upgrade to Go 1.24

## v1.6.0

* feat(s3): Pass the retry options to the AWS SDK on top of our custom retrier.

## v1.5.1

* build(deps): rollback github.com/aws/aws-sdk-go-v2 to v1.32.0

## v1.5.0

* build(deps): bump github.com/Scalingo/go-utils/logger to v1.4.0
* build(deps): bump github.com/aws/aws-sdk-go-v2 to v1.36.0
* build(deps): github.com/aws/aws-sdk-go-v2/credentials to v1.17.57
* build(deps): github.com/aws/aws-sdk-go-v2/feature/s3/manager to v1.17.57
* build(deps): github.com/aws/aws-sdk-go-v2/service/s3 to v1.75.2
* build(deps): github.com/aws/smithy-go to v1.22.2
* feat(errors): Use github.com/Scalingo/go-utils/errors/v2
* feat(s3): Stop using deprecated API to access customer S3 endpoint
* feat(s3): Set default checksum behavior for request and response to "when_required"

## v1.4.0

* feat(s3): Add method `GetWithRetries` taking a writer as argument and which
  will handle gracefully if a `GetObject` sent to the object storage is failing
  and has to be retried.

## v1.3.3

* build(deps): bump github.com/aws/aws-sdk-go-v2/feature/s3/manager from 1.11.37 to 1.11.46
* build(deps): bump github.com/aws/aws-sdk-go-v2/service/s3 from 1.29.1 to 1.29.6
* build(deps): bump github.com/aws/aws-sdk-go-v2/credentials from 1.12.23 to 1.13.3
* build(deps): bump github.com/aws/aws-sdk-go-v2 from 1.17.1 to 1.17.3

## v1.3.2

* feat: `Delete` may return `ObjectNotFound` [#425](https://github.com/Scalingo/go-utils/pull/425)

## v1.3.1

* feat(list): `List` can take a `ListOpts` which can be used to define a `MaxKeys` field. This field limits the amount of objects in the returned list. By default this variable equals 0 which means that S3 will return the maximum objects he can (max is storage.S3ListMaxKeys = 1,000) [#427](https://github.com/Scalingo/go-utils/pull/427)

## v1.3.0

* feat: add methods:
    * List method that returns a slice of string filtered with a prefix.
    * Move method.
* build(deps): bump github.com/aws/aws-sdk-go-v2/credentials from 1.12.21 to 1.12.23
* build(deps): bump github.com/aws/aws-sdk-go-v2/feature/s3/manager from 1.11.34 to 1.11.37
* build(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1

## v1.2.1

* fix: s3.Info Prevent Panic for Content-Type and Checksum if correspondent values returned by s3 are nil pointers [403](https://github.com/Scalingo/go-utils/pull/403)

## v1.2.0

* feat: add Info method, it returns object information and ObjectNotFound custom error in case of object not found [393](https://github.com/Scalingo/go-utils/pull/393)
* build(deps): bump github.com/Scalingo/go-utils/logger from 1.1.1 to 1.2.0
* build(deps): bump github.com/aws/aws-sdk-go-v2 from 1.16.4 to 1.16.16
* build(deps): bump github.com/aws/aws-sdk-go-v2/credentials from 1.12.4 to 1.12.21
* build(deps): bump github.com/aws/aws-sdk-go-v2/feature/s3/manager from 1.11.14 to 1.11.34
* build(deps): bump github.com/aws/aws-sdk-go-v2/service/s3 from 1.26.10 to 1.27.11
* build(deps): bump github.com/aws/smithy-go from 1.11.2 to 1.13.3
* build(deps): bump github.com/stretchr/testify from 1.7.1 to 1.8.0

## v1.1.2

* chore(go): use go 1.17
* build(deps): bump github.com/aws/aws-sdk-go-v2 from 1.9.2 to 1.16.4
* build(deps): bump github.com/aws/aws-sdk-go-v2/credentials from 1.6.3 to 1.12.4
* build(deps): bump github.com/aws/aws-sdk-go-v2/feature/s3/manager from 1.5.4 to 1.11.14
* build(deps): bump github.com/aws/aws-sdk-go-v2/service/s3 from 1.16.1 to 1.26.10
* build(deps): bump github.com/aws/smithy-go from 1.10.0 to 1.11.2
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1

## v1.1.1

* Change github.com/awslabs/smithy-go to github.com/aws/smithy-go

## v1.1.0

* Bump github.com/aws/aws-sdk-go-v2 from 0.27.0 to 1.9.2
* Bump github.com/aws/aws-sdk-go-v2/credentials from 0.1.2 to 1.4.3
* Bump github.com/aws/aws-sdk-go-v2/feature/s3/manager from 0.2.0 to 1.5.4
* Bump github.com/aws/aws-sdk-go-v2/service/s3 0.31.0 to 1.16.1
* Bump github.com/stretchr/testify from 1.6.1 to 1.7.0
* Bump github.com/ncw/swift from 1.0.52 to 1.0.53
* Bump github.com/golang/mock from 1.4.4 to 1.6.0
* Bump github.com/Scalingo/go-utils/logger from 1.0.0 to 1.1.0
* Bump go version to 1.16
* Add options to s3 client to control multipart upload:
	* `func WithPartSize(size int64)`
	* `func WithUploadConcurrency(concurrency int)`

## v1.0.0

* Initial breakdown of go-utils into subpackages
