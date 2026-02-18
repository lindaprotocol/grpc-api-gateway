.PHONY: proto build clean test

proto:
	@echo "Generating protos..."
	@./scripts/gen-proto.sh

build: proto
	@echo "Building gateway..."
	@go build -o bin/gateway cmd/gateway/main.go

clean:
	@rm -rf bin/ pkg/api/*.pb.go

test:
	@go test -v ./...

run: build
	@./bin/gateway -grpc-server-endpoint=localhost:50051 -http-port=:18890

deps:
	@go mod download
	@go mod tidy

.PHONY: proto build clean test run deps
