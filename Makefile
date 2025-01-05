BINARY_NAME ?= tmp/main
GO_VERSION ?= 1.23.4
GO_FLAGS ?= -v

.PHONY: all build run run-hot fmt lint test test-coverage clean docker docker-run dev dev-cleanup help ensure-tools version ci clean-docker

# Default target
all: build

# ---- Build & Run Targets ----

build:
	@if [ ! -f cmd/main.go ]; then echo "Error: cmd/main.go not found!"; exit 1; fi
	@echo "Building the Go application..."
	go build $(GO_FLAGS) -o $(BINARY_NAME) cmd/main.go

run: build
	@echo "Running the application..."
	./$(BINARY_NAME)

run-hot: build
	@echo "Running the application with hot reload..."
	air -c ~/.air.toml

# ---- Code Quality Targets ----

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run

# ---- Test Targets ----

test:
	@echo "Running tests..."
	go test ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

# ---- Cleanup Targets ----

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) coverage.out

clean-docker:
	@echo "Removing Docker images and containers..."
	docker system prune -f

# ---- Docker Targets ----

docker:
	@echo "Building Docker image for multi-platform support..."
	docker buildx build --platform linux/amd64,linux/arm64 -t $(BINARY_NAME) .

docker-run: docker
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# ---- Development Environment ----

dev:
	@echo "Setting up development environment with Docker Compose..."
	docker compose up

dev-cleanup:
	@echo "Cleaning up development environment with Docker Compose..."
	docker compose down -v

# ---- Misc Targets ----

help:
	@echo "Available targets:"
	@echo "  all              Build the application"
	@echo "  build            Build the application binary"
	@echo "  run              Run the application locally"
	@echo "  run-hot          Run the application with hot reload"
	@echo "  fmt              Format Go code"
	@echo "  lint             Lint Go code"
	@echo "  test             Run tests"
	@echo "  test-coverage    Run tests with coverage report"
	@echo "  clean            Clean up the build artifacts"
	@echo "  clean-docker     Remove Docker images and containers"
	@echo "  docker           Dockerize the Go application"
	@echo "  docker-run       Build and run the app inside Docker"
	@echo "  dev              Set up the local development environment"
	@echo "  dev-cleanup      Clean up the local development environment"
	@echo "  help             Show this help message"
	@echo "  ensure-tools     Ensure required tools are installed"
	@echo "  version          Display Go and application version"
	@echo "  ci               Run CI pipeline checks"

ensure-tools:
	@echo "Ensuring required tools are installed..."
	@command -v air > /dev/null || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	@command -v golangci-lint > /dev/null || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }

version:
	@echo "Go version: $(GO_VERSION)"
	@echo "Application binary: $(BINARY_NAME)"

ci: fmt lint test
	@echo "CI pipeline checks passed!"
