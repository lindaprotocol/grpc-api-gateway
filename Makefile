.PHONY: proto build clean test run deps docker-up docker-down lint help init

PROJECT_NAME = grpc-api-gateway
GO = go
GOFLAGS = -v
BIN_DIR = bin
CMD_DIR = cmd
MAIN_GATEWAY = ./cmd/gateway
MAIN_INDEXER = ./cmd/indexer
MAIN_APIKEY = ./cmd/apikey

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@echo '  init          Initialize project (deps + proto)'
	@echo '  deps          Download dependencies'
	@echo '  proto         Generate protobuf files'
	@echo '  build         Build all binaries'
	@echo '  clean         Clean build artifacts'
	@echo '  run-gateway   Run gateway locally'
	@echo '  run-indexer   Run indexer locally'

init: deps proto
	@echo "Project initialized!"

deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies downloaded!"

proto:
	@echo "Generating protobuf files..."
	@chmod +x scripts/gen-proto.sh 2>/dev/null || true
	@./scripts/gen-proto.sh
	@echo "Proto generation complete!"

build: proto
	@echo "Building binaries..."
	mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/gateway $(MAIN_GATEWAY)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/indexer $(MAIN_INDEXER)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/apikey $(MAIN_APIKEY)
	@echo "Build complete! Binaries in $(BIN_DIR)/"

build-gateway:
	@echo "Building gateway..."
	mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/gateway $(MAIN_GATEWAY)
	@echo "Gateway build complete!"

build-indexer:
	@echo "Building indexer..."
	mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/indexer $(MAIN_INDEXER)
	@echo "Indexer build complete!"

clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	rm -rf pkg/lindapb/*.pb.go
	rm -rf pkg/lindapb/core/*.pb.go
	@echo "Clean complete!"

test:
	@echo "Running tests..."
	$(GO) test ./... -cover
	@echo "Tests complete!"

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Formatting complete!"

vet:
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "Vet complete!"

run-gateway: build-gateway
	@echo "Starting gateway..."
	./$(BIN_DIR)/gateway -config ./internal/config/config.yaml

run-indexer: build-indexer
	@echo "Starting indexer..."
	./$(BIN_DIR)/indexer -config ./internal/config/config.yaml

run: run-gateway