#!/usr/bin/env bash

source .envrc

if kind get kubeconfig --name ${INITIUM_REPO_NAME} > /dev/null 2>&1; then
tilt up
fi
