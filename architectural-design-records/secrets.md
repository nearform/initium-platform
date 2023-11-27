## Initium secrets

**Management of secrets and namespaces in the Initium project**

### Decision required

In the Initium project we reached a landmark where the question of how custom / user provided secrets get injected to the dynamically namespaced Knative services cannot be avoided anymore as we found that the only way to pass image pull secrets to the app is through the app’s own namespace. This in itself poses an issue right now as the CLI workflows implicitly and dynamically create / destroy separate namespaces based on the new feature's PR branch, without handling any secret creation - replication.

### Current Status

Currently we use a partitioned namespace approach to enable distinguishing the created Knative service instances through the namespace domain segment of the service instance FQDN (e.g. [appname.branchname.example.com](http://appname.branchname.example.com)). This makes app deployments bound to separate namespaces, which was a feasible choice for the seamless cleanup of the namespace after a PR merge (and automatic branch deletion).

On the other hand there’s no current solution for injecting secrets to app deployments, which is a much needed feature. To resolve the dilemma we have a couple of options to consider.

### Options

#### Single namespace approach

![Single namespace](https://raw.githubusercontent.com/nearform/initium-platform/main/architectural-design-records/img/single-namespace.jpg)

One namespace for all deployments, all secrets shared between branches.
##### Pros
* Simplified workflow - just create all in one namespace
##### Cons
* Need to refactor current deployment flow and app routing
* Namespace cluttered with deployments
* Question of secret versioning

#### Secret replication

![Secret replication](https://raw.githubusercontent.com/nearform/initium-platform/main/architectural-design-records/img/multi-namespace.jpg)

Sync secrets from a single namespace. This can be any open source replicator or ClusterSecret implementation.
##### Pros
* Keep current simple workflows and routing
* Structured / isolated deployments
##### Cons
*  Added complexity of new replicator component
*  Need to ensure replication is finished before Knative service launch

### Recommendation

The case for picking a solution that avoids the need to make any fundamental changes to the current Initium workflows is strong, since apart from secrets, all other aspects of the solution fit well in the namespace separated design - as it’s indicated in the pros / cons section, having a single namespace draws a lot of concerns and shortcomings which would need to be addressed later down the line.

### Decision

Picking a custom open source operator to be part of a toolchain is not easy, but the best bet would be to pick one that brings a new resource type like ClusterSecret and watches live namespace events to sync secrets immediately. For this, [https://clustersecret.io/](https://clustersecret.io/) is a good option as it’s the most lightweight of the alternatives and gives the exact functionality we need. Once it’s deployed, administrators can setup the required ClusterSecret resources as part of the Knative application’s main namespace and replicate the secrets from there to the branch based, newly created ones.

### Next steps

-   Add ClusterSecret operator to the list of apps
-   Add instructions on how to setup arbitrary secrets (like private docker registry credentials) with ClusterSecret definitions
-   Test the new component with docker registry secrets and close down [https://github.com/nearform/initium-cli/issues/84](https://github.com/nearform/initium-cli/issues/84) after merging [https://github.com/nearform/initium-cli/pull/94](https://github.com/nearform/initium-cli/pull/94).
