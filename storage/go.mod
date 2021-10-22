module github.com/Scalingo/go-utils/storage

go 1.16

require (
	github.com/Scalingo/go-utils/logger v1.1.0
	github.com/aws/aws-sdk-go-v2 v1.9.2
	github.com/aws/aws-sdk-go-v2/credentials v1.4.3
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.5.4
	github.com/aws/aws-sdk-go-v2/service/s3 v1.16.1
	github.com/awslabs/smithy-go v0.3.0
	github.com/golang/mock v1.6.0
	github.com/ncw/swift v1.0.53
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20211020174200-9d6173849985 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
