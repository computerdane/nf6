#!/usr/bin/env bash

socket_dir=/tmp

createdb -h $socket_dir nf6
psql -h $socket_dir -d nf6 -f ./db/init.sql
