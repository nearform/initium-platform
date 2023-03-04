# Service account for CI Bot

There are various ways to integrate a CI system with kubernetes clusters. OpenID Connect is one option, cloud providers offer 
integration points with their identity and access management systems, etc. For this project however, we aim to be agnostic 
to these services and selected a solution that would be applicable to all by creating a k8s `ServiceAccount` and `ServiceAccountToken`.  
This is an optional configuration for a cluster and not the only valid solution to solve authentication and 
authorization patterns for a CI system.

This cluster can deploy a `ServiceAccount` with the sole purpose of being utilized by the 
CI System to interact with the cluster.

The `ServiceAccount` is associated with a `kubernetes.io/service-account-token`.
Once this token is generated, the cluster administrator must export this token to the
CI System as a secret.  Once done, the pipeline for the CI system can use this token to authenticate 
to the cluster.

In the case of GitHub, follow the [official guide](https://docs.github.com/en/actions/security-guides/encrypted-secrets#about-encrypted-secrets)
for creating a secret. Depending on your specific needs the secret could be applicable to a single repo,
an environment, or the entire GitHub organization.

### Deploy the Service Account

There is a Makefile Goal that can be run against the cluster to create the aforementioned `ServiceAccount`, token,
and permission sets. To run this, simply execute the following command:
```bash
$ make create-ci-service-account
```

The permission sets should be modified for your use case.