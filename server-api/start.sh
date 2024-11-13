#!/usr/bin/env bash

SSL_DIR=~/.nf6/server-api/ssl
mkdir -p "$SSL_DIR"
go run server-api/main.go -ssl-dir="$SSL_DIR"
