# k8s-kurated-addons addon list

This document is a list of addons, what they are, how to use them and their purpose in our repository. This is going to be updated as the repository grows.

It is important to emphasize that none of the following addons are strictly **required**. That's why most of them can be disabled by adding the `enabled: false` to the app-of-apps `values.yaml` file.

## Summary
- ArgoCD
- cert-manager
- Dex
- Istio
- Knative
- kube-prometheus-stack
- Additional Notes

### ArgoCD

ArgoCD is a declarative, GitOps continuous delivery tool for Kubernetes.

ArgoCD follows the GitOps pattern of using Git repositories as the source of truth for defining the desired application state. Kubernetes manifests can be specified in several ways:

- `kustomize` applications
- `helm` charts
- `jsonnet` files
- Plain directory of YAML/json manifests
- Any custom config management tool configured as a config management plugin

ArgoCD automates the deployment of the desired application states in the specified target environments. Application deployments can track updates to branches, tags, or pinned to a specific version of manifests at a Git commit. See tracking strategies for additional details about the different tracking strategies available.

We use ArgoCD in our repository for managing all the addons that will be installed on the Kubernetes clusters. It is possible to run all the addons on the same `k8s-kurated-addons` revision, or pass down a specific revision to each addon, using the `app-of-apps/values.yaml` file `targetRevision` field.

More information at [ArgoCD Docs](https://argo-cd.readthedocs.io/en/stable/).

### cert-manager

*cert-manager is disabled by default*

`cert-manager` is a cloud-native certificate management solution designed to work on Kubernetes. It integrates with AWS Certificate Manager, GCP Certificate Manager, CloudFlare, Let's Encrypt, as well as local issuers and other providers to create SSL/TLS certificates. It is a member of CNCF since 2020.

`cert-manager` main responsibilities are to issue certificates and ensure they are valid and up to date, as well as attempt to renew them at a configured time before expiry.

We use `cert-manager` in this repository for managing all the SSL/TLS certificates a Kubernetes cluster might need. It is listed on our addon dictionary because most clusters need working SSL/TLS certificates for their services that are exposed to the internet.

The way `cert-manager` is set up in this repository, getting it to work once installed is just a matter of setting up a ClusterIssuer custom resource that will integrate with the desired provider (Let's Encrypt, for example), and configure secrets and desired domains.

More info at [cert-manager Docs](https://cert-manager.io/docs/).

### dex

Dex is a Federated OpenID Connect Provider, and a Sandbox project at CNCF.

Dex acts as a portal to other identity providers through “connectors.” This lets Dex defer authentication to LDAP servers, SAML providers, or established identity providers like GitHub, Google, and Active Directory. Clients write their authentication logic once to talk to Dex, then Dex handles the protocols for a given backend.

Once the user has dex up and running, the next step is to write applications that use dex to drive authentication. Apps that interact with dex generally fall into one of two categories:

- Apps that request OpenID Connect ID tokens to authenticate users.
    - Used for authenticating an end user.
    - Must be web based.
    - Standard OAuth2 clients. Users show up at a website, and the application wants to authenticate those end users by pulling claims out of the ID token.
- Apps that consume ID tokens from other apps.
    - Needs to verify that a client is acting on behalf of a user.
    - These consume ID tokens as credentials.
    - This lets another service handle OAuth2 flows, then use the ID token retrieved from dex to act on the end user’s behalf with the app.
    - An example of an app that falls into this category is the Kubernetes API server .

More information at [Dex Docs](https://dexidp.io/docs/getting-started/).

### Istio

Istio is an open source service mesh, which is a dedicated infrastructure layer that you can add to your applications. It allows you to transparently add capabilities like observability, traffic management, and security, without adding them to your own code.

Istio provides:
- Secure service-to-service communication in a cluster with TLS encryption, strong identity-based authentication and authorization
- Automatic load balancing for HTTP, gRPC, WebSocket, and TCP traffic
- Fine-grained control of traffic behavior with rich routing rules, retries, failovers, and fault injection
- A pluggable policy layer and configuration API supporting access controls, rate limits and quotas
- Automatic metrics, logs, and traces for all traffic within a cluster, including cluster ingress and egress

More information at [Istio Docs](https://istio.io/latest/).

### Knative

Knative is a platform-agnostic solution for running serverless deployments. It has two main components called `Serving` and `Eventing`, which empower teams working with Kubernetes. They work together to automate and manage tasks and applications.

##### Serving
Knative Serving defines a set of objects as Kubernetes Custom Resource Definitions (CRDs). These resources are used to define and control how your serverless workload behaves on the cluster.

Common use cases for Knative serving are:

- Rapid deployment of serverless containers.
- Autoscaling, including scaling pods down to zero.
- Support for multiple networking layers, such as Contour, Kourier, and Istio, for integration into existing environments.

The primary Knative Serving resources are:
- Services, which automatically manage the whole lifecycle of your workload. They control the creation of other objects to ensure that your app has a route, a configuration, and a new revision for each update of the service.

- Routes, which map a network endpoint to one or more revisions.

- Configurations, which maintain the desired state for your deployment. It provides a clean separation between code and configuration and follows the Twelve-Factor App methodology. Modifying a configuration creates a new revision.

- Revisions, which is a point-in-time snapshot of the code and configuration for each modification made to the workload. Revisions are immutable objects and can be retained for as long as useful. Knative Serving Revisions can be automatically scaled up and down according to incoming traffic.

##### Eventing
Knative Eventing is a collection of APIs that enable you to use an event-driven architecture with your applications. You can use these APIs to create components that route events from event producers to event consumers, known as sinks, that receive events. Sinks can also be configured to respond to HTTP requests by sending a response event.

Knative Eventing uses standard HTTP POST requests to send and receive events between event producers and sinks. These events conform to the CloudEvents specifications, which enables creating, parsing, sending, and receiving events in any programming language.

Common use cases of Knative Eventing are:

- Publishing an event without creating a consumer.
    - You can send events to a broker as an HTTP POST, and use binding to decouple the destination configuration from your application that produces events.

- Consuming an event without creating a publisher.
    - You can use a trigger to consume events from a broker based on event attributes. The application receives events as an HTTP POST.

More information at [Knative Docs](https://knative.dev/docs/).

### kube-prometheus-stack

> **IMPORTANT:** This addon requires >= ArgoCD 2.5.x

`kube-prometheus-stack` is a collection of Kubernetes manifests, Grafana dashboards, and Prometheus rules combined with documentation and scripts to provide easy to operate end-to-end Kubernetes cluster monitoring with Prometheus using the Prometheus Operator.

We use `kube-prometheus-stack` as the main observability stack deployed on the Kubernetes cluster. It can also be tweaked with values like Grafana login credentials, and Prometheus rules, as well as ingress configurations.

More information at [kube-prometheus-stack Docs](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack).

#### ArgoCD < 2.5.x

This addon requires server deployment, which is unavailable until ArgoCD 2.5 ( see https://github.com/argoproj/argo-cd/issues/820 ). To disable this addon, you can use the snippet below in app-of-apps `values.yaml` file:

```yaml
  - name: kube-prometheus-stack
    enabled: false
```

### Additional Notes

We are constantly evaluating new addons that might become standards in the industry. That's not high priority, though, since our main goal is to keep this repository straight to the point and minimize overhead on the users' clusters.

If you want to contribute with the repo, see [CONTRIBUTING.md](CONTRIBUTING.md).
