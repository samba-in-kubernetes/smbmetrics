# Allow developer to override some defaults
-include devel.mk

# Version
VERSION?=0.0.1

# Commit info
COMMIT_ID=$(shell git describe --abbrev=40 --always --exclude='*' --dirty=+ 2>/dev/null)
GIT_VERSION=$(shell git describe --match='v[0-9]*.[0-9]' --match='v[0-9]*.[0-9].[0-9]' 2>/dev/null || echo "(unset)")

GO_CMD:=go
GOFMT_CMD:=gofmt
BUILDAH_CMD:=buildah
YAMLLINT_CMD:=yamllint

# Image URL to use all building/pushing image targets
TAG?=latest
IMG?=quay.io/samba.org/samba-metrics:$(TAG)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell $(GO_CMD) env GOBIN))
GOBIN=$(shell $(GO_CMD) env GOPATH)/bin
else
GOBIN=$(shell $(GO_CMD) env GOBIN)
endif

# Get current GOARCH
GOARCH?=$(shell $(GO_CMD) env GOARCH)

# Local (alternative) GOBIN for auxiliary build tools
GOBIN_ALT:=$(CURDIR)/.bin

# Common link-flags for go programs
GOLDFLAGS="-X main.Version=$(GIT_VERSION) -X main.CommitID=$(COMMIT_ID)"


CONTAINER_BUILD_OPTS?=
CONTAINER_CMD?=
ifeq ($(CONTAINER_CMD),)
	CONTAINER_CMD:=$(shell docker version >/dev/null 2>&1 && echo docker)
endif
ifeq ($(CONTAINER_CMD),)
	CONTAINER_CMD:=$(shell podman version >/dev/null 2>&1 && echo podman)
endif
# handle the case where podman is present but is (defaulting) to remote and is
# not not functioning correctly. Example: mac platform but not 'podman machine'
# vms are ready
ifeq ($(CONTAINER_CMD),)
	CONTAINER_CMD:=$(shell podman --version >/dev/null 2>&1 && echo podman)
ifneq ($(CONTAINER_CMD),)
$(warning podman detected but 'podman version' failed. \
	this may mean your podman is set up for remote use, but is not working)
endif
endif

# Helper function to re-format yamls using helper script
define yamls_reformat
	YQ=$(YQ) $(CURDIR)/hack/yq-fixup-yamls.sh $(1)
endef

all: build

# Build executable
.PHONY: build
build:
	CGO_ENABLED=0 $(GO_CMD) build -o bin/smbmetrics \
		-ldflags $(GOLDFLAGS) cmd/main.go

# Run unit tests
.PHONY: test
test: build vet
	$(GO_CMD) test ./... -coverprofile cover.out


# Run go fmt to reformat code
.PHONY: reformat
reformat:
	$(GO_CMD) fmt ./...

# Run go vet against code
.PHONY: vet
vet: reformat
	$(GO_CMD) vet ./...

# Format yaml files for yamllint standard
.PHONY: yaml-fmt
yaml-fmt: yq
	$(call yamls_reformat, $(CURDIR))

# Build the container image
.PHONY: image-build
image-build: Dockerfile
	$(CONTAINER_CMD) build \
		--build-arg=GIT_VERSION="$(GIT_VERSION)" \
		--build-arg=COMMIT_ID="$(COMMIT_ID)" \
		--build-arg=ARCH="$(GOARCH)" \
		$(CONTAINER_BUILD_OPTS) -f $< -t $(IMG) .

# Push the container image
.PHONY: image-push
image-push:
	$(CONTAINER_CMD) push $(IMG)

# Check the code
.PHONY: check check-golangci-lint check-format check-yaml

check: check-golangci-lint vet check-yaml

check-golangci-lint: golangci-lint
	$(GOLANGCI_LINT) -c .golangci.yaml run ./...

check-yaml:
	$(YAMLLINT_CMD) -c ./.yamllint.yaml ./

# Find or download auxiliary build tools
.PHONY: build-tools golangci-lint yq
build-tools: golangci-lint yq

define installtool
	@GOBIN=$(GOBIN_ALT) GO_CMD=$(GO_CMD) $(CURDIR)/hack/install-tools.sh $(1)
endef


GOLANGCI_LINT=$(GOBIN_ALT)/golangci-lint
golangci-lint:
ifeq (, $(shell command -v $(GOLANGCI_LINT) ;))
	@$(call installtool, --golangci-lint)
	@echo "golangci-lint installed in $(GOBIN_ALT)"
endif

YQ=$(GOBIN_ALT)/yq
yq:
ifeq (, $(shell command -v $(YQ) ;))
	@$(call installtool, --yq)
	@echo "yq installed in $(YQ)"
endif


