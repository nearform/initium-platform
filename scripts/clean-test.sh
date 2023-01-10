#!/usr/bin/env bash

set -e

source .envrc

if [ -d "$GOPATH" ]; then
  echo "Deleting ${GOPATH}..."
  sudo rm -rf $GOPATH
fi
