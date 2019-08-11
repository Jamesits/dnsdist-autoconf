#!/bin/bash
set -euo pipefail

! update-remote-config.sh

dnsdist-autoconf -config /etc/dnsdist -docker

if dnsdist --check-config; then
    supervisorctl restart dnsdist
else
    echo "Something happened, unable to parse auto-generated config file."
fi
