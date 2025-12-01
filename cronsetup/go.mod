module github.com/Scalingo/go-utils/cronsetup

go 1.24.0

toolchain go1.24.4

require (
	github.com/Scalingo/go-utils/errors/v3 v3.1.1
	github.com/Scalingo/go-utils/logger v1.11.0
	github.com/gofrs/uuid/v5 v5.4.0
	github.com/sirupsen/logrus v1.9.3
	go.etcd.io/etcd/client/v3 v3.6.6
)

require (
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.etcd.io/etcd/api/v3 v3.6.6 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.6 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/errgo.v1 v1.0.1 // indirect
)

// In Dev you can uncomment the following line to use the local 'logger' package
// replace github.com/Scalingo/go-utils/logger => ../logger
