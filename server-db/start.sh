#!/usr/bin/env bash

source ./server-db/env.sh

mkdir -p "$DIR/data"
chmod 700 "$DIR/data"

initdb -D "$DIR/data"
postgres -D "$DIR/data" -k /tmp
