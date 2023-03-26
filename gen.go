package pkg

//go:generate protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/urn.proto proto/snowflake.proto proto/access.proto
