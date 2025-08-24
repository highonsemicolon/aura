MAKEFLAGS += --no-print-directory

BINARY_NAME ?= tmp/main
GO_VERSION ?= 1.24.5
GO_FLAGS ?= -v

APP_NAME := aura
VERSION := 0.1
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_HOST := $(shell hostname)
MODULE := $(shell go list -m)
LDFLAGS := 

print-go-version:
	@echo $(GO_VERSION)


tidy:
	@go mod tidy

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

build:
	@if [ ! -f services/app/main.go ]; then echo "Error: services/app/main.go not found!"; exit 1; fi
	@go build $(GO_FLAGS)  -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./services/app

test:
	@go test -v ./...
	
run:
	@go run -ldflags "$(LDFLAGS)" ./services/app

run-hot:
	air \
	--build.pre_cmd="make tidy" \
	--build.cmd="make build" \
	--build.bin="$(BINARY_NAME)" \
	--build.send_interrupt=true \
	--build.kill_delay=2s

proto: $(wildcard proto/*.proto)
	protoc \
    --go_out=gen/ \
    --go-grpc_out=gen/ \
    proto/*.proto
	@echo "Protobuf files generated in gen/ directory."

client:
	go run -ldflags "$(LDFLAGS)" ./services/client

setup-helm:
	@echo "Setting up Helm..."
	@helm repo add highonsemicolon https://highonsemicolon.github.io/charts
	@helm repo update

deploy:
	helm upgrade --install \
		$(APP_NAME) \
		highonsemicolon/pathfinder \
		--set image.tag=$(VERSION) \
		-f helm/values.yaml
