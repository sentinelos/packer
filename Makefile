SHA := $(shell git describe --match=none --always --abbrev=8 --dirty)
TAG := $(shell git describe --tag --always --dirty)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
ARTIFACTS := bin
WITH_DEBUG ?= false
WITH_RACE ?= false
REGISTRY ?= ghcr.io
VENDOR ?= sentinelos
REGISTRY_AND_VENDOR ?= $(REGISTRY)/$(VENDOR)
IMAGE ?= packer
GOLANGCILINT_VERSION ?= v1.50.1
GOFUMPT_VERSION ?= v0.4.0
GO_VERSION ?= $(shell grep -m 1 go go.mod | cut -d\  -f2)
TOOLCHAIN ?= docker.io/golang:$(GO_VERSION)-alpine3.16
GOIMPORTS_VERSION ?= v0.3.0
GO_BUILDFLAGS ?=
GO_LDFLAGS ?=
CGO_ENABLED ?= 0
TESTPKGS ?= ./...

# nerdctl build settings

BUILD := nerdctl build
PLATFORM ?= linux/amd64,linux/arm64
PROGRESS ?= auto
COMMON_ARGS = --file=Containerfile
COMMON_ARGS += --progress=$(PROGRESS)
COMMON_ARGS += --build-arg=SHA="$(SHA)"
COMMON_ARGS += --build-arg=TAG="$(TAG)"
COMMON_ARGS += --build-arg=VENDOR="$(VENDOR)"
COMMON_ARGS += --build-arg=REGISTRY="$(REGISTRY)"
COMMON_ARGS += --build-arg=TOOLCHAIN="$(TOOLCHAIN)"
COMMON_ARGS += --build-arg=CGO_ENABLED="$(CGO_ENABLED)"
COMMON_ARGS += --build-arg=GO_BUILDFLAGS="$(GO_BUILDFLAGS)"
COMMON_ARGS += --build-arg=GO_LDFLAGS="$(GO_LDFLAGS)"
COMMON_ARGS += --build-arg=GOLANGCILINT_VERSION="$(GOLANGCILINT_VERSION)"
COMMON_ARGS += --build-arg=GOFUMPT_VERSION="$(GOFUMPT_VERSION)"
COMMON_ARGS += --build-arg=GOIMPORTS_VERSION="$(GOIMPORTS_VERSION)"
COMMON_ARGS += --build-arg=TESTPKGS="$(TESTPKGS)"
CI_ARGS ?=
PUSH ?= false

# extra variables
RUN_TESTS ?= TestIntegration

# help menu

export define HELP_MENU_HEADER
# Getting Started

To build this project, you must have the following installed:

- git
- make
- nerdctl (1.0.0 or higher)

## Artifacts

All artifacts will be output to ./$(ARTIFACTS). Images will be tagged with the
registry "$(REGISTRY)", username "$(VENDOR)", and a dynamic tag (e.g. $(IMAGE):$(TAG)).
The registry and username can be overridden by exporting REGISTRY, and VENDOR
respectively.

endef

ifneq (, $(filter $(WITH_RACE), t true TRUE y yes 1))
GO_BUILDFLAGS += -race
CGO_ENABLED := 1
GO_LDFLAGS += -linkmode=external -extldflags '-static'
endif

ifneq (, $(filter $(WITH_DEBUG), t true TRUE y yes 1))
GO_BUILDFLAGS += -tags sentinelos.debug
else
GO_LDFLAGS += -s -w
endif

all: unit-tests integration.test integration lint packer packer-image

%-local: ## Builds the specified target. The build result will be output to the specified local destination.
	@$(MAKE) $*-target TARGET_ARGS="--output=type=local,dest=$(ARTIFACTS) $(TARGET_ARGS)"

%-image: ## Builds the specified target. The build result will be loaded into image.
	@$(MAKE) $*-target TARGET_ARGS="--tag=$(REGISTRY_AND_VENDOR)/$(IMAGE):$(TAG) --output type=image,name=$(REGISTRY_AND_VENDOR)/$(IMAGE):$(TAG),push=$(PUSH) $(TARGET_ARGS)"

%-target: ## Builds the specified target. The build result will only remain in the build cache.
	@$(BUILD) --target=$* $(COMMON_ARGS) $(PLATFORM_ARGS) $(TARGET_ARGS) $(CI_ARGS) .

lint-golangci-lint:  ## Runs golangci-lint linter.
	@$(MAKE) $@-target

lint-gofumpt:  ## Runs gofumpt linter.
	@$(MAKE) $@-target

.PHONY: fmt
fmt:  ## Formats the source code
	@nerdctl run --rm -it -v $(PWD):/src -w /src $(TOOLCHAIN) \
		sh -c "export GO111MODULE=on; export GOPROXY=https://proxy.golang.org; \
		go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION) && \
		go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION) && \
		gofumpt -w . && \
        goimports -w -local github.com/sentinelos/packer ."

lint-govulncheck:  ## Runs govulncheck linter.
	@$(MAKE) $@-target

lint-goimports:  ## Runs goimports linter.
	@$(MAKE) $@-target

.PHONY: base
base:  ## Prepare base toolchain
	@$(MAKE) $@-target

.PHONY: unit-tests
unit-tests:  ## Performs unit tests
	@$(MAKE) $@-local

.PHONY: unit-tests-race
unit-tests-race:  ## Performs unit tests with race detection enabled.
	@$(MAKE) $@-target

.PHONY: coverage
coverage:  ## Upload coverage data to codecov.io.
	bash -c "bash <(curl -s https://codecov.io/bash) -f $(ARTIFACTS)/coverage.txt -X fix"

.PHONY: packer-linux-amd64
packer-linux-amd64:  ## Builds executables for packer platform linux/amd64.
	@$(MAKE) packer-linux-amd64-local PLATFORM_ARGS="--platform=linux/amd64"

.PHONY: packer-linux-arm64
packer-linux-arm64:  ## Builds executables for packer platform linux/arm64.
	@$(MAKE) packer-linux-arm64-local PLATFORM_ARGS="--platform=linux/arm64"

.PHONY: packer
packer: packer-linux-amd64 packer-linux-arm64  ## Builds executables for packer.

.PHONY: lint-markdown
lint-markdown:  ## Runs markdownlint.
	@$(MAKE) $@-target

.PHONY: lint
lint: lint-golangci-lint lint-gofumpt lint-govulncheck lint-goimports lint-markdown  ## Run all linters for the project.

.PHONY: packer-image
packer-image:  ## Builds image for packer.
	@$(MAKE) $@-image PLATFORM_ARGS="--platform=$(PLATFORM)"

.PHONY: integration.test
integration.test:
	@$(MAKE) $@-local

.PHONY: integration
integration: integration.test packer
	@$(MAKE) packer-image PUSH=true
	cp $(ARTIFACTS)/packer-linux-amd64 $(ARTIFACTS)/packer
	cd pkg/integration && PATH="$$PWD/../../../$(ARTIFACTS):$$PATH" integration.test -test.v -test.run $(RUN_TESTS)

.PHONY: clean
clean:  ## Cleans up all artifacts.
	@rm -rf $(ARTIFACTS)/*

.PHONY: help
help:  ## This help menu.
	@echo "$$HELP_MENU_HEADER"
	@grep -E '^[a-zA-Z%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
