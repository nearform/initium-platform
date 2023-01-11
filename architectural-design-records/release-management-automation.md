## Summary

Use a Release-please GitHub action to automatically release new versions of our software on GitHub.

### Problem

Releasing new versions of software can be a time-consuming task, especially for teams that release software frequently. It can be easy to miss steps or make errors during the release process, and releasing software manually can take up valuable development time.

### Decision

We will use a Release-please GitHub action to automatically release new versions of our software on GitHub. The action will be triggered by a merge to main branch or specific commit tag, it will check all the CI checks pass, update the version of package, create a release and push it to the appropriate git tag and branch. Additionally, there'll be an action that will upload the app-of-apps.yaml to the release assets.

### Status

In-use

## Details

### Constraints

* The action should be configured to respect the semantic versioning of the packages.
* Should be able to define the release branch, where it should push the release.
* Must be able to define the git tag format.
* Should be able to validate the release notes.
* Should be able to push the release notes to GitHub release notes.

### Argument

Using a Release-please GitHub action will save time and reduce errors by automating the process of releasing new versions of software. Additionally, the action can be configured to follow best practices for releasing software, such as respecting semantic versioning, validating release notes, and pushing the release notes to GitHub releases.

### Implications

* The team will need to make sure that the git tag and branch are correctly defined in the workflow configuration.
* New versions of software will be released automatically, and the team can focus on other tasks.
* The action might break the release process if not configured properly

## Notes

* The Release-please GitHub action can be configured to use different release strategies, such as major, minor, and patch releases.
* The action can also be configured to add "release" labels to pull requests, making it easy to identify releases that are waiting to be deployed.

## Record history
11/03/2023 - Problem defined and decisions recorded
