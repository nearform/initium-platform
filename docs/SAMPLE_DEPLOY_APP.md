# Deploying a Sample App on our Kubernetes cluster with Knative

We will walk through the steps of deploying a sample application using Knative on Kubernetes cluster where different addons are installed and the cluster is ready for serving traffic.

## Prerequisites

- Our type of Kubernetes cluster with all addons should be installed as it is described [here](https://github.com/nearform/initium-platform#readme)
- The kubectl command-line tool for interacting with the Kubernetes cluster.
- Already configured `kubectl` context for your CLI client. If the Kubernetes cluster is installed with the `make` commands from the description [here](https://github.com/nearform/initium-platform#readme), then the context setup question is already solved.
- (Optional) Knative CLI for interacting with your Knative installation.

## Step 1: Deploy a sample app with Knative

We have a very simple Golang application where the output is just a simple string. With only one single command that app will be deployed on the cluster.

```bash
kubectl apply -f examples/sample-deploy-app/sample-deploy-app.yaml
```

The `sample-deploy-app.yaml` file defines the configuration for the Knative Service, then it specifies the Docker image for the app, an ENV variable and a container port.

Wait for the Knative Service to be created and for the app to be deployed, the following command will help you with the first verification:

```bash
kubectl get ksvc
```

You should see the sample app listed with a Ready status.

## Step 2: Test the App

To test the app response, you can run a simple `curl` command against the deployed 'Hello World' service:

Depending on the OS type that you are using, the `INITIUM_LB_ENDPOINT` variable is exported differently, for Linux users we have:

```bash
export INITIUM_LB_ENDPOINT="$(kubectl get service -n istio-ingress istio-ingressgateway -o go-template='{{(index .status.loadBalancer.ingress 0).ip}}'):80"
```

For the MacOS users we have (accessing the Kind cluster is different for MacOS and Windos and it is described [here](https://kind.sigs.k8s.io/docs/user/loadbalancer/)):

```bash
export PROXY_HTTP_CONTAINER_NAME="istio-lb-proxy-80"
export INITIUM_LB_ENDPOINT="127.0.0.1:80"
export INITIUM_LB_INT_HTTP_PORT=$(kubectl get service -n istio-ingress istio-ingressgateway -o go-template='{{range .spec.ports}}{{if (eq .port 80)}}{{.nodePort}}{{end}}{{end}}')
export INITIUM_REPO_NAME="initium-platform"
docker run \
      -d \
      --restart always \
      --name $PROXY_HTTP_CONTAINER_NAME \
      --publish 127.0.0.1:80:$INITIUM_LB_INT_HTTP_PORT \
      --link ${INITIUM_REPO_NAME}-control-plane:target \
      --network kind \
      alpine/socat -dd tcp-listen:$INITIUM_LB_INT_HTTP_PORT,fork,reuseaddr tcp-connect:target:$INITIUM_LB_INT_HTTP_PORT
```

Then you need to export one extra variable and the `curl` command is ready to use:

```bash
export SAMPLE_APP_URL=$(kubectl get ksvc -o json | jq -r '.items[] | select(.metadata.name == "helloworld") | .status.url'  | sed 's#http://##')
curl -H "Host: $SAMPLE_APP_URL" "http://$INITIUM_LB_ENDPOINT"
```

The responded string of 'Hello World' is enough for verification that the application is successfully deployed on the Kubernetes cluster.

## Step 3: Clean Up

To delete the Knative Service and the running app, run the following command:

```bash
kubectl delete -f examples/sample-deploy-app/sample-deploy-app.yaml
```

This will delete the Knative Service and all the resources associated with it, including the app pods.
