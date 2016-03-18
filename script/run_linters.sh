#!/bin/bash

source script/packages.sh

set -eo pipefail

echo "Running go vet..."
go vet github.com/roxot/polly
for package in ${packages[@]}; do
  go vet github.com/roxot/polly/$package
done

echo "Running golint..."
golint github.com/roxot/polly
for package in ${packages[@]}; do
  golint github.com/roxot/polly/$package
done
