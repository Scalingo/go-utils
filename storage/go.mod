module github.com/Scalingo/go-utils/storage

go 1.17

require (
	github.com/Scalingo/go-utils/logger v1.2.0
	github.com/aws/aws-sdk-go-v2 v1.17.5
	github.com/aws/aws-sdk-go-v2/credentials v1.13.11
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.50
	github.com/aws/aws-sdk-go-v2/service/s3 v1.30.1
	github.com/aws/smithy-go v1.13.5
	github.com/golang/mock v1.6.0
	github.com/ncw/swift/v2 v2.0.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.23 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
