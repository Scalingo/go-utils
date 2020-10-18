module github.com/Scalingo/go-utils/storage

go 1.14

require (
	github.com/Scalingo/go-utils/logger v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.35.9
	github.com/aws/aws-sdk-go-v2 v0.27.0
	github.com/aws/aws-sdk-go-v2/credentials v0.1.2
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v0.1.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.2.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v0.27.0
	github.com/awslabs/smithy-go v0.2.0
	github.com/golang/mock v1.4.4
	github.com/ncw/swift v1.0.52
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.2.2
)

replace github.com/Scalingo/go-utils/logger => ../logger
