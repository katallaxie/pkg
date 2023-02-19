#!/bin/bash
# This script is executed after the creation of a new project.

sudo apt-get update
sudo apt install -y direnv protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
