//go:build tools

package tools

import (
	_ "github.com/actgardner/gogen-avro/v10/cmd/gogen-avro"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/vektra/mockery/v2/cmd"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
