SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset

.PHONY: help
help: ## Display list of all targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: all
all: create connect ## Create AWS infrastructure and ssh to ec2 instance

.PHONY: create
create: ## Create ec2 instance, security group, ssh keypair, and k3d cluster
	hack/create.sh

.PHONY: connect
connect: ## ssh to the ec2 instance
	hack/connect.sh

.PHONY: destroy
destroy: ## Tear down AWS infrastructure
	hack/destroy.sh