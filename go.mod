module github.com/ontavu/cams-api

go 1.26.7

require (
	google.golang.org/grpc v1.61.1
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0 // indirect
)

// The two protoc plugins, pinned by the module graph rather than by whatever
// is on PATH. protoc-gen-go lives in the same module as the protobuf runtime,
// so pinning the plugin pins the runtime: a newer protoc-gen-go-grpc emits code
// against a newer grpc than the one required above, and the two have to be
// bumped together. `go install tool` builds them; see AGENTS.md.
tool (
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
