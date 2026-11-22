#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

exec ${DIR}/headscale/headscale -c /var/snap/tailscale/current/config/headscale.yaml "$@"
