# Kubernetes Authentication Through Dex and GitHub

## Introduction
Document contains procedure to configure Kubernetes authentication using Dex, GitHub and [Kubelogin plugin](https://github.com/int128/kubelogin).
As access to local `Kind` cluster is already possible with generated kubeconfig, intention is to provide an authentication method
which can be used in general use case.

In this example we will use GitHub personal account. Users can also use their GitHub Enterprise account which is further
explained in [GitHub Dex documentation](https://dexidp.io/docs/connectors/github/#github-enterprise)

## Prerequisites
* It is assumed that local cluster with Dex addon is already deployed. Needed steps are explained in the [bootstrap section](https://github.com/nearform/k8s-kurated-addons#bootstrap).
During cluster bootstrap OIDC parameters are already configured. Parameters can be checked in the [kind manifest](https://github.com/nearform/k8s-kurated-addons/blob/main/manifests/kind.yaml).
* Kubelogin plugin is installed using [setup instructions](https://github.com/int128/kubelogin#setup).
* For now, we will use default configuration for Dex addon and make changes as we progress.

## GitHub Configuration
1. [Create](https://github.com/settings/applications/new ) a new OAuth app setting in GitHub
   * Application name: Dex
   * Homepage URL: any URL related to the app
   * Authorization callback URL: https://dex.kube.local:32000/callback

2. Open registered app and generate a new client secret. Registered application can be found under `Account Settings/Developer settings/Oauth Apps`.
Save the GitHub client id and the GitHub client secret for later use.

## Configure Secrets
1. Generate Kubelogin client secret which will be used for communication between Kubelogin and Dex
2. Export prepared values as environment variables:
```bash
export GITHUB_CLIENT_ID=[github-client-id]
export GITHUB_CLIENT_SECRET=[github-client-secret]
export KUBELOGIN_CLIENT_SECRET=[kubelogin-client-secret]
```
3. Create `github-client` Kubernetes secret
```bash
kubectl create secret generic github-client \
  --namespace=dex \
  --from-literal=client-id=${GITHUB_CLIENT_ID} \
  --from-literal=client-secret=${GITHUB_CLIENT_SECRET}
```
4. Create `kubelogin-client` Kubernetes secret
```bash
kubectl create secret generic kubelogin-client \
  --namespace=dex \
  --from-literal=client-secret=${KUBELOGIN_CLIENT_SECRET}
```

## Configure Dex
1. Uncomment GitHub section in Dex [values.yaml](https://github.com/nearform/k8s-kurated-addons/blob/main/addons/dex/values.yaml) file. Also, remove default section at the start of the file. Resulting config should look like below:
```yaml
dex-source:
  config:
    issuer: https://dex.kube.local:32000
    storage:
      type: kubernetes
      config:
        inCluster: true
    connectors:
    - type: github
      id: github
      name: GitHub
      config:
        clientID: $GITHUB_CLIENT_ID
        clientSecret: $GITHUB_CLIENT_SECRET
        redirectURI: https://dex.kube.local:32000/callback

    staticClients:
    - id: kubelogin
      redirectURIs:
        - 'http://localhost:8000'
      name: 'Example App'
      secretEnv: KUBELOGIN_CLIENT_SECRET
  envVars:
  - name: GITHUB_CLIENT_ID
    valueFrom:
      secretKeyRef:
        name: github-client
        key: client-id
  - name: GITHUB_CLIENT_SECRET
    valueFrom:
      secretKeyRef:
        name: github-client
        key: client-secret
  - name: KUBELOGIN_CLIENT_SECRET
    valueFrom:
      secretKeyRef:
        name: kubelogin-client
        key: client-secret

istioConfig:
  gateway: istio-ingress/kube-gateway
  host: dex.kube.local
  port: 5556
  dexServiceName: dex-dex-source
```
2. Commit changes and make sure ArgoCD deployed Dex with the new configuration.

## Accessing Kind Cluster
