#!/bin/bash -xe
DIR=$(cd "$(dirname "$0")" && pwd)
${DIR}/../build/snap/meta/hooks/install --help
${DIR}/../build/snap/bin/cli --help
