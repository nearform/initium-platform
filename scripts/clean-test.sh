#!/usr/bin/env bash

set -e

source .envrc

if [ -d "$GOPATH" ]; then
  echo "Deleting ${GOPATH}..."
  sudo rm -rf $GOPATH
fi

if [ -d "$HOME/.kube/cache/oidc-login" ]; then
  echo "Deleting $HOME/.kube/cache/oidc-login..."
  sudo rm -rf $HOME/.kube/cache/oidc-login
fi
