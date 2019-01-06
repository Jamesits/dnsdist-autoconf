#!/bin/bash
set -euo pipefail

if [[ "$1" == supervisord ]] || [ "$1" == dnsdist ]; then
    dnsdist-autoconf -config /etc/dnsdist/autoconf.toml -output /etc/dnsdist/dnsdist.conf -docker
fi

exec "$@"
