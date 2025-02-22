BINARY_NAME ?= tmp/main
GO_VERSION ?= 1.23.4
GO_FLAGS ?= -v


tidy:
	go mod tidy

fmt:
	@echo "Formatting Go code..."
	go fmt ./...
lint:
	@echo "Running linter..."
	golangci-lint run --allow-parallel-runners --no-config 

codegen:
	oapi-codegen \
	-generate types,gin,strict-server,spec \
	-package api \
	-o ./src/api/api.gen.go ./openapi/spec.yaml
	
static-fix:
	@echo "Fixing go-staticcheck ST1005 errors in generated code..."
	go run openapi/fix_errors.go ./src/api

run:
	go run ./src/main.go

run-hot:
	air --build.cmd="go run ./src/main.go"

build:
	@if [ ! -f src/main.go ]; then echo "Error: src/main.go not found!"; exit 1; fi
	@echo "Building the Go application..."
	go build $(GO_FLAGS) -o $(BINARY_NAME) src/main.go

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