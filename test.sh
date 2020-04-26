#!/bin/bash

set -e
set -u

./resources/scripts/wait-for-it.sh $MYSQL_HOST -t 30 -s

go test -tags=integration -p=1 -race ./...
