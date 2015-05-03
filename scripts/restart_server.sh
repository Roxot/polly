#!/bin/bash

./$POLLY_HOME/scripts/shutdown_server.sh

echo "Building server..."
go install pollyserver

echo "Starting server..."
pollyserver 2> /dev/null &

echo "Server started."
