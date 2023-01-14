SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset

.PHONY: all
all: clean build unit-test up

.PHONY: help
help: ## Display list of all targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: ## Compile the binary
	hack/build.sh

.PHONY: clean
clean: ## Delete compiled binary from root directory
	hack/clean.sh

.PHONY: connect
connect: ## Connect to ec2 instance via SSH
	hack/connect.sh

.PHONY: down
down: ## Teardown cluster
	./ec2-k3s down

.PHONY: unit-test
unit-test: ## Run unit tests
	go test -failfast  -v ./...

.PHONY: up
up: ## Provision cluster
	./ec2-k3s up --region us-east-1
