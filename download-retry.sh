#!/bin/bash -e
URL=$1
TARGET=$2
for i in $(seq 1 10); do
  if wget --tries=3 --timeout=60 -q -O "$TARGET" "$URL"; then
    exit 0
  fi
  echo "retry download $URL"
  sleep 10
done
echo "download failed: $URL"
exit 1
