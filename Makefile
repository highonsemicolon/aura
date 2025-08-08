MAKEFLAGS += --no-print-directory

BINARY_NAME ?= tmp/main
GO_VERSION ?= 1.24.5
GO_FLAGS ?= -v

APP_NAME := Aura
VERSION := 0.1
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_HOST := $(shell hostname)
MODULE := $(shell go list -m)
LDFLAGS := 

tidy:
	@go mod tidy

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

build:
	@if [ ! -f cmd/app/main.go ]; then echo "Error: cmd/app/main.go not found!"; exit 1; fi
	@go build $(GO_FLAGS)  -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/app
	
run:
	@go run -ldflags "$(LDFLAGS)" ./cmd/app

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
