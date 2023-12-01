module github.com/Scalingo/go-utils/storage

go 1.20

require (
	github.com/Scalingo/go-utils/logger v1.2.0
	github.com/aws/aws-sdk-go-v2 v1.23.4
	github.com/aws/aws-sdk-go-v2/credentials v1.16.8
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.15.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.47.1
	github.com/aws/smithy-go v1.18.1
	github.com/golang/mock v1.6.0
	github.com/ncw/swift/v2 v2.0.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.2.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.2.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.16.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
