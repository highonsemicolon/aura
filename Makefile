BINARY_NAME ?= tmp/main
GO_VERSION ?= 1.23.6
GO_FLAGS ?= -v

APP_NAME := Aura
VERSION := 1.3.0
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_HOST := $(shell hostname)
MODULE := $(shell go list -m)
LDFLAGS := -X '$(MODULE)/src/api.appName=$(APP_NAME)' \
		   -X '$(MODULE)/src/api.version=$(VERSION)' \
           -X '$(MODULE)/src/api.commitHash=$(COMMIT_HASH)' \
           -X '$(MODULE)/src/api.buildTime=$(BUILD_TIME)' \
           -X '$(MODULE)/src/api.buildHost=$(BUILD_HOST)'


tidy:
	@go mod tidy

fmt:
	@echo "Formatting Go code..."
	go fmt ./...
lint:
	@echo "Running linter..."
	golangci-lint run --allow-parallel-runners --no-config 

codegen:
	oapi-codegen -config oapi-config.yaml ./openapi/spec.yaml
	
static-fix:
	@echo "Fixing go-staticcheck ST1005 errors in generated code..."
	go run openapi/fix_errors.go ./src/api

run:
	go run -ldflags "$(LDFLAGS)" ./src

run-hot:
	air \
	--build.pre_cmd="make tidy" \
	--build.cmd="make build" \
	--build.bin="./tmp/main"

build:
	@if [ ! -f src/main.go ]; then echo "Error: src/main.go not found!"; exit 1; fi
	@echo "Building the Go application..."
	@go build $(GO_FLAGS)  -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./src

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) coverage.out

help:
	@echo "Available targets:"
	@echo "  build            Build the application binary"
	@echo "  run              Run the application locally"
	@echo "  run-hot          Run the application with hot reload"
	@echo "  fmt              Format Go code"
	@echo "  lint             Lint Go code"
	@echo "  clean            Clean up the build artifacts"
	@echo "  help             Show this help message"