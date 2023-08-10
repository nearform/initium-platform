#!/usr/bin/env bash

KIND_PATH="manifests/kind"
KIND_COMPUTED_PATH="kind.local.yaml"

source .envrc

if ! kind get kubeconfig --name ${INITIUM_REPO_NAME} > /dev/null 2>&1; then
helm template ${KIND_PATH} --set k8s_version="${INITIUM_K8S_VERSION}" \
                           --set repo_host_path="${INITIUM_REPO_HOST_PATH}" \
                           --set repo_node_path="${INITIUM_REPO_NODE_PATH}" \
                           --set name="${INITIUM_REPO_NAME}" \
                           --set repo_name="${INITIUM_REPO_NAME}" > ${KIND_COMPUTED_PATH}

kind create cluster --config ${KIND_COMPUTED_PATH}
fi
