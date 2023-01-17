#!/usr/bin/env bash

source .envrc

# Generate report of current env vars
echo "======================================================"
printenv | grep "KKA_.*"
echo "======================================================"

# Run Tilt CI
tilt ci --file Tiltfile"${CI_TYPE}"

if [ "${KKA_DEPLOY_MINIMAL}" == "false" ]; then
  # Login on ArgoCD
  argocd login --core --name k8s-kurated-addons
  kubectl config set-context --current --namespace=argocd

  # Ensure ArgoCD apps are all healthy and in sync
  echo ">> Waiting for k8s-kurated-addons to be healty and in sync..."
  while [ true ]
  do
    ALL_HEALTHY=true
    readarray -t apps_health < <(argocd app get k8s-kurated-addons -o json | jq -r '.status.resources | .[]? | select(.kind | contains("Application")) | .health.status')

    if (( ${#apps_health[@]} == 0 )); then
      ALL_HEALTHY=false
    else
      for status in ${apps_health[@]}
      do
        if [ "$status" != "Healthy" ]; then
          ALL_HEALTHY=false
        fi
      done
    fi

    ALL_SYNCED=true
    readarray -t apps_sync < <(argocd app get k8s-kurated-addons -o json | jq -r '.status.resources | .[]? | select(.kind | contains("Application")) | .status')
    if (( ${#apps_sync[@]} == 0 )); then
      ALL_SYNCED=false
    else
      for status in ${apps_sync[@]}
      do
        if [ "$status" != "Synced" ]; then
          ALL_SYNCED=false
        fi
      done
    fi

    if [[ "$ALL_HEALTHY" == "true" && "$ALL_SYNCED" == "true" ]]; then
      break
    else
      argocd app get k8s-kurated-addons

      # Print the 10 last lines of logs of apps not currently healthy
      readarray -t apps < <(argocd app get k8s-kurated-addons -o json | jq -r '.status.resources | .[]? | select(.kind | contains("Application")) | select(.health.status | contains("Progressing")) | .name')
      for app in ${apps[@]}
      do
        echo ">> Printing last 10 log lines for $app..."
        argocd app logs $app --tail 10
      done
    fi

    sleep 15
  done

  # List all the installed apps
  argocd app list

  # Logout on ArgoCD
  argocd logout k8s-kurated-addons
  kubectl config set-context --current --namespace=default
fi
