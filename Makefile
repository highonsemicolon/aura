MAKEFLAGS += --no-print-directory

print-%:
	@echo '$*=$($*)'

# ==== Project metadata ====
APP_NAME        ?= aura
BINARY_DIR      ?= .bin
BUILD_DIR       ?= .build
GO_VERSION 		?= $(shell if [ -f go.work ]; then grep '^go ' go.work | awk '{print $$2}'; else grep '^go ' services/app/go.mod | awk '{print $$2}'; fi)
GO_FLAGS        ?= -trimpath -buildvcs=false
GOBIN           ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
  GOBIN := $(shell go env GOPATH)/bin
endif

# Versions & metadata
VERSION         ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo 0.1.0)
COMMIT_HASH     ?= $(shell git rev-parse --short=12 HEAD)
BUILD_TIME      := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_HOST      := $(shell hostname)

# Go modules (top-level and nested)
MODULES         := $(shell go list -m -f '{{.Dir}}')
SERVICES        := $(shell ls -d services/* 2>/dev/null | xargs -n1 basename)

# Linker flags for versioning (expect main package to expose these variables)
# Example in Go: var (
#   Version = "dev"; Commit = ""; BuildTime = ""; BuiltBy = "make"
# )
LDFLAGS         := -X 'main.Version=$(VERSION)' \
                   -X 'main.Commit=$(COMMIT_HASH)' \
                   -X 'main.BuildTime=$(BUILD_TIME)' \
                   -X 'main.BuiltBy=$(BUILD_HOST)'

.DEFAULT_GOAL := help

help:
	@echo "Targets:"; \
	echo "  tidy              - go mod tidy across modules"; \
	echo "  fmt               - go fmt across modules"; \
	echo "  lint              - run golangci-lint"; \
	echo "  test              - unit tests"; \
	echo "  cover             - tests with coverage + threshold"; \
	echo "  build             - build all services"; \
	echo "  build-one SERVICE=app - build a single service"; \
	echo "  run SERVICE=app   - run a single service"; \
	echo "  dev SERVICE=app   - run a service with hot reload (Air)"; \
	echo "  proto             - generate protobuf via buf (fallback to protoc)"; \
	echo "  docker-build      - build docker images for all services"; \
	echo "  docker-push       - push docker images (requires REGISTRY/IMAGE_OWNER)"; \
	echo "  release           - build binaries w/ metadata";

print-go-version:
	@go version

# ---- Hygiene ----

tidy:
	@for m in $(MODULES); do \
		cd $$m && go mod tidy; \
	done

fmt:
	@for m in $(MODULES); do \
		cd $$m && go fmt ./...; \
	done

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v2.4.0; };
	@for m in $(MODULES); do \
		echo "→ golangci-lint $$m"; \
		(cd $$m && golangci-lint run ./...); \
	done

# ---- Development ----

AIR_BIN ?= $(GOBIN)/air

$(AIR_BIN):
	@echo "Installing Air..."
	@go install github.com/cosmtrek/air@latest

dev: $(AIR_BIN)
ifndef SERVICE
	$(error Usage: make dev SERVICE=app)
endif
	@SERVICE_NAME=$(SERVICE) $(AIR_BIN) -c .air.toml


# ---- Tests ----

COVER_PROFILE := $(BUILD_DIR)/coverage.out
COVER_HTML    := $(BUILD_DIR)/coverage.html
COVER_MIN     ?= 00

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

test-one:
ifndef SERVICE
	$(error Usage: make test-one SERVICE=app)
endif
	@echo "→ Running tests for $(SERVICE)"
	@if [ -d services/$(SERVICE) ]; then \
		go test -race -v ./services/$(SERVICE)/...; \
	else \
		echo "Service $(SERVICE) not found"; exit 1; \
	fi

cover-one: $(BUILD_DIR)
ifndef SERVICE
	$(error Usage: make cover-one SERVICE=app)
endif
	@echo "→ Running coverage for $(SERVICE)"
	@if [ -d services/$(SERVICE) ]; then \
		COVER_FILE=$(BUILD_DIR)/$(SERVICE).cover.out; \
		go test -coverprofile=$$COVER_FILE -covermode=atomic -race ./services/$(SERVICE)/...; \
		go tool cover -func=$$COVER_FILE | tee $(BUILD_DIR)/$(SERVICE).coverage.txt; \
	else \
		echo "Service $(SERVICE) not found"; exit 1; \
	fi

build-one:
ifndef SERVICE
	$(error Usage: make build-one SERVICE=app)
endif
	@if [ -f services/$(SERVICE)/main.go ]; then \
		echo "Building $(SERVICE)"; \
		GOOS=$${GOOS:-$(shell go env GOOS)} \
		GOARCH=$${GOARCH:-$(shell go env GOARCH)} \
		go build $(GO_FLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$(SERVICE) ./services/$(SERVICE); \
	else \
		echo "Service $(SERVICE) not found"; exit 1; \
	fi

build: $(BINARY_DIR)
	@for svc in $(SERVICES); do \
		if [ -f services/$$svc/main.go ]; then \
			echo "Building $$svc"; \
			GOOS=$${GOOS:-linux} GOARCH=$${GOARCH:-amd64} \
			go build $(GO_FLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$$svc ./services/$$svc; \
		fi; \
	done

$(BINARY_DIR):
	@mkdir -p $(BINARY_DIR)

run:
ifndef SERVICE
	$(error Usage: make run SERVICE=app)
endif
	@go run -ldflags "$(LDFLAGS)" ./services/$(SERVICE)

release: clean build
	@echo "Binaries in $(BINARY_DIR)"

clean:
	@rm -rf $(BINARY_DIR) $(BUILD_DIR)

# ---- Protobuf ----
PROTO_DIR      := apis
BUF := $(GOBIN)/buf

$(BUF):
	@echo "Installing Buf..."
	@go install github.com/bufbuild/buf/cmd/buf@latest

proto: $(BUF)
	@echo "Generating protobuf with protoc"; \
	buf generate

proto-lint: $(BUF)
	@echo "→ Linting protobuf definitions"
	buf lint

proto-breaking: $(BUF)
	@echo "→ Checking protobuf breaking changes"
	buf breaking --against '.git#branch=main,subdir=apis'

proto-all: proto-lint proto-breaking proto

# ---- Docker ----

REGISTRY     ?= ghcr.io
IMAGE_OWNER  ?= $(shell basename $(shell dirname $(shell git remote get-url origin 2>/dev/null || echo unknown/unknown)))
IMAGE_TAG    ?= $(VERSION)
DOCKER_BUILDKIT ?= 1
PLATFORMS    ?= linux/amd64
PUSH         ?= false
LOCAL ?= true

# Build all service images using a shared Dockerfile template
# Expect per-service build context at repo root with ARG SERVICE=<name>

ifeq ($(LOCAL),true)
    CACHE_FROM=
    CACHE_TO=
else
    CACHE_FROM=--cache-from=type=registry,ref=$(REGISTRY)/$(IMAGE_OWNER)/$(APP_NAME)-$(SERVICE):buildcache
    CACHE_TO=--cache-to=type=registry,ref=$(REGISTRY)/$(IMAGE_OWNER)/$(APP_NAME)-$(SERVICE):buildcache,mode=max
endif

docker-build-one:
ifndef SERVICE
	$(error Usage: make docker-build-one SERVICE=app)
endif
	@if [ -f services/$(SERVICE)/main.go ]; then \
		IMG_BASE=$(REGISTRY)/$(IMAGE_OWNER)/$(APP_NAME)-$(SERVICE); \
		IMG_VERSION=$$IMG_BASE:$(IMAGE_TAG); \
		IMG_LATEST=$$IMG_BASE:latest; \
		echo "=== Building $(SERVICE) -> $$IMG_VERSION ==="; \
		docker buildx build \
			--build-arg GO_VERSION=$(GO_VERSION) \
			--platform $(PLATFORMS) \
			--build-arg SERVICE=$(SERVICE) \
			--build-arg VERSION=$(VERSION) \
			--build-arg COMMIT=$(COMMIT_HASH) \
			--build-arg LDFLAGS="$(LDFLAGS)" \
			$(CACHE_FROM) \
    		$(CACHE_TO) \
			-t $$IMG_VERSION \
			-t $$IMG_LATEST \
			-f Dockerfile \
			$(if $(filter true,$(PUSH)),--push,) \
			.; \
	else \
		echo "Service $(SERVICE) not found"; exit 1; \
	fi

docker-build:
	@for svc in $(SERVICES); do \
		$(MAKE) docker-build-one SERVICE=$$svc; \
	done

# ---- Helm ----

HELM_CHART ?= highonsemicolon/pathfinder
HELM_VALUES_DIR ?= helm/values
HELM_NAMESPACE ?= default

helm-deploy:
	@for svc in $(SERVICES); do \
		HELM_RELEASE="$(APP_NAME)-$$svc"; \
		echo "Deploying $$HELM_RELEASE using chart $(HELM_CHART)"; \
		helm upgrade --install "$$HELM_RELEASE" "$(HELM_CHART)" \
			--set image.tag=$(IMAGE_TAG) \
			-f $(HELM_VALUES_DIR)/$$svc.yaml \
			--namespace $(HELM_NAMESPACE) \
			--create-namespace; \
	done
