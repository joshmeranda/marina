SOURCES=go.mod go.sum $(shell find pkg -type f -name '*.go')

# # # # # # # # # # # # # # # # # # # #
# Go commands                         #
# # # # # # # # # # # # # # # # # # # #
GO_BUILD=go build -ldflags "-X main.Version=${VERSION}"
GO_FMT=go fmt
GO_TEST=go test

ifdef VERBOSE
	GO_BUILD += -v
	GO_FMT += -x
	GO_TEST += -test.v

	HELM_PACKAGE += --debug

	RM += --verbose
endif

VERSION=$(shell hack/version.sh)

$(info using tag '${VERSION}')

LOCALBIN=$(shell pwd)/bin

# # # # # # # # # # # # # # # # # # # #
# Help text for easier Makefile usage #
# # # # # # # # # # # # # # # # # # # #
.PHONY: help

help:
	@echo "Usage: make [TARGETS]... [VALUES]"
	@echo ""
	@echo "Targets:"
	@echo "  marina             build marina binary"
	@echo "  marina-server      build marina-server binary"
	@echo "  docker             build docker image"
	@echo "  generate           run code generation"
	@echo "  lint               run linting (can run seperate linitng with lint-go and lint-helm)"
	@echo "  clean              clean built and generated files"
	@echo ""
	@echo "Values:"
	@echo "  VERBOSE            if set, various recipes are run with verbose output"
	@echo "  PUSH               if set, run docker push after building"

# # # # # # # # # # # # # # # # # # # #
# Binary building                     #
# # # # # # # # # # # # # # # # # # # #

.PHONY: marina marina-server

marina: ${LOCALBIN}/marina
${LOCALBIN}/marina: ./cmd/marina/main.go ${SOURCES}
	${GO_BUILD} -o $@ ./cmd/marina

marina-server: ${LOCALBIN}/marina-server

${LOCALBIN}/marina-server: ./cmd/marina-server/main.go ${SOURCES}
	${GO_BUILD} -o $@ ./cmd/marina-server

docker:
	docker build --tag joshmeranda/marina:${VERSION} .
	[ -n "${PUSH}" ] && docker push joshmeranda/marina:${VERSION} || true

# # # # # # # # # # # # # # # # # # # #
# code generation                     #
# # # # # # # # # # # # # # # # # # # #
PROTOS=$(shell find pkg/apis -type f -name '*.proto')

.PHONY: generate

# we have to generate everything eevry time because there is no good way to handle protobuf dependencies.
# PROTOS_GO=$(PROTOS_RAW:.proto=.pb.go)
# PROTOS_GRPC_GO=$(PROTOS_RAW:.proto=_grpc.pb.go)
#
# %.pb.go: %.proto
# 	protoc --go_out=. --go_opt=paths=source_relative $<
#
# %_grpc.pb.go: %.proto
# 	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative $<
#
# generate: ${PROTOS_GO} ${PROTOS_GRPC_GO}

generate:
	protoc -I=. \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		${PROTOS}


# # # # # # # # # # # # # # # # # # # #
# Linting recipes                     #
# # # # # # # # # # # # # # # # # # # #

.PHONY: lint lint-go lint-helm

lint-go:
	go vet ./...
	golangci-lint run

lint-helm:
	helm lint --quiet chart

lint: lint-go lint-helm

# # # # # # # # # # # # # # # # # # # #
# Test Environment recipes            #
# # # # # # # # # # # # # # # # # # # #

.PHONY: test envtest

ENVTEST=${LOCALBIN}/setup-envtest
ENV_TEST_K8S_VERSION=1.26.0

envtest: $(ENVTEST)
$(ENVTEST):
	test -s $(ENVTEST) || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

# # # # # # # # # # # # # # # # # # # #
# Testing recipes                     #
# # # # # # # # # # # # # # # # # # # #

crds:
	hack/update_crd.sh

test: crds envtest
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" ${GO_TEST} ./...

# # # # # # # # # # # # # # # # # # # #
# Project management recipes          #
# # # # # # # # # # # # # # # # # # # #

.PHONY: clean

clean:
	${RM} --recursive bin
