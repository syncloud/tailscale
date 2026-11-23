#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ROOT=$( cd "${DIR}/.." && pwd )

PROJECT=$1
DISTRO=$2
SPEC=${3:?usage: ui.sh project distro spec}

getent hosts tailscale.${DISTRO}.com | sed "s/tailscale.${DISTRO}.com/auth.${DISTRO}.com/g" | tee -a /etc/hosts

ARTIFACT=${ROOT}/artifact/screenshots
mkdir -p ${ARTIFACT}
trap 'chmod -R a+r ${ARTIFACT} 2>/dev/null; exit' EXIT INT TERM

cd ${ROOT}/web
./e2e/wait-app.sh tailscale.${DISTRO}.com
${ROOT}/npm.sh ci --no-audit --no-fund

CI=true \
PLAYWRIGHT_DOMAIN=${DISTRO}.com \
PLAYWRIGHT_ARTIFACT_DIR=${ARTIFACT} \
DEVICE_USER=user \
DEVICE_PASSWORD=Password1 \
  npx playwright test --project=${PROJECT} "${SPEC}"
