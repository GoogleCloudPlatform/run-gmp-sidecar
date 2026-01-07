# read PKG_VERSION from VERSION file
include VERSION

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

TOOLS_DIR := collector/internal/tools

.EXPORT_ALL_VARIABLES:

.DEFAULT_GOAL := presubmit

# --------------------------
#  Helper Commands
# --------------------------

# The non-stable OpenTelemetry Collector version.
OTEL_VERSION = v0.113.0
# The OpenTelemetry Collector contrib repo version.
# Equal to OTEL_VERSION most of the time, only deviates in rare cases.
OTEL_CONTRIB_VERSION = v0.113.0

DISTROGEN_BIN ?= distrogen
LIST_DIRECT_MODULES = go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all
INCLUDE_COLLECTOR_CORE_COMPONENTS = grep "^go.opentelemetry.io" | grep -v "^go.opentelemetry.io/otel"
INCLUDE_CONTRIB_COMPONENTS = grep "^github.com/open-telemetry/opentelemetry-collector-contrib"
GO_GET_ALL = xargs --no-run-if-empty -t -I '{}' go get {}

.PHONY: update-components
update-components: install-tools core-components contrib-components

.PHONY: core-components
core-components:
	$(LIST_DIRECT_MODULES) | \
		$(INCLUDE_COLLECTOR_CORE_COMPONENTS) | \
		$(DISTROGEN_BIN) otel_component_versions --otel_version $(OTEL_VERSION) | \
		$(GO_GET_ALL)

.PHONY: contrib-components
contrib-components:
	$(LIST_DIRECT_MODULES) | \
		$(INCLUDE_CONTRIB_COMPONENTS) | \
		$(GO_GET_ALL)@$(OTEL_CONTRIB_VERSION)

.PHONY: update-mdatagen
update-mdatagen:
	cd $(TOOLS_DIR) && \
		go get -u go.opentelemetry.io/collector/cmd/mdatagen@$(OTEL_VERSION) && \
		go mod tidy

.PHONY: update-opentelemetry
update-opentelemetry: update-components update-mdatagen generate

# --------------------------
#  Tools
# --------------------------

.PHONY: install-tools
install-tools:
	cd $(TOOLS_DIR) && \
		go install \
			github.com/client9/misspell/cmd/misspell \
			github.com/golangci/golangci-lint/cmd/golangci-lint \
			github.com/google/addlicense \
			go.opentelemetry.io/collector/cmd/mdatagen \
			golang.org/x/tools/cmd/goimports \
			github.com/GoogleCloudPlatform/opentelemetry-operations-collector/cmd/distrogen

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

.PHONY: test
test:
	$(MAKE) build-collector
	go test -tags=$(GO_BUILD_TAGS) $(GO_TEST_VERBOSE) -p 1 -race ./...

.PHONY: test_quiet
test_verbose:
	$(MAKE) GO_TEST_VERBOSE=-v test

.PHONY: test_update
test_update:
	go test ./confgenerator -update

.PHONY: generate
generate:
	go generate ./...

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
