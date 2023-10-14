#!/bin/bash

set -x
set +e

kill -9 $(lsof -t -i :3001)> /dev/null 2>&1

if [ $1 == 'run' ]; then
    go run main.go
fi