#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

${DIR}/bin/wait-for-configure.sh
export SSL_CERT_FILE=/var/snap/platform/current/syncloud.ca.crt
exec ${DIR}/headscale/headscale serve -c ${SNAP_DATA}/config/headscale.yaml
