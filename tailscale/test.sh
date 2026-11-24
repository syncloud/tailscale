#!/bin/bash -xe
DIR=$( cd "$( dirname "$0" )" && pwd )
${DIR}/../build/snap/tailscale/tailscale version
${DIR}/../build/snap/tailscale/tailscaled --version
