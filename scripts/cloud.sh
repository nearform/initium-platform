#!/usr/bin/env bash

source .envrc

# Install ArgoCD
kubectl create namespace argocd
helm install argocd argo/argo-cd
kubectl apply -f https://github.com/nearform/k8s-kurated-addons/releases/download/v0.0.1/app-of-apps.yaml
