#!/usr/bin/env bash

set -e

source .envrc

PROJECT_TYPE="${1:-unit}"
PROJECT_PATH="${2:-tests}"
PROXY_HTTP_CONTAINER_NAME="istio-lb-proxy-80"

if [ "$PROJECT_TYPE" == "integration" ]; then
  export INITIUM_LB_ENDPOINT="$(kubectl get service -n istio-ingress istio-ingressgateway -o go-template='{{(index .status.loadBalancer.ingress 0).ip}}'):80"

  case "$OSTYPE" in
    darwin*)
      INITIUM_LB_INT_HTTP_PORT=$(kubectl get service -n istio-ingress istio-ingressgateway -o go-template='{{range .spec.ports}}{{if (eq .port 80)}}{{.nodePort}}{{end}}{{end}}')

    # Run a reverse proxy that will forward traffic between host and the LB, you can reach the LB HTTP port via 127.0.0.1:80
    docker run \
      -d \
      --restart always \
      --name $PROXY_HTTP_CONTAINER_NAME \
      --publish 127.0.0.1:80:$INITIUM_LB_INT_HTTP_PORT \
      --link ${INITIUM_REPO_NAME}-control-plane:target \
      --network kind \
      alpine/socat -dd tcp-listen:$INITIUM_LB_INT_HTTP_PORT,fork,reuseaddr tcp-connect:target:$INITIUM_LB_INT_HTTP_PORT || true

      export INITIUM_LB_ENDPOINT="127.0.0.1:80"
    ;;
    msys*)
    ;;
    linux*)
    ;;
  esac
fi

echo "Running test on ${PROJECT_PATH}"

cd $PROJECT_PATH
go mod init test || true
go install . || true
go test . -v -timeout=30m

if [ "$PROJECT_TYPE" == "integration" ]; then
  case "$OSTYPE" in
    darwin*)
      docker stop $PROXY_HTTP_CONTAINER_NAME
      docker rm $PROXY_HTTP_CONTAINER_NAME
    ;;
    msys*)
    ;;
    linux*)
    ;;
  esac
fi
