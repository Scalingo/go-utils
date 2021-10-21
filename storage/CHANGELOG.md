# Changelog

## To be Released

* Bump github.com/aws/aws-sdk-go-v2 from 0.27.0 to 0.31.0
* Bump github.com/stretchr/testify from 1.6.1 to 1.7.0
* Bump github.com/aws/aws-sdk-go-v2/credentials from 0.1.2 to 0.2.0
* Bump github.com/ncw/swift from 1.0.52 to 1.0.53
* Bump github.com/golang/mock from 1.4.4 to 1.6.0
* Bump go version to 1.16
* Add options to s3 client to control multipart upload:
	* `func WithPartSize(size int64)`
	* `func WithUploadConcurrency(concurrency int)`

## v1.0.0

* Initial breakdown of go-utils into subpackages
