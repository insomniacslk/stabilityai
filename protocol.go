package stabilityai

// This file generates the protocol code from the Protobuf speficification.
//
// You need `protoc`, the Go protobuf plugin, and the Go gRPC plugin. Install
// the `protobuf-compiler` (or equivalent) package on your system, the install
// the plugin by running
//   `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
// and
//   `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

//go:generate mkdir -p generation
//go:generate curl https://raw.githubusercontent.com/Stability-AI/stability-sdk/ecma_clients/src/proto/generation.proto -o generation/generation.proto
//go:generate protoc -I. --go_out=generation --go-grpc_out=generation --experimental_allow_proto3_optional generation/generation.proto
