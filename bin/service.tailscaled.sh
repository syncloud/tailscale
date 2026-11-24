#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

${DIR}/bin/wait-for-configure.sh
exec ${DIR}/tailscale/tailscaled --statedir=${SNAP_DATA}/tailscaled --socket=${SNAP_DATA}/tailscaled.sock --port=41641
