#!/bin/bash

echo "Starting debugger..."
dlv exec /app/goblncr-dbg \
    --headless --listen=0.0.0.0:2345 \
    --continue --accept-multiclient \
    --allow-non-terminal-interactive=true --log # --log-output=rpc