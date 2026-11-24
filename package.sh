#!/bin/bash -xe

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

NAME=$1
VERSION=$2
ARCH=$(dpkg --print-architecture)
SNAP_DIR=${DIR}/build/snap

${DIR}/apt.sh squashfs-tools

cp -r ${DIR}/bin/* ${SNAP_DIR}/bin
cp -r ${DIR}/config ${SNAP_DIR}
cp -r ${DIR}/meta ${SNAP_DIR}
echo ${VERSION} > ${SNAP_DIR}/version

echo "version: $VERSION" >> ${SNAP_DIR}/meta/snap.yaml
echo "architectures:" >> ${SNAP_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${SNAP_DIR}/meta/snap.yaml

PACKAGE=${NAME}_${VERSION}_${ARCH}.snap
echo ${VERSION} > ${DIR}/version
echo ${PACKAGE} > ${DIR}/package.name
mksquashfs ${SNAP_DIR} ${DIR}/${PACKAGE} -noappend -comp xz -no-xattrs -all-root
mkdir -p ${DIR}/artifact
cp ${DIR}/${PACKAGE} ${DIR}/artifact
