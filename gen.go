package main

//go:generate echo "Generating..."
//go:generate protoc -I. --go_out=paths=source_relative:. urn/urn.proto
//go:generate protoc -I. --go_out=paths=source_relative:. snowflake/snowflake.proto
