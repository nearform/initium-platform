# Kubernetes Authentication Through Dex and GitHub

## Introduction
The document contains a procedure to configure Kubernetes authentication using Dex, GitHub, and [Kubelogin plugin](https://github.com/int128/kubelogin).
As the access to local `Kind` cluster is already possible with generated kubeconfig, the intention is to provide an
authentication method that can be used in a general use case.

In this example, we will use GitHub personal account. Users can also use their GitHub Enterprise account which is further
explained in [GitHub Dex documentation](https://dexidp.io/docs/connectors/github/#github-enterprise)

## Prerequisites
* It is assumed that local cluster with Dex addon is already deployed. The needed steps are explained in the [bootstrap section](https://github.com/nearform/initium-platform#bootstrap).
During cluster bootstrap OIDC parameters are already configured. Parameters can be checked in the [kind manifest](https://github.com/nearform/initium-platform/blob/main/manifests/kind/templates/cluster.yaml).
* The Kubelogin plugin is installed using [setup instructions](https://github.com/int128/kubelogin#setup).
* For now, we will use default configuration for Dex addon and make changes as we progress.

## Local Environment Configuration
1. Retrieve the IP address allocated by MetalLB to the service of type `Loadbalancer`. This address will make Dex issuer reachable on the local network.
```bash
kubectl get svc -n istio-ingress istio-ingressgateway \
  -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```
2. Make dex domain resolves to the Loadbalancer IP address by configuring `/etc/hosts` file:
```bash
172.18.255.200  dex.kube.local
```
3. Add self-signed CA created during cluster bootstrap to trusted certificate store of your machine. For example, if using Ubuntu:
```bash
sudo apt-get install -y ca-certificates
sudo cp .ssl/ca.pem /usr/local/share/ca-certificates/ca.crt
sudo update-ca-certificates
```

**WSL Users Only**:
As WSL networking will not allow Loadbalancer address to be reachable from the Windows host, additional steps are needed:
4. Start TCP forwarding between the WSL LAN address and the Loadbalancer address. Utility `socat` can be used for that purpose.
```bash
sudo socat TCP4-LISTEN:443,fork,reuseaddr TCP4:172.18.255.200:443 &
```
5. Instead of using the Loadbalancer address in WSL hosts file, configure the WSL LAN address in the Windows hosts file, for example:
```bash
172.23.41.241  dex.kube.local
```
6. Make sure **not** to use `generateResolvConf=false` otherwise above setting will not be reflected in WSL

## GitHub Configuration
1. [Create](https://github.com/settings/applications/new ) a new OAuth app setting in GitHub
   * Application name: Dex
   * Homepage URL: any URL related to the app
   * Authorization callback URL: https://dex.kube.local/callback

2. Open the registered app and generate a new client secret. Registered application can be found under `Account Settings/Developer settings/Oauth Apps`.
Save the GitHub client id and the GitHub client secret for later use.

## Configure Secrets
1. Generate Kubelogin client secret that will be used for communication between Kubelogin and Dex
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
1. Uncomment the GitHub section in Dex [values.yaml](https://github.com/nearform/initium-platform/blob/main/addons/dex/values.yaml) file. Also, remove the default section at the start of the file. Resulting config should look like the below:
```yaml
dex-source:
  nameOverride: dex
  config:
    issuer: https://dex.kube.local
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
        redirectURI: https://dex.kube.local/callback

    staticClients:
    - id: kubelogin
      redirectURIs:
        - 'http://localhost:8000'
      name: 'Kubelogin'
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

virtualService:
  name: dex
  namespace: dex
  gateway: istio-ingress/kube-gateway
  host: dex.kube.local
  port: 5556
  serviceName: dex
```
2. Commit changes and make sure ArgoCD deployed Dex with the new configuration.

## Accessing Kind Cluster
1. Create cluster role binding that binds GitHub user (its email address) to a read-only cluster role. At this stage changes to the cluster
are still made using kubeconfig generated at cluster bootstrap.
```bash
kubectl create clusterrolebinding oidc-cluster-viewer --clusterrole=view --user='[github_email_address]'
```
2. Set up the Kubeconfig
```bash
kubectl config set-credentials oidc \
  --exec-api-version=client.authentication.k8s.io/v1beta1 \
  --exec-command=kubectl \
  --exec-arg=oidc-login \
  --exec-arg=get-token \
  --exec-arg=--oidc-issuer-url=https://dex.kube.local \
  --exec-arg=--oidc-client-id=kubelogin \
  --exec-arg=--oidc-extra-scope=email \
  --exec-arg=--oidc-client-secret=[plain-text-kubelogin-client-secret]
```
3. Verify cluster access by using `oidc` user. In this step `kubelogin` will open a browser and redirect user to GiHub login page.
```bash
kubectl --user=oidc get pods -A
```
4. Authenticate, authorize Dex application to use your GitHub account and grant using your email address.
If everything is successful you should see pods running in the cluster.
5. Switch the default context to oidc user and try deleting a pod.
```bash
kubectl config set-context --current --user=oidc
kubectl delete pod argocd-application-controller-0 -n argocd
```
This operation should be forbidden as we now use a read-only role.
