#!/usr/bin/env bash

mkdir -p ~/.nf6/server-api/ssl
cd ~/.nf6/server-api/ssl

openssl genrsa -out root.key 4096
chmod 400 root.key
openssl req -new -x509 -key root.key -sha256 -out root.crt
chmod 444 root.crt
