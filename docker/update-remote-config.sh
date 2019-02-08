#!/bin/bash
set -euo pipefail

if [ ! -z "$REMOTE_CONFIG" ]; then
    echo "Trying to update config from $REMOTE_CONFIG ..."
    if hash wget 2>/dev/null; then
        wget -O /etc/dnsdist/autoconf.toml "$REMOTE_CONFIG"
    else
        curl -o /etc/dnsdist/autoconf.toml "$REMOTE_CONFIG"
    fi
else
    echo "No remote config URL found, no need to update"
fi

echo "Config updated successfully."