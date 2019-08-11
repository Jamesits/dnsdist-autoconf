#!/bin/bash
set -Eeuo pipefail

if [ ! -f /tmp/started ]; then
    exit 0
fi

dig @127.0.0.1 baidu.com
dig @127.0.0.1 google.com

exit 0