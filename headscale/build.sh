#!/bin/bash -xe
DIR=$( cd "$( dirname "$0" )" && pwd )
VERSION=$1
ARCH=$(dpkg --print-architecture)
BUILD_DIR=${DIR}/../build/snap/headscale

${DIR}/../apt.sh wget ca-certificates
mkdir -p ${BUILD_DIR}
${DIR}/../download-retry.sh https://github.com/juanfont/headscale/releases/download/v${VERSION}/headscale_${VERSION}_linux_${ARCH} ${BUILD_DIR}/headscale
chmod +x ${BUILD_DIR}/headscale
