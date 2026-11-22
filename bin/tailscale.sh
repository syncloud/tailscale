#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

exec ${DIR}/tailscale/tailscale --socket=/var/snap/tailscale/current/tailscaled.sock "$@"
