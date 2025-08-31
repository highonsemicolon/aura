MAKEFLAGS += --no-print-directory

BINARY_NAME ?= tmp/main
GO_VERSION ?= 1.24.5
GO_FLAGS ?= -v
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
  GOBIN := $(shell go env GOPATH)/bin
endif

APP_NAME := aura
VERSION := 0.1
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_HOST := $(shell hostname)
MODULES := $(shell go list -m -f '{{.Dir}}')
LDFLAGS := 
PROTO_FILES := $(wildcard apis/*/proto/*.proto)


print-go-version:
	@echo $(GO_VERSION)

tidy:
	@go mod tidy

lint:
	@echo "Linting (per module)..."
	@for m in $(MODULES); do \
		echo "→ Linting $$m"; \
		(cd $$m && go vet ./...); \
	done

fmt:
	@echo "Formatting all modules..."
	@for m in $(MODULES); do \
		echo "→ Formatting $$m"; \
		(cd $$m && go fmt ./...); \
	done

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

proto:
	@for file in $(PROTO_FILES); do \
		service_dir=$$(dirname $$file | sed 's|/proto$$||'); \
		mkdir -p $$service_dir/gen; \
		protoc -I $$service_dir/proto \
			--plugin=protoc-gen-go=$(GOBIN)/protoc-gen-go \
			--plugin=protoc-gen-go-grpc=$(GOBIN)/protoc-gen-go-grpc \
			--go_out=$$service_dir/gen --go_opt=paths=source_relative \
			--go-grpc_out=$$service_dir/gen --go-grpc_opt=paths=source_relative \
			$$file; \
	done
	@echo "Protobuf files generated in apis/*/gen directories."

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
