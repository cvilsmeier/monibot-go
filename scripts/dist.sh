#!/bin/sh

if [ ! -f go.mod ]
then
    echo "go.mod not found, must be in src directory"
    exit 1
fi

rm -rf _dist
mkdir _dist

echo "go build ./cmd/moni"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o _dist ./cmd/moni

TARFILE="moni-linux-amd64.tar.gz"
echo "tar $TARFILE"
tar -czf _dist/$TARFILE README.md CHANGELOG.md -C _dist moni
