# GitOps Architectural Decision Record

Contents:

* [Summary](#summary)
  * [Problem](#problem)
  * [Decision](#decision)
  * [Status](#status)
* [Details](#details)
  * [Constraints](#constraints)
  * [Argument](#argument)
  * [Implications](#implications)
* [Notes](#notes)
* [Record history](#Record history)


## Summary

In initium-platform, one of the key functionalities is managing the addons. Ideally, we needed a GitOps solution for handling that. During the idealization of this repo, our engineers suggested ArgoCD as the technology for solving the problem, since it is the solution of choice for GitOps in Nearform. It also fits the GitOps principles, such as declarative definitions, versioning of the system (in our case, the addons), checking for changes automatically and drift consolidation.

### Problem

Decide which tool is going to be used for GitOps, and handling the management of the addons for initium-platform.

### Decision

ArgoCD was accepted as a great addition to initium-platform.

### Status

ArgoCD is currently being used in initium-platform.

## Details

### Constraints

* We need a technology that is declarative and supports specific revisions for each addon.
* The technology should also look out for changes and solve drifts.
* The technology should be part of Nearform's tech radar.
* The technology must be well-known by the community and with active contributors.
* The technology must be open-source.
* The software should also be easy to use and configure.
* A developer should be able to learn the addon in a few days, so initium-platform won't be an overhead to the implementation.

### Argument

ArgoCD was deployed and tested on a few developers' dev environments.
The process of installing and configuring ArgoCD, as well as running an addon was easy - less than a day of work.
ArgoCD fits the constraints mentioned above.

### Implications

ArgoCD is the most opinionated part of the solution and it is nested into further assumptions.

## Notes

n/a

## Record history
11/11/2022 - Problem defined and decisions recorded
