SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset

.PHONY: help
help: ## Display list of all targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: all
all: create-infra wait-infra create-cluster wait-cluster ## Create AWS infrastructure and k3s cluster

.PHONY: build
build: ## Compile the program into a static go binary
	hack/build.sh

.PHONY: clean
clean: ## Delete the build directory containing compiled binaries
	hack/clean.sh

.PHONY: connect
connect: ## Connect to ec2 instance via SSH
	hack/connect.sh

.PHONY: create-cluster
create-cluster: ## Create k3s cluster on the ec2 instance
	hack/create-cluster.sh

.PHONY: create-infra
create-infra: ## Create ec2 instance, security group, ssh keypair
	hack/create-infra.sh

.PHONY: destroy
destroy: ## Tear down AWS infrastructure
	hack/destroy.sh

.PHONY: wait-cluster
wait-cluster: ## Wait for the cluster to be ready
	hack/wait-cluster.sh

.PHONY: wait-infra
wait-infra: ## Wait for the ec2 instance to be ready
	hack/wait-infra.sh