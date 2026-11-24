#!/bin/bash -xe
DIR=$( cd "$( dirname "$0" )" && pwd )
VERSION=$1
ARCH=$(dpkg --print-architecture)
BUILD_DIR=${DIR}/../build/snap/tailscale

${DIR}/../apt.sh wget ca-certificates
mkdir -p ${BUILD_DIR}
${DIR}/../download-retry.sh https://pkgs.tailscale.com/stable/tailscale_${VERSION}_${ARCH}.tgz /tmp/tailscale.tgz
tar xzf /tmp/tailscale.tgz -C /tmp
cp /tmp/tailscale_${VERSION}_${ARCH}/tailscale /tmp/tailscale_${VERSION}_${ARCH}/tailscaled ${BUILD_DIR}
