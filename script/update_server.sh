#!/bin/bash

echo "Updating codebase..."
git pull
$POLLY_HOME/script/restart_server.sh
