## Summary

Use a Renovate bot to automatically update dependencies in our GitHub repository.

### Problem

Keeping dependencies up to date can be a time-consuming task, especially for large projects with many dependencies. It can be easy to miss updates, and updating dependencies manually can introduce errors and inconsistencies. As the number of charts that the project supports will grow,
keeping up-to-date with the latest releases manually is a time-consuming effort.

### Decision

We will use a Renovate bot to automatically update dependencies in our GitHub repository. The bot will check for updates regularly and open pull requests for any updates that are found and merge automatically if the checks will pass and the change is backward compatible (patch or minor). The pull requests with major changes will then be reviewed and merged by the project maintainers.

### Status

In-use

## Details

### Constraints

* The bot should be configured to respect the semantic versioning of the packages.
* Should not update packages that are not compatible with the current version.
* Must be able to configure the schedule of updates.
* Should be able to lock certain package versions.

### Argument

Using a Renovate bot will save time and reduce errors by automating the process of updating dependencies. Additionally, the bot can be configured to follow best practices for dependency management, such as respecting semantic versioning and not updating packages that are not compatible with the current version.

### Implications

* The team of maintainers will need to review and merge pull requests created by the bot regularly.
* Dependencies will be kept up to date and the team can focus on other tasks.
* The bot can break the current build if not configured properly

## Notes
* The Renovate bot can be configured to only update certain types of dependencies (e.g. only direct dependencies).
* Some add-ons will need additional renovate bot configuration update to group the versions update.

## Record history
11/03/2023 - Problem defined and decisions recorded
