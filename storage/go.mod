module github.com/Scalingo/go-utils/storage

go 1.16

require (
	github.com/Scalingo/go-utils/logger v1.0.0
	github.com/aws/aws-sdk-go-v2 v0.31.0
	github.com/aws/aws-sdk-go-v2/credentials v0.2.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v0.2.0
	github.com/aws/aws-sdk-go-v2/service/s3 v0.31.0
	github.com/awslabs/smithy-go v0.3.0
	github.com/golang/mock v1.5.0
	github.com/ncw/swift v1.0.53
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20201112073958-5cba982894dd // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
