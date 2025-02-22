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
