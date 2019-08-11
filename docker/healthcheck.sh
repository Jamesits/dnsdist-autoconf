#!/bin/bash

set -Eeuo pipefail

dig @127.0.0.1 baidu.com
dig @127.0.0.1 google.com

exit 0