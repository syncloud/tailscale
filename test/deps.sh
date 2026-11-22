#!/bin/bash -e
DIR=$( cd "$( dirname "$0" )" && pwd )
TAILSCALE=$1
ARCH=$(dpkg --print-architecture)

${DIR}/../apt.sh sshpass openssh-client wget
pip install -r ${DIR}/requirements.txt
${DIR}/../download-retry.sh https://pkgs.tailscale.com/stable/tailscale_${TAILSCALE}_${ARCH}.tgz /tmp/tailscale.tgz
mkdir -p /tmp/tailscale-client
tar xzf /tmp/tailscale.tgz -C /tmp/tailscale-client --strip-components=1
