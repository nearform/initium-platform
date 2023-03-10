###############################################################################
# ENVIRONMENT CONFIGURATION
###############################################################################
MAKEFLAGS += --no-print-directory
SHELL=/bin/bash

.PHONY: default

###############################################################################
# GOALS ( safe defaults )
###############################################################################

default: generate-certs kind-up tilt-up

clean: tilt-down kind-down
	@./scripts/clean-test.sh

ci: generate-certs kind-up
	@./scripts/ci.sh

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

argocd:
	@./scripts/cloud.sh

###############################################################################
# GOALS ( bootstrap)
###############################################################################

plugin_install: ## install asdf plugins
	@asdf plugin add argocd https://github.com/beardix/asdf-argocd.git || true
	@asdf plugin add awscli https://github.com/MetricMike/asdf-awscli.git || true
	@asdf plugin add eksctl https://github.com/elementalvoid/asdf-eksctl.git || true
	@asdf plugin add golang https://github.com/kennyp/asdf-golang.git || true
	@asdf plugin add golangci-lint https://github.com/hypnoglow/asdf-golangci-lint.git || true
	@asdf plugin add jq https://github.com/AZMCode/asdf-jq.git || true
	@asdf plugin add helm https://github.com/Antiarchitect/asdf-helm.git || true
	@asdf plugin add kind https://github.com/reegnz/asdf-kind.git || true
	@asdf plugin add kubectl https://github.com/asdf-community/asdf-kubectl.git || true
	@asdf plugin add pre-commit https://github.com/jonathanmorley/asdf-pre-commit.git || true
	@asdf plugin add tilt https://github.com/virtualstaticvoid/asdf-tilt.git || true

plugin_uninstall: ## uninstall asdf plugins
	@asdf plugin remove argocd || true
	@asdf plugin remove awscli || true
	@asdf plugin remove eksctl || true
	@asdf plugin remove golang || true
	@asdf plugin remove golangci-lint || true
	@asdf plugin remove jq || true
	@asdf plugin remove helm || true
	@asdf plugin remove kind || true
	@asdf plugin remove kubectl || true
	@asdf plugin remove pre-commit || true
	@asdf plugin remove tilt || true

asdf_install: plugin_install ## install all plugins and packages present in .tool-versions file
	@asdf install

asdf_uninstall: plugin_uninstall ## uninstall all plugins and packages present in .tool-versions file

pre-commit: ## set up the git hook scripts
	@pre-commit install

bootstrap: asdf_install pre-commit ## setup all the needed plugins and install pre-commit hook

###############################################################################
# GOALS
###############################################################################

generate-certs: ## Generate TLS assets for to be used by cluster API server and istio ingress gateway
	@./scripts/generate-certs.sh

kind-up: ## Create the K8s cluster
	@./scripts/kind-up.sh

kind-down: ## Destroy the K8s cluster
	@./scripts/kind-down.sh

tilt-up: ## Start Tilt and deploy all the apps
	@./scripts/tilt-up.sh

tilt-down: ## Stop Tilt and destroy all the apps
	@./scripts/tilt-down.sh

tilt-validate: ## Validate the Tiltfile
	@./scripts/tilt-validate.sh

unit-test: ## Run unit tests
	@./scripts/run-test.sh unit

integration-test: ## Run integration tests
	@./scripts/run-test.sh integration examples/sample-app

integration-test-podinfo: ## Run integration tests
	@./scripts/run-test.sh integration examples/sample-podinfo-app

validate: ## Run static checks
	@ASDF_DEFAULT_TOOL_VERSIONS_FILENAME=$(CURDIR)/.tool-versions pre-commit run --color=always --show-diff-on-failure --all-files

create-ci-service-account: ## Create a k8s service account that would be used by CI systems
	@./scripts/create-ci-service-account.sh
