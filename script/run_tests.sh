#!/bin/bash

source script/packages.sh

set -eo pipefail

go test github.com/roxot/polly
for package in ${packages[@]}; do
  go test github.com/roxot/polly/$package
done
