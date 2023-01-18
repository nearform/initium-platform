# Known Issues

### Knative Deletion Race Condition

Knative Operator inserts finalizers in `knativeserving.operator.knative.dev` and `knativeeventing.operator.knative.dev`
resources. As all three ArgoCD applications (knative-serving, knative-eventing, knative-operator) are deleted at the
same time, Knative Operator is removed before it can fulfill finalizer condition when `knativeserving.operator.knative.dev`
and `knativeeventing.operator.knative.dev` resources are to be deleted. This results in knative-serving and knative-eventing
apps stuck in deletion.

Similar issue is present in the ArgoCD project: [Keycloak deletion deadlock](https://github.com/argoproj/argo-cd/issues/9296).

**Resolution**

Proposed solution is to add [post-delete hook](https://github.com/argoproj/argo-cd/issues/7575) to ArgoCD. Until this
feature is available, manual solution is offered.

After deletion action is executed, it is necessary to remove finalizers from `knativeserving.operator.knative.dev` and
`knativeeventing.operator.knative.dev` resources. This will allow application removal process to finish.
```bash
kubectl patch knativeeventing knative-eventing -n knative-eventing -p '{"metadata":{"finalizers":null}}' --type=merge
kubectl patch knativeserving knative-serving -n knative-serving -p '{"metadata":{"finalizers":null}}' --type=merge
```
