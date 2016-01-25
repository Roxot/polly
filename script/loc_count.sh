#!/bin/bash

(find "$POLLY_HOME/src/polly" -name "*.go" -print0 && find "$POLLY_HOME/src/pollyserver" -name "*.go" -print0 && find "$POLLY_HOME/scripts" -name "*.sh" -print0) | xargs -0 wc -l
