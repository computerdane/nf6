#!/usr/bin/env bash

go build -o nf.bin main.go
./nf.bin "$@"
