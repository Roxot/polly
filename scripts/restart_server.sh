#!/bin/bash

PROCIDS=`fuser ../server.log 2>/dev/null`

for PID in `echo $PROCIDS`
do
    if [ $PID != $$ ]; then
        echo "Killing current running instance..."
        kill -9 $PID
    fi
        done

echo "Building server..."
go install pollyserver

echo "Starting server..."
pollyserver > /dev/null 2>&1 &

echo "Server started."
