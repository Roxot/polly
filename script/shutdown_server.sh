#!/bin/bash

SERVER_LOG=`echo $POLLY_HOME/server.log | sed 's/ /\\ /g'`
PROC_IDS=`fuser "$SERVER_LOG" 2> /dev/null`

for PID in `echo $PROC_IDS`
do
    if [ $PID != $$ ]; then
        echo "Killing current running instance..."
        kill -9 $PID
    fi
        done

