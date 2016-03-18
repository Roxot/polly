#!/bin/bash

source script/packages.sh

update_status() {
  cur_status=$1
  exit_status=$2
  if [[ $cur_status -gt 0 || $exit_status -gt 0 ]]; then
    echo "1"
  fi
}

script_status=0

echo "Running go vet..."
go vet github.com/roxot/polly
script_status=$(update_status $script_status $?)
for package in ${packages[@]}; do
  go vet github.com/roxot/polly/$package
  script_status=$(update_status $script_status $?)
done

# Allow code style violations for now
echo "Running golint..."
golint github.com/roxot/polly
for package in ${packages[@]}; do
  golint github.com/roxot/polly/$package
done

exit $script_status
