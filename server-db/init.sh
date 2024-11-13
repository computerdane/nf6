#!/usr/bin/env bash

source ./server-db/env.sh

createdb -h /tmp nf6
psql -h /tmp -d nf6 -f ./server-db/init.sql
