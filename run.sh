#! /bin/sh
set -e

export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0

docker build -t scraping .