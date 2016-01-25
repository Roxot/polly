#!/bin/bash

echo "Updating codebase..."
git pull
$POLLY_HOME/scripts/restart_server.sh
