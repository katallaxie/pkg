#!/bin/bash
# This script is executed after the creation of a new project.

sudo apt-get update
sudo apt install -y direnv protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

cp scripts/pre-commit.sh .git/hooks/pre-commit
cp scripts/pre-push.sh .git/hooks/pre-push
chmod 755 .git/hooks/pre-commit
chmod 755 .git/hooks/pre-push
