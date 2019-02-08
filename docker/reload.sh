#!/bin/bash
set -euo pipefail

update-remote-config

dnsdist-autoconf -config /etc/dnsdist/autoconf.toml -output /etc/dnsdist/dnsdist.conf -docker

if dnsdist --check-config; then
    supervisorctl restart dnsdist
else
    echo "Something happened, unable to parse auto-generated config file."
fi
