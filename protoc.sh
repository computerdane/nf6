#!/usr/bin/env bash

mkdir -p nf6
protoc --go_out=./nf6 --go_opt=paths=source_relative --go-grpc_out=./nf6 --go-grpc_opt=paths=source_relative nf6.proto
