# read PKG_VERSION from VERSION file
include VERSION

MAKEFLAGS += --no-print-directory

SPEC_FILE = spec.yaml

OTEL_VERSION = 0.113.0
OTEL_CONTRIB_VERSION = 0.113.0

# if GOOS is not supplied, set default value based on user's system, will be overridden for OS specific packaging commands
GOOS ?= $(shell go env GOOS)

ALL_SRC := $(shell find . -name '*.go' -type f | sort)
ALL_DOC := $(shell find . \( -name "*.md" -o -name "*.yaml" \) -type f | sort)
GIT_SHA := $(shell git rev-parse --short HEAD)

BUILD_INFO_IMPORT_PATH := github.com/GoogleCloudPlatform/run-gmp-sidecar/collector/internal/version
ENTRYPOINT_BUILD_INFO_IMPORT_PATH := github.com/GoogleCloudPlatform/run-gmp-sidecar/confgenerator
BUILD_X1 := -X $(BUILD_INFO_IMPORT_PATH).GitHash=$(GIT_SHA)
BUILD_X2 := -X $(BUILD_INFO_IMPORT_PATH).Version=$(PKG_VERSION)
BUILD_X3 := -X $(ENTRYPOINT_BUILD_INFO_IMPORT_PATH).Version=$(PKG_VERSION)
LD_FLAGS := -ldflags "${BUILD_X1} ${BUILD_X2} ${BUILD_X3}"

# Needs to be an absolute path.
TOOLS_DIR := $(PWD)/collector/internal/tools

.EXPORT_ALL_VARIABLES:

.DEFAULT_GOAL := presubmit

# --------------------------
#  Helper Commands
# --------------------------

.PHONY: update-components-old
update-components-old:
	grep -o github.com/open-telemetry/opentelemetry-collector-contrib/[[:lower:]]*/[[:lower:]]* go.mod | xargs -I '{}' go get {}
	go mod tidy
	cd $(TOOLS_DIR) && go get -u github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen
	cd $(TOOLS_DIR) && go mod tidy

OTEL_VER ?= latest
.PHONY: update-components
update-components:
	go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all | \
		grep "^go.opentelemetry.io" | \
		xargs -t -I '{}' go get {}@$(OTEL_VER)
	go get -u github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus@latest
	go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all | \
		grep "^github.com/open-telemetry/opentelemetry-collector-contrib" | \
		xargs -t -I '{}' go get {}@$(OTEL_VER)
	go mod tidy
	cd $(TOOLS_DIR) && go get -u go.opentelemetry.io/collector/cmd/mdatagen@$(OTEL_VER)
	cd $(TOOLS_DIR) && go mod tidy

# We can bring this target back when https://github.com/open-telemetry/opentelemetry-collector/issues/8063 is resolved.
update-opentelemetry:
	$(MAKE) update-components
	$(MAKE) install-tools
	$(MAKE) GO_BUILD_TAGS=gpu generate

# --------------------------
#  Tools
# --------------------------

DISTROGEN_BIN ?= $(TOOLS_DIR)/distrogen
MDATAGEN_BIN ?= $(TOOLS_DIR)/mdatagen

.PHONY: install-tools
install-tools:
	cd $(TOOLS_DIR) && \
		go install \
			github.com/client9/misspell/cmd/misspell \
			github.com/golangci/golangci-lint/cmd/golangci-lint \
			github.com/google/addlicense \
			go.opentelemetry.io/collector/cmd/mdatagen \
			golang.org/x/tools/cmd/goimports

.PHONY: addlicense
addlicense:
	addlicense -c "Google LLC" -l apache ./**/*.go

.PHONY: checklicense
checklicense:
	@output=`addlicense -check ./**/*.go` && echo checklicense finished successfully || (echo checklicense errors: $$output && exit 1)

.PHONY: lint
lint:
	golangci-lint run --allow-parallel-runners --build-tags=$(GO_BUILD_TAGS) --timeout=20m

.PHONY: misspell
misspell:
	@output=`misspell -error $(ALL_DOC)` && echo misspell finished successfully || (echo misspell errors:\\n$$output && exit 1)

$(DISTROGEN_BIN):
	$(MAKE) tools-dir
	GOBIN=$(TOOLS_DIR) go install github.com/GoogleCloudPlatform/opentelemetry-operations-collector/cmd/distrogen@$(shell cat .distrogen/VERSION)

$(MDATAGEN_BIN):
	$(MAKE) tools-dir
	GOBIN=$(TOOLS_DIR) bash ./scripts/download_mdatagen.sh v$(OTEL_VERSION)

# This is a PHONY target cause if you make it as a normal recipe
# it gets very confused because the creation date of the .tools
# directory is newer than the tools inside it.
.PHONY: tools-dir
tools-dir:
	@mkdir -p $(TOOLS_DIR)

.PHONY: distrogen
distrogen: $(DISTROGEN_BIN)

# --------------------------
#  CI
# --------------------------

# Adds license headers to files that are missing it, quiet tests
# so full output is visible at a glance.
.PHONY: precommit
precommit: addlicense lint misspell test

# Checks for the presence of required license headers, runs verbose
# tests for complete information in CI job.
.PHONY: presubmit
presubmit: checklicense lint misspell test_verbose

# --------------------------
#  Build and Test
# --------------------------

GO_BUILD_OUT ?= ./bin/rungmpcol
.PHONY: build-collector
build-collector:
	CGO_ENABLED=0 go build -tags=$(GO_BUILD_TAGS) -o $(GO_BUILD_OUT) $(LD_FLAGS) -buildvcs=false ./collector/cmd/rungmpcol

OTELCOL_BINARY = google-cloud-run-gmp-sidecar-$(GOOS)
.PHONY: build-collector-full-name
build-collector-full-name:
	$(MAKE) GO_BUILD_OUT=./bin/$(OTELCOL_BINARY) build-collector

ENTRYPOINT_BINARY = run-gmp-entrypoint
.PHONY: build-run-gmp-entrypoint
build-run-gmp-entrypoint:
	CGO_ENABLED=0 go build -tags=$(GO_BUILD_TAGS) -o ./bin/$(ENTRYPOINT_BINARY) $(LD_FLAGS) -buildvcs=false entrypoint.go

.PHONY: build
build:
	$(MAKE) build-collector
	$(MAKE) build-run-gmp-entrypoint

.PHONY: test-collector
test:
	$(MAKE) build-collector
	go test -tags=$(GO_BUILD_TAGS) $(GO_TEST_VERBOSE) -p 1 -race ./...

.PHONY: test_quiet
test_verbose:
	$(MAKE) GO_TEST_VERBOSE=-v test

.PHONY: test_update
test_update:
	go test ./confgenerator -update

.PHONY: go-generate
go-generate:
	go generate ./...

###################
# Distro Generation
###################

GEN_OTEL = $(DISTROGEN_BIN) generate --spec $(SPEC_FILE) \
								 --registry ./components/registry.yaml \
								 --templates ./templates

.PHONY: gen
gen: distrogen
	@$(GEN_OTEL)

.PHONY: regen
regen: distrogen
	@$(GEN_OTEL) --force

regen-v: distrogen
	@$(GEN_OTEL) --force -v


# --------------------
#  Docker
# --------------------

# set default docker build image name
BUILD_IMAGE_NAME ?= rungmpcol-build
BUILD_IMAGE_REPO ?= gcr.io/stackdriver-test-143416/opentelemetry-operations-collector:test

.PHONY: docker-build-image
docker-build-image:
	docker build -t $(BUILD_IMAGE_NAME) .

.PHONY: docker-push-image
docker-push-image:
	docker tag $(BUILD_IMAGE_NAME) $(BUILD_IMAGE_REPO)
	docker push $(BUILD_IMAGE_REPO)

.PHONY: docker-build-and-push
docker-build-and-push: docker-build-image docker-push-image

# Usage:   make TARGET=<target> docker-run
# Example: make TARGET=build-collector docker-run
TARGET ?= build_collector
.PHONY: docker-run
docker-run:
	docker run -e PKG_VERSION -v $(CURDIR):/mnt -w /mnt $(BUILD_IMAGE_NAME) /bin/bash -c "make $(TARGET)"
