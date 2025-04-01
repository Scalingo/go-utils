module github.com/Scalingo/go-utils/storage

go 1.23.3
toolchain go1.24.1

require (
	github.com/Scalingo/go-utils/errors/v2 v2.4.0
	github.com/Scalingo/go-utils/logger v1.5.0
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/credentials v1.17.65
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.70
	github.com/aws/aws-sdk-go-v2/service/s3 v1.79.0
	github.com/aws/smithy-go v1.22.3
	github.com/golang/mock v1.6.0
	github.com/ncw/swift/v2 v2.0.3
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.31.0 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
