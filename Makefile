# todo: we don't really care about all this buyndle nonsense since we are deploying with helm rather than
# dealing with this OLM stuff

# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= $(shell hack/version.sh)
$(info using tag '${VERSION}')

# todo: double check that we have all the necessary kubebuilder recipes

# SOURCES is the list of source files for the project
SOURCES=go.mod go.sum $(shell find . 	-type f -name '*.go')

GO_BUILD=go build
GO_FMT=go fmt
GO_TEST=go test

ifdef VERBOSE
	GO_BUILD += -v
	GO_FMT += -x
	GO_TEST += -test.v

	HELM_PACKAGE += --debug

	RM += --verbose
endif

GATEWAY_IMAGE_TAG ?= joshmeranda/marina-gateway:${VERSION}
OPERATOR_IMAGE_TAG ?= joshmeranda/marina-operator:${VERSION}

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.30.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean up build artifacts
	${RM} --recursive bin crds

##@ Development

PROTOS=$(shell find gateway/api -type f -name '*.proto')
.PHONY: proto-generate
proto-generate: ## Generate grpc protobuf api code.
	protoc -I=. \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		${PROTOS}

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -coverprofile cover.out

.PHONY: lint lint-go
lint-go: ## Run golangci-lint against code.
	golangci-lint run

lint: lint-go lint-helm ## Run all linting.

##@ Build

.PHONY: marina gateway operator

marina: ${LOCALBIN}/marina ## Build marina binary.
${LOCALBIN}/marina: ./cmd/marina/main.go ${SOURCES}
	GOBIN=${GOBIN} ${GO_BUILD} -o $@ -ldflags "-X github.com/joshmeranda/marina/cmd/marina/app.Version=${VERSION}" ./cmd/marina

gateway: ${LOCALBIN}/gateway ## Build gateway binary.
${LOCALBIN}/gateway: ./cmd/gateway/main.go ${SOURCES}
	GOBIN=${GOBIN} ${GO_BUILD} -o $@ -ldflags "-X github.com/joshmeranda/marina/cmd/gateway/app.Version=${VERSION}" ./cmd/gateway

operator: ${LOCALBIN}/operator ## Build operator binary.
${LOCALBIN}/operator: ./cmd/operator/main.go ${SOURCES}
	GOBIN=${GOBIN} ${GO_BUILD} -o $@ -ldflags "-X github.com/joshmeranda/marina/cmd/operator/app.Version=${VERSION}" ./cmd/operator

##@ Docker

.PHONY: docker docker-gateway docker-operator

docker: docker-gateway docker-operator ## Build all docker images.

docker-push: docker ## Build and push all docker image
	docker push ${GATEWAY_IMAGE_TAG}
	docker push ${OPERATOR_IMAGE_TAG}

docker-operator: ## Builder operator docker image.
	docker build --file Dockerfile.operator --tag ${OPERATOR_IMAGE_TAG} .

docker-gateway: ## Builder gateway docker image.
	docker build --file Dockerfile.gateway --tag ${GATEWAY_IMAGE_TAG} .

##@ Build Dependencies

## Tool Binaries
ENVTEST ?= $(LOCALBIN)/setup-envtest
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen-$(CONTROLLER_TOOLS_VERSION)

## Tool Versions
ENVTEST_VERSION ?= release-0.18
CONTROLLER_TOOLS_VERSION ?= v0.15.0

.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)
	mv $(LOCALBIN)/controller-gen $(CONTROLLER_GEN)
