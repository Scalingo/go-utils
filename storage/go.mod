module github.com/Scalingo/go-utils/storage

go 1.20

require (
	github.com/Scalingo/go-utils/logger v1.2.0
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/credentials v1.13.30
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.75
	github.com/aws/aws-sdk-go-v2/service/s3 v1.38.0
	github.com/aws/smithy-go v1.14.2
	github.com/golang/mock v1.6.0
	github.com/ncw/swift/v2 v2.0.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.15.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.2 // indirect
	golang.org/x/sys v0.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
