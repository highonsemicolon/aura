# Set the name of the binary and the Go version
BINARY_NAME=temp/main
GO_VERSION=1.18

# Set default Go build flags
GO_FLAGS=-v

# Default target: Build the application
all: build

# Build the application binary
build:
	@echo "Building the Go application..."
	go build $(GO_FLAGS) -o $(BINARY_NAME) cmd/main.go

# Run the application locally
run: build
	@echo "Running the application..."
	./$(BINARY_NAME)

# Build the application with hot reload using air
run-hot: build
	@echo "Running the application with hot reload..."
	air -c ~/.air.toml

# Format Go code using gofmt
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Lint Go code using golangci-lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Clean up the build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) coverage.out

# Dockerize the Go application
docker:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Build and run the app inside Docker
docker-run: docker
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Set up the local development environment (e.g., Docker Compose)
dev-env:
	@echo "Setting up development environment with Docker Compose..."
	docker-compose up --build
