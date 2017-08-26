#!/bin/bash

protoc \
  -I api \
  -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
  --go_out=plugins=grpc:./api \
  --grpc-gateway_out=logtostderr=true:./api \
  --swagger_out=logtostderr=true:./api \
  api/login.proto
