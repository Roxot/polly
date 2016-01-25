#!/bin/bash

SHUTDOWN_SCRIPT=`echo $POLLY_HOME/scripts/shutdown_server.sh | sed 's/ /\\ /g'`
"$SHUTDOWN_SCRIPT"

echo "Building server..."
go install pollyserver

echo "Starting server..."
pollyserver 2> /dev/null &

echo "Server started."
