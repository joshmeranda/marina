# todo: we don't really care about all this buyndle nonsense since we are deploying with helm rather than
# dealing with this OLM stuff

# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= $(shell hack/version.sh)
$(info using tag '${VERSION}')

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

# IMAGE_TAG_BASE defines the docker.io namespace and part of the image name for remote images.
# This variable is used to construct full image tags for bundle and catalog images.
#
# For example, running 'make bundle-build bundle-push catalog-build catalog-push' will build and push both
# marina.io/marina-operator-bundle:$VERSION and marina.io/marina-operator-catalog:$VERSION.
# IMAGE_TAG_BASE ?= marina.io/marina-operator
SERVER_IMAGE_TAG ?= joshmeranda/marina-server:${VERSION}

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.26.0

# MARINA_OPERATOR_URL is the git repository URL for the marina-operator project.
MARINA_OPERATOR_URL=git@github.com:joshmeranda/marina-operator.git

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

PROTOS=$(shell find pkg/apis -type f -name '*.proto')
.PHOTO: proto-generate
proto-generate: ## Generate grpc protobuf api code.
	protoc -I=. \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		${PROTOS}

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: crds fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -coverprofile cover.out

.PHONY: lint lint-go
lint-go: ## Run golangci-lint against code.
	golangci-lint run

lint: lint-go lint-helm ## Run all linting.

##@ Build

.PHONY: marina marina-server

marina: ${LOCALBIN}/marina ## Build marina binary.
${LOCALBIN}/marina: ./cmd/marina/main.go ${SOURCES}
	GOBIN=${GOBIN} ${GO_BUILD} -o $@ -ldflags "-X github.com/joshmeranda/marina/cmd/marina.Version=${VERSION}" ./cmd/marina

marina-server: ${LOCALBIN}/marina-server ## Build marina-server binary.
${LOCALBIN}/marina-server: ./cmd/marina-server/main.go ${SOURCES}
	GOBIN=${GOBIN} ${GO_BUILD} -o $@ -ldflags "-X github.com/joshmeranda/marina/cmd/server.Version=${VERSION}" ./cmd/marina-server

##@ Docker

.PHONY: docker-marina-server

docker-marina-server: ## Builder dockeri mage with marina-server.
	docker build --file Dockerfile --tag ${SERVER_IMAGE_TAG} .

##@ Build Dependencies

## Tool Binaries
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

crds: ## Pull crds from the marina-operator project.
ifdef OPERATOR_VERSION
	@echo "Pulling crds from ${MARINA_OPERATOR_URL} at tag ${OPERATOR_VERSION}"
	git clone ${MARINA_OPERATOR_URL} && cd marina-operator && git checkout ${OPERATOR_VERSION} && cp -r config/crd/bases ../crds && cd .. && rm --force --recursive marina-operator
else
	@echo "Pulling crds from ${MARINA_OPERATOR_URL}"
	git clone ${MARINA_OPERATOR_URL} && cp -r marina-operator/config/crd/bases crds && rm --force --recursive marina-operator
endif