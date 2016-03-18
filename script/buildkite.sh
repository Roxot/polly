#!/bin/bash

set -eo pipefail

# Set up environment (can't be done with hooks)
echo "~~~ Setting up environment"
export GOPATH=$(cd ../../../../ && pwd)
echo "Set GOPATH to $GOPATH"
export PATH=$PATH:$GOPATH/bin
echo "Updated PATH to $PATH"

# Fetch and test for dependencies
echo "~~~ Fetching dependencies"
./script/dependencies.sh

# Build the server
echo "~~~ Building server"
./script/build.sh

# Run tests
echo "~~~ Running tests"
./script/run_tests.sh

# Run vet
echo "~~~ Running linters"
./script/run_linters.sh
