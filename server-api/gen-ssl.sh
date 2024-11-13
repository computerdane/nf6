#!/usr/bin/env bash

export OPENSSL_CONF=$(pwd)/server-api/openssl.cnf

mkdir -p ~/.nf6/server-api/ssl
cd ~/.nf6/server-api/ssl

# generate ca key and cert
openssl genpkey -algorithm ED25519 > ca.key
openssl req -new -x509 -key ca.key -out ca.crt

# generate server key, csr, and cert
openssl genpkey -algorithm ED25519 > server.key
openssl req -new -key server.key -out server.req
openssl x509 -req -CA ca.crt -CAkey ca.key -in server.req -out server.crt
rm server.req
