#!/usr/bin/env bash

data_dir=$HOME/.local/share/nf6-db-dev
socket_dir=/tmp

mkdir -p $data_dir
chmod 700 $data_dir

initdb -D $data_dir || true
postgres -D $data_dir -k $socket_dir