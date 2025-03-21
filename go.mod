module github.com/Orzelius/cosi-testing

go 1.24.0

replace github.com/cosi-project/runtime => ../cosi-runtime/

require (
	github.com/cosi-project/runtime v0.9.4
	github.com/fsnotify/fsnotify v1.8.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.11.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/siderolabs/gen v0.8.0 // indirect
	github.com/siderolabs/go-pointer v1.0.0 // indirect
	github.com/siderolabs/protoenc v0.2.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.10.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250227231956-55c901821b1e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250227231956-55c901821b1e // indirect
	google.golang.org/grpc v1.70.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
