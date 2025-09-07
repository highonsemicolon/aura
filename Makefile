MAKEFLAGS += --no-print-directory

# ==== Project metadata ====
APP_NAME        ?= aura
BINARY_DIR      ?= .bin
BUILD_DIR       ?= .build
GO_VERSION      ?= 1.25
GO_FLAGS        ?= -trimpath -buildvcs=false
GOBIN           ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
  GOBIN := $(shell go env GOPATH)/bin
endif

# Versions & metadata
VERSION         ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo 0.1.0)
COMMIT_HASH     := $(shell git rev-parse --short=12 HEAD)
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

PROTO_DIRS      := $(shell find apis -type d -name proto 2>/dev/null)
PROTOC_GEN_GO    = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GRPC  = $(GOBIN)/protoc-gen-go-grpc

.DEFAULT_GOAL := help

help:
	@echo "Targets:"; \
	echo "  tidy              - go mod tidy across modules"; \
	echo "  fmt               - go fmt across modules"; \
	echo "  lint              - run golangci-lint"; \
	echo "  test              - unit tests"; \
	echo "  cover             - tests with coverage + threshold"; \
	echo "  build             - build all services"; \
	echo "  run SERVICE=app   - run a single service"; \
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

# ---- Tests ----

COVER_PROFILE := $(BUILD_DIR)/coverage.out
COVER_HTML    := $(BUILD_DIR)/coverage.html
COVER_MIN     ?= 00

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

cover: $(BUILD_DIR)
	@rm -f $(COVER_PROFILE)
	@for m in $(MODULES); do \
		echo "→ Testing $$m"; \
		MOD_NAME=$$(basename $$m); \
		COVER_FILE=$(BUILD_DIR)/$$MOD_NAME.cover.out; \
		go test -coverprofile=$$COVER_FILE -covermode=atomic -race $$m/... || exit 1; \
		if [ -f $$COVER_FILE ]; then \
			if [ ! -f $(COVER_PROFILE) ]; then \
				cp $$COVER_FILE $(COVER_PROFILE); \
			else \
				tail -n +2 $$COVER_FILE >> $(COVER_PROFILE); \
			fi \
		fi; \
	done
	@go tool cover -func=$(COVER_PROFILE) | tee $(BUILD_DIR)/coverage.txt
	@TOTAL=$$(go tool cover -func=$(COVER_PROFILE) | tail -n 1 | awk '{print $$3}' | sed 's/%//'); \
	if [ $${TOTAL%.*} -lt $(COVER_MIN) ]; then \
		echo "\nCoverage $$TOTAL% is below threshold $(COVER_MIN)%"; exit 1; \
	fi
	@go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "Coverage OK (>= $(COVER_MIN)%)"


# Simple multi-module test target
test:
	@for m in $(MODULES); do \
		echo "→ Running tests for $$m"; \
		go test -race -v $$m/... || exit 1; \
	done

# ---- Build / Run ----

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

proto:
	@echo "Generating protobuf with protoc"; \
	for d in $(PROTO_DIRS); do \
		service_dir=$$(dirname $$d); \
		mkdir -p $$service_dir/gen; \
		protoc -I $$d \
			--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
			--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GRPC) \
			--go_out=$$service_dir/gen --go_opt=paths=source_relative \
			--go-grpc_out=$$service_dir/gen --go-grpc_opt=paths=source_relative \
			$$d/*.proto; \
	done; \

# ---- Docker ----

REGISTRY     ?= ghcr.io
IMAGE_OWNER  ?= $(shell basename $(shell dirname $(shell git remote get-url origin 2>/dev/null || echo unknown/unknown)))
IMAGE_TAG    ?= $(VERSION)-$(COMMIT_HASH)
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
		IMG=$(REGISTRY)/$(IMAGE_OWNER)/$(APP_NAME)-$(SERVICE):$(IMAGE_TAG); \
		echo "=== Building $(SERVICE) -> $$IMG ==="; \
		docker buildx build \
			--platform $(PLATFORMS) \
			--build-arg SERVICE=$(SERVICE) \
			--build-arg VERSION=$(VERSION) \
			--build-arg COMMIT=$(COMMIT_HASH) \
			$(CACHE_FROM) \
    		$(CACHE_TO) \
			-t $$IMG \
			-f Dockerfile \
			$(if $(filter true,$(PUSH)),--push,) \
			.; \
	else \
		echo "Service $(SERVICE) not found"; exit 1; \
	fi

docker-push-one:
ifndef SERVICE
	$(error Usage: make docker-push-one SERVICE=app)
endif
	@if [ "$(PUSH)" = "true" ]; then \
		IMG=$(REGISTRY)/$(IMAGE_OWNER)/$(APP_NAME)-$(SERVICE):$(IMAGE_TAG); \
		echo "=== Pushing $$IMG ==="; \
		docker push $$IMG; \
		if [ "$(BRANCH)" = "main" ]; then \
			LATEST=$(REGISTRY)/$(IMAGE_OWNER)/$(APP_NAME)-$(SERVICE):latest; \
			docker tag $$IMG $$LATEST; \
			docker push $$LATEST; \
		fi; \
	else \
		echo "Service $(SERVICE) not found"; exit 1; \
	fi

docker-build:
	@for svc in $(SERVICES); do \
		$(MAKE) docker-build-one SERVICE=$$svc; \
	done

docker-push:
	@for svc in $(SERVICES); do \
		$(MAKE) docker-push-one SERVICE=$$svc; \
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
