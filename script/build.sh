#!/bin/bash

set -eo pipefail

echo "Building github.com/roxot/polly/cmd/pollyserver..."
go install github.com/roxot/polly/cmd/pollyserver
