#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

${DIR}/bin/wait-for-configure.sh
/bin/rm -f ${SNAP_COMMON}/web.socket
exec ${DIR}/bin/backend --config ${SNAP_DATA}/config/backend.json
