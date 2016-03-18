#!/bin/bash

set -eo pipefail

echo "Go get godep..."
go get github.com/tools/godep

echo "Go get golint..."
go get -u github.com/golang/lint/golint

echo "Fetching build dependencies..."
godep restore
