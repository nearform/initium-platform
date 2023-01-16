# App Of Apps Deployment Options

## Introduction
Before ArgoCD can install addons it needs bootstrapping. For that purpose `App of Apps` pattern is used.
Below are presented 3 options for `app-of-apps` deployment with the focus on the end users installing it on their own cluster.
Decision made will also affect the usage of the solution with local cluster (tilt powered).

## Option 1: Deployment with Manifest and Inline Values (Current Solution)

Release process generates several `app-of-apps.yaml` files. Each app-of-apps.yaml file present different
group of addons deployed together (addon flavor). Each app-of-apps flavor is deployed using kubectl:
```bash
kubectl apply -f app-of-apps.yaml
```
Example `app-of-apps.yaml` is shown below:
```yaml
spec:
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  project: default
  source:
    path: app-of-apps
    repoURL: https://github.com/nearform/k8s-kurated-addons.git
    targetRevision: 0.0.1
    helm:
      values: |
        repoURL: https://github.com/nearform/k8s-kurated-addons.git
        subChartsRevision: 0.0.1
        apps:
          dex:
            excluded: false
          istio-base:
            excluded: false
          istiod:
            excluded: false
          istio-ingressgateway:
            excluded: false
          knative-operator:
            excluded: false
          knative-serving:
            excluded: false
          knative-eventing:
            excluded: false
          kube-prometheus-stack:
            excluded: false
```
Characteristics:
* `app-of-apps` source points to k8s-kurated-addons repo: https://github.com/nearform/k8s-kurated-addons.git
* On the defined path (app-of-apps) there is a helm chart responsible for deploying multiple application resources.
* Values file for `app-of-apps` is located on the same path in the `values.yaml` file, and it's applied by default (no need for referencing it)
* Inline values (values defined in the body of app-of-apps manifest) define which addons are deployed and are base for making different addon flavors
* Inline values have priority over the values defined in `values.yaml`
* Disable/enable of a specific addon is done by changing `app-of-apps` manifest and redeploying it using kubectl


**Pro**
* Simplified deployment
* Simplified release process

**Contra**
* Using inline values that override values in values.yaml and passing those values to individual applications makes it
difficult to understand and troubleshoot
* Going away from the GitOps principle: disable/enable of a specific addon requires manual `kubectl apply`


## Option 2: Deployment with Helm and File Based Values

New chart `argocd-resources` is created and generated as an artifact. Different files, that represent group of addons deployed together (addons flavor),
are named as `values-[flavor_name].yaml` and kept in the repository in the app-of-apps path. Specific flavor is installed using `argocd-resources` helm chart
while specifying which addons flavor to use:
```bash
helm repo add argocd-resources https://nearform.github.io/helm-charts
helm install --set applications[0].source.helm.valueFiles={values-flavor1.yaml} k8-kurated-addons argocd-resources
```
Example default `values.yaml` is shown below:
```yaml
applications:
  - name: k8s-kurated-addons
    namespace: argocd
    additionalLabels: {}
    additionalAnnotations: {}
    finalizers:
    - resources-finalizer.argocd.argoproj.io
    project: default
    source:
      repoURL: https://github.com/nearform/k8s-kurated-addons.git
      targetRevision: HEAD
      path: app-of-apps
      helm:
        passCredentials: true
        valueFiles:
          - values.yaml
    destination:
      server: https://kubernetes.default.svc
```
Characteristics:
* `app-of-apps` source points to k8s-kurated-addons repo: https://github.com/nearform/k8s-kurated-addons.git
* On the defined path (app-of-apps) there is a helm chart responsible for deploying multiple application resources.
* Specific flavor is selected by pointing valueFiles reference - done by `--set` in the above example
* There are no inline values
* Disable/enable of specific addon is done by committing change to specific values.yaml file where ArgoCD manages deployment

**Pro**
* By not using inline values, state of what is deployed can be easily tracked using git
* Disable/enable of specific addon is controlled by values.yaml file and ArgoCD
* Using Helm chart and not just a manifest make it easier to extend the solution for future use (project, application sets)

**Contra**
* It's not clear what is released by this project
* Deployment requires Helm


## (DRAFT) Option 3: Generate Manifest Using Helm Chart and Use Inline Values

New simplified Helm chart is created - it contains only application resource and limited set of exposed variables. Using `helm template` release process
generates `app-of-apps.yaml` files. Each app-of-apps.yaml file present different group of addons deployed together (addon flavor).
Each app-of-apps flavor is deployed using kubectl:
```bash
kubectl apply -f app-of-apps.yaml
```
Example `app-of-apps.yaml` is shown below:
```yaml
spec:
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  project: default
  source:
    path: app-of-apps
    repoURL: https://github.com/nearform/k8s-kurated-addons.git
    targetRevision: 0.0.1
    helm:
      values: |
        repoURL: https://github.com/nearform/k8s-kurated-addons.git
        subChartsRevision: 0.0.1
        apps:
          dex:
            excluded: false
          istio-base:
            excluded: false
          istiod:
            excluded: false
          istio-ingressgateway:
            excluded: false
          knative-operator:
            excluded: false
          knative-serving:
            excluded: false
          knative-eventing:
            excluded: false
          kube-prometheus-stack:
            excluded: false
```
Characteristics:
* `app-of-apps` source points to k8s-kurated-addons repo: https://github.com/nearform/k8s-kurated-addons.git
* On the defined path (app-of-apps) there is a helm chart responsible for deploying multiple application resources.
* Values file for `app-of-apps` is located on the same path in the `values.yaml` file, and it's applied by default (no need for referencing it)
* Inline values (values defined in the body of app-of-apps manifest) define which addons are deployed and are base for making different addon flavors
* Inline values have priority over the values defined in `values.yaml`
* Disable/enable of specific addon is done by changing `app-of-apps` manifest and redeploying it using kubectl
