#!/bin/sh
set -e
stat go.mod > /dev/null   # must be in src/
go test ./... -count 1 
staticcheck ./... 
go run internal/check/check.go 
echo "check ok"
