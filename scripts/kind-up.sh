#!/usr/bin/env bash

KIND_PATH="manifests/kind.yaml"
KIND_COMPUTED_PATH="kind.local.yaml"

source .envrc

if ! kind get kubeconfig --name ${KKA_REPO_NAME} > /dev/null 2>&1; then
cat ${KIND_PATH} | envsubst > ${KIND_COMPUTED_PATH}
kind create cluster --config ${KIND_COMPUTED_PATH} --name ${KKA_REPO_NAME}
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl apply -f https://github.com/nearform/k8s-kurated-addons/releases/download/v0.0.1/app-of-apps.yaml
fi
