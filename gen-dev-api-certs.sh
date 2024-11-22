#!/usr/bin/env bash

tls_dir=$HOME/.local/share/dev-nf6-api/tls

go run cli/main.go gentls -d $tls_dir --ca
go run cli/main.go gentls -d $tls_dir -n server
go run cli/main.go gentls -d $tls_dir -n server --cert
