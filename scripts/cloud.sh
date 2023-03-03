#!/usr/bin/env bash

source .envrc

# Generate report of current env vars
echo "======================================================"
printenv | grep "KKA_.*"
echo "======================================================"

# Install ArgoCD
kubectl create namespace argocd
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
helm install argocd argo/argo-cd --namespace=argocd -f ./addons/argocd/values.yaml --version $ARGOCD_VERSION
# Ensure ArgoCD apps are all healthy and in sync
echo ">> Waiting for argocd to be healty and in sync..."

while true; do
    # Check the status of the ArgoCD deployment
    deployment_status=$(kubectl get deployment -n argocd argocd-server -o jsonpath='{.status.readyReplicas}')

    if [ "$deployment_status" != "1" ]; then
      echo "ArgoCD deployment is not ready."
    else
      # Check the status of the ArgoCD pod
      pod_status=$(kubectl get pod -n argocd -l app.kubernetes.io/name=argocd-server -o jsonpath='{.items[].status.phase}')

      if [ "$pod_status" != "Running" ]; then
        echo "ArgoCD pod is not running. Current pod status: $pod_status"
      else
        echo "ArgoCD is installed and ready to use."
        break;
      fi
    fi
  sleep 30
done
