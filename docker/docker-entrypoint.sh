#!/bin/bash
set -euo pipefail

! update-remote-config.sh

if [[ "$1" == supervisord ]] || [ "$1" == dnsdist ]; then
    dnsdist-autoconf -config /etc/dnsdist -docker
fi

# indicate dnsdist has been started for the healthcheck script
touch /tmp/started

exec "$@"
