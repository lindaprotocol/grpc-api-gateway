FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make protobuf protobuf-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate protos and build
RUN make proto
RUN go build -o bin/gateway cmd/gateway/main.go
RUN go build -o bin/indexer cmd/indexer/main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binaries
COPY --from=builder /app/bin/gateway /app/bin/
COPY --from=builder /app/bin/indexer /app/bin/

# Copy config
COPY internal/config/config.yaml /app/config/

# Create non-root user
RUN addgroup -g 1000 -S linda && adduser -u 1000 -S linda -G linda
USER linda

EXPOSE 18890 2112

CMD ["/app/bin/gateway", "-config", "/app/config/config.yaml"]