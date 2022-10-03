# Changelog

## To be Released

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
