# Service account for CI Bot
This cluster will deploy a `ServiceAccount` with the 
sole purpose of being utilized by the CI System to interact with the cluster.

The `ServiceAccount` is associated with a `kubernetes.io/service-account-token`.
Once this token is generated, the cluster administrator must export this token to the
CI System as a secret.

In the case of GitHub, follow the [official guide](https://docs.github.com/en/actions/security-guides/encrypted-secrets#about-encrypted-secrets)
for creating a secret. Depending on your specific needs the secret could be applicable to a single repo, 
an environment, or the entire GitHub organization.