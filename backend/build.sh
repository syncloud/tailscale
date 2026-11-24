#!/bin/bash -e
DIR=$(cd "$(dirname "$0")" && pwd)
cd "$DIR"

go test ./... -count=1

mkdir -p ../build/snap/bin
CGO_ENABLED=0 go build -buildvcs=false -o ../build/snap/bin/backend .
