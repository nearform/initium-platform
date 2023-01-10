#!/usr/bin/env bash

KIND_PATH="manifests/kind.yaml"
KIND_COMPUTED_PATH="kind.local.yaml"

source .envrc

if ! kind get kubeconfig --name ${KKA_REPO_NAME} > /dev/null 2>&1; then
cat ${KIND_PATH} | envsubst > ${KIND_COMPUTED_PATH}
kind create cluster --config ${KIND_COMPUTED_PATH} --name ${KKA_REPO_NAME}
fi
