# Deploying a Sample App with sample Prometheus metrics

This application will be deployed with `knative` same as the other sample application where the whole process is described [here](https://github.com/nearform/k8s-kurated-addons/blob/main/docs/SAMPLE_DEPLOY_APP.md).

For exposing the application in Prometheus a `PodMonitor` CRD is used, [here](https://github.com/prometheus-operator/prometheus-operator#customresourcedefinitions) you can read more about the Prometheus CustomResourceDefinitions.

## Prerequisites

- Our type of Kubernetes cluster with all addons should be installed as it is described [here](https://github.com/nearform/k8s-kurated-addons#readme)
- The kubectl command-line tool for interacting with the Kubernetes cluster.
- Already configured `kubectl` context for your CLI client. If the Kubernetes cluster is installed with the `make` commands from the description [here](https://github.com/nearform/k8s-kurated-addons#readme), then the context setup question is already solved.

## Step 1: Deploy a sample app with Knative

This is a sample Golang application where the sample Prometheus metrics are exported using the official Prometheus [client](https://github.com/prometheus/client_golang).

```bash
kubectl apply -f examples/sample-prometheus-app/sample-prometheus-app.yaml
```

## Step 2: Setup verification

First, we need to get the service endpoint, then let's make a few `curl` requests with the `curl` command mentioned below.

Depending on the OS type that you are using, the `KKA_LB_ENDPOINT` variable is exported differently, for Linux users we have:

```bash
export KKA_LB_ENDPOINT="$(kubectl get service -n istio-ingress istio-ingressgateway -o go-template='{{(index .status.loadBalancer.ingress 0).ip}}'):80"
```

For the MacOS users we have (accessing the Kind cluster is different for MacOS and Windos and it is described [here](https://kind.sigs.k8s.io/docs/user/loadbalancer/)):

```bash
export PROXY_HTTP_CONTAINER_NAME="istio-lb-proxy-80"
export KKA_LB_ENDPOINT="127.0.0.1:80"
export KKA_LB_INT_HTTP_PORT=$(kubectl get service -n istio-ingress istio-ingressgateway -o go-template='{{range .spec.ports}}{{if (eq .port 80)}}{{.nodePort}}{{end}}{{end}}')
export KKA_REPO_NAME="k8s-kurated-addons"
docker run \
      -d \
      --restart always \
      --name $PROXY_HTTP_CONTAINER_NAME \
      --publish 127.0.0.1:80:$KKA_LB_INT_HTTP_PORT \
      --link ${KKA_REPO_NAME}-control-plane:target \
      --network kind \
      alpine/socat -dd tcp-listen:$KKA_LB_INT_HTTP_PORT,fork,reuseaddr tcp-connect:target:$KKA_LB_INT_HTTP_PORT
```

Then you need to export one extra variable and the `curl` command is ready to use:

```bash
export SAMPLE_APP_URL=$(kubectl -n prometheus get ksvc -o json | jq -r '.items[] | select(.metadata.name == "sample-prometheus-app") | .status.url'  | sed 's#http://##')
curl -H "Host: $SAMPLE_APP_URL" "http://$KKA_LB_ENDPOINT"
```

Let's open the Prometheus targets for our internal Prometheus Stack. First, we need to expose the service using the `port-forward` feature:

```bash
kubectl -n prometheus port-forward prometheus-kube-prometheus-stack-kube-prometheus-0 9090
```

Then just open http://127.0.0.1:9090/targets in your browser and a new target should be listed.

## Step 3: Grafana

The manifest from where the Prometheus sample app is deployed contains a Grafana sample dashboard as well. For verifying if the metrics from the sample Prometheus app are correctly fetched in Grafana just expose the Grafana service locally and open the dashboard in browser.

```bash
export GRAFANA_POD=$(kubectl -n prometheus get pods -l "app.kubernetes.io/name=grafana" -o json | jq -r '.items[].metadata.name')
kubectl -n prometheus port-forward $GRAFANA_POD 3000
```

Get the Grafana's admin username and password:

```bash
export GRAFANA_SECRET=$(kubectl -n prometheus get secrets -l "app.kubernetes.io/name=grafana" -o json | jq -r '.items[].metadata.name')
kubectl -n prometheus get secrets kube-prometheus-stack-grafana -o jsonpath='{.data.admin-user}' | base64 -d
kubectl -n prometheus get secrets kube-prometheus-stack-grafana -o jsonpath='{.data.admin-password}' | base64 -d
```

Open http://127.0.0.1:3000 in brower and search for a dashboard named as `Sample Prometheus App`.

## Step 4: Clean Up

To delete the deployed Sample Prometheus app just execute the following command:

```bash
kubectl delete -f examples/sample-prometheus-app/sample-prometheus-app.yaml
```
