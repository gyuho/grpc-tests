#!/usr/bin/env bash
set -e

if ! [[ $(protoc --version) =~ "3.3.0" ]]; then
  echo "could not find protoc 3.3.0, is it installed + in PATH?"
  exit 255
fi

curl -o ./git-grpc-go.json https://api.github.com/repos/grpc/grpc-go/git/refs/heads/master

go get -v -u github.com/golang/protobuf/protoc-gen-go

echo "Generating proto"

protoc \
  --go_out=plugins=grpc:. \
  service.proto

echo "Generated proto"
