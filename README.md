# LindaGrid API Gateway

A high-performance gRPC to HTTP/JSON gateway providing comprehensive API access to the Linda blockchain. This service acts as a unified entry point for all Linda blockchain APIs, including FullNode, Solidity Node, JSON-RPC, and custom Lindascan endpoints.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Architecture](#architecture)
- [Installation](#installation)
- [Configuration](#configuration)
- [Building from Source](#building-from-source)
- [Running the Gateway](#running-the-gateway)
- [Docker Deployment](#docker-deployment)
- [API Documentation](#api-documentation)
- [Authentication](#authentication)
- [Rate Limiting](#rate-limiting)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Overview

The LindaGrid API Gateway is a robust middleware solution that translates gRPC services to HTTP/JSON APIs, providing a unified interface for interacting with the Linda blockchain. It aggregates multiple blockchain API services into a single, coherent gateway:

- **FullNode HTTP API** - Real-time blockchain interactions
- **Solidity Node HTTP API** - Confirmed/finalized blockchain data
- **JSON-RPC API** - Ethereum-compatible JSON-RPC interface
- **Event Query Service** - Indexed event and transaction queries
- **Lindascan Custom APIs** - Explorer-specific endpoints

## Features

- 🚀 **High Performance** - Built on gRPC with HTTP/JSON translation via grpc-gateway
- 🔐 **Security First** - API key authentication, JWT support, and allowlist controls
- ⚡ **Rate Limiting** - Configurable rate limiting with multiple strategies
- 💾 **Caching** - Redis-based response caching for improved performance
- 📊 **Blockchain Indexing** - PostgreSQL-based indexer for historical data
- 🎯 **Comprehensive API Coverage** - All Linda blockchain APIs in one place
- 📈 **Monitoring** - Prometheus metrics and structured logging
- 🐳 **Docker Support** - Containerized deployment with docker-compose

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Core Language** | Go 1.23+ | High-performance backend |
| **API Gateway** | grpc-gateway v2 | gRPC to HTTP/JSON translation |
| **Web Framework** | Gin | HTTP routing and middleware |
| **Database** | PostgreSQL 15+ | Blockchain data indexing |
| **Cache** | Redis 7+ | Rate limiting and response caching |
| **Protocol Buffers** | Protobuf 3 | API contract definition |
| **Authentication** | JWT + API Keys | Secure access control |
| **Logging** | Logrus | Structured logging |
| **Metrics** | Prometheus | System monitoring |

## Architecture

```
┌─────────────┐      ┌──────────────┐     ┌─────────────────┐
│   Clients   │────▶│  API Gateway │────▶│  Linda FullNode │
│ (Web/Mobile)│      │   (Golang)   │     │   (gRPC/HTTP)   │
└─────────────┘      └──────────────┘     └─────────────────┘
                            │                      │
                            ▼                      ▼
                     ┌──────────────┐     ┌─────────────────┐
                     │   Redis      │     │  Solidity Node  │
                     │   Cache      │     │   (gRPC/HTTP)   │
                     └──────────────┘     └─────────────────┘
                            │                      │
                            ▼                      ▼
                     ┌──────────────┐     ┌─────────────────┐
                     │  PostgreSQL  │     │  Event Service  │
                     │   Indexer    │     │   (HTTP/JSON)   │
                     └──────────────┘     └─────────────────┘
```

### Components

1. **API Gateway Core** - Handles HTTP requests, authentication, rate limiting
2. **Blockchain Clients** - gRPC clients for FullNode and Solidity Node
3. **Indexer Service** - Background service for blockchain data indexing
4. **Redis Cache** - In-memory cache for rate limiting and responses
5. **PostgreSQL** - Persistent storage for indexed blockchain data

## Installation

### Prerequisites

- **Go 1.23+** - [Download](https://golang.org/dl/)
- **PostgreSQL 15+** - [Download](https://www.postgresql.org/download/)
- **Redis 7+** - [Download](https://redis.io/download/)
- **Protocol Buffers** - [Download](https://github.com/protocolbuffers/protobuf/releases)
- **Git** - For cloning the repository

### Quick Start

```bash
# Clone the repository
git clone https://github.com/lindaprotocol/grpc-api-gateway.git
cd grpc-api-gateway

# Install dependencies
make deps

# Generate protobuf files
make proto

# Build the binaries
make build

# Set up database and Redis (see Configuration section)

# Run the gateway
make run
```

## Configuration

The gateway is configured via a YAML file located at `internal/config/config.yaml`. Here's a comprehensive configuration example:

```yaml
# internal/config/config.yaml
server:
  http_port: 18890
  grpc_port: 50051
  enable_tls: false
  cert_file: ""
  key_file: ""

environment: "production"  # production, staging, development

linda:
  fullnode_endpoint: "localhost:50051"  # FullNode gRPC endpoint
  solidity_endpoint: "localhost:50061"  # Solidity Node gRPC endpoint
  event_endpoint: "localhost:8080"      # Event Service HTTP endpoint
  grpc_timeout: 30s
  max_msg_size: 10485760  # 10MB

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "lindascan"
  password: "your_password"
  dbname: "lindascan"
  sslmode: "disable"
  max_connections: 100
  idle_connections: 10

redis:
  addr: "localhost:6379"
  password: ""  # Set if Redis requires authentication
  db: 0
  pool_size: 100
  min_idle_conns: 10

auth:
  api_key_enabled: true
  jwt_enabled: true
  jwt_secret: "your-secret-key"  # Change in production!
  default_rate_limit_qps: 15
  default_daily_limit: 100000

rate_limit:
  enabled: true
  default_qps: 15
  default_burst: 30
  strategy: "token_bucket"  # token_bucket, sliding_window, leaky_bucket

cors:
  allowed_origins:
    - "https://lindascan.org"
    - "https://*.lindascan.org"

indexer:
  enabled: true
  block_batch_size: 100
  sync_interval: 5s
  start_block: 0
  max_workers: 10

logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # json, text
```

### Environment Variables

Key configuration can be overridden with environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `LINDA_FULLNODE_ENDPOINT` | FullNode gRPC endpoint | `localhost:50051` |
| `LINDA_SOLIDITY_ENDPOINT` | Solidity Node gRPC endpoint | `localhost:50061` |
| `DATABASE_PASSWORD` | PostgreSQL password | - |
| `REDIS_PASSWORD` | Redis password | - |
| `JWT_SECRET` | JWT signing secret | - |

## Building from Source

### Using Make (Recommended)

```bash
# Download dependencies
make deps

# Generate protobuf files
make proto

# Build all binaries
make build

# Build individual components
make build-gateway  # Build only gateway
make build-indexer  # Build only indexer
make build-apikey   # Build only API key tool

# Clean build artifacts
make clean
```

### Manual Build

```bash
# Download dependencies
go mod download
go mod tidy

# Generate protobufs
chmod +x scripts/gen-proto.sh
./scripts/gen-proto.sh

# Build binaries
go build -o bin/gateway ./cmd/gateway
go build -o bin/indexer ./cmd/indexer
go build -o bin/apikey ./cmd/apikey
```

## Running the Gateway

### 1. Set Up Database

```bash
# Create PostgreSQL database and user
sudo -u postgres psql

CREATE USER lindascan WITH PASSWORD 'your_password';
CREATE DATABASE lindascan OWNER lindascan;
GRANT ALL PRIVILEGES ON DATABASE lindascan TO lindascan;
\q
```

### 2. Start Redis

```bash
# Install Redis if not already installed
sudo apt update
sudo apt install redis-server -y

# Start Redis
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

### 3. Start the Gateway

```bash
# Run the gateway (will auto-migrate database)
make run

# Or run with custom config
./bin/gateway -config ./internal/config/config.yaml
```

### 4. Start the Indexer (Optional)

```bash
# Run the indexer in a separate terminal
make run-indexer

# Or directly
./bin/indexer -config ./internal/config/config.yaml
```

## Docker Deployment

### Using Docker Compose

```bash
# Build and start all services
docker-compose -f docker/docker-compose.yml up -d

# View logs
docker-compose -f docker/docker-compose.yml logs -f

# Stop services
docker-compose -f docker/docker-compose.yml down
```

### Building Individual Docker Images

```bash
# Build gateway image
docker build -f docker/Dockerfile -t lindagrid/gateway .

# Run container
docker run -p 18890:18890 -v $(pwd)/config:/app/config lindagrid/gateway
```

## API Documentation

The gateway exposes multiple API endpoints organized into sections:

### API Categories

| Category | Base Path | Description |
|----------|-----------|-------------|
| **Lindagrid V1** | `/v1/` | Event query service and indexed data |
| **FullNode HTTP API** | `/wallet/` | Real-time blockchain interactions |
| **Solidity Node API** | `/walletsolidity/` | Confirmed/finalized blockchain data |
| **JSON-RPC API** | `/jsonrpc` | Ethereum-compatible JSON-RPC |
| **Lindascan Custom** | `/api/` | Explorer-specific endpoints |
| **External** | `/external/` | Tag system and file uploads |
| **Monitoring** | `/monitor/` | Node health and metrics |

### API Categories Detail

#### Lindagrid V1 (Event Query Service)
- `GET /v1/transactions` - List transactions with pagination
- `GET /v1/transactions/{hash}` - Get transaction by hash
- `GET /v1/events` - Query blockchain events
- `GET /v1/events/transaction/{transactionId}` - Get events by transaction
- `GET /v1/events/{contractAddress}` - Get events by contract
- `GET /v1/blocks` - List blocks
- `GET /v1/blocks/{hash}` - Get block by hash

#### FullNode HTTP API
- `POST /wallet/getaccount` - Get account information
- `POST /wallet/createtransaction` - Create a transaction
- `POST /wallet/broadcasttransaction` - Broadcast signed transaction
- `POST /wallet/getnowblock` - Get latest block
- `GET /wallet/listnodes` - List connected nodes

#### FullNode Solidity HTTP API
- `POST /walletsolidity/getaccount` - Get confirmed account info
- `POST /walletsolidity/gettransactionbyid` - Get confirmed transaction
- `POST /walletsolidity/getnowblock` - Get latest confirmed block

#### Full Node JSON-RPC API
- `POST /jsonrpc` - Ethereum-compatible JSON-RPC endpoint

### Example API Calls

```bash
# Get account information
curl -X POST http://localhost:18890/wallet/getaccount \
  -H "Content-Type: application/json" \
  -d '{"address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr", "visible": true}'

# Get latest block
curl -X POST http://localhost:18890/wallet/getnowblock

# Query events (Lindagrid V1)
curl "http://localhost:18890/v1/events?limit=10&contract=TG3XXyExBkPp9nzdajDZsozEu4BkaSJozs"

# JSON-RPC call
curl -X POST http://localhost:18890/jsonrpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

## Authentication

### API Key Authentication

Generate an API key using the included tool:

```bash
# Generate a new API key
./bin/apikey generate user123 "My Application"

# List API keys for a user
./bin/apikey list user123

# Revoke an API key
./bin/apikey revoke key_id_here
```

Use the API key in requests:

```bash
curl -X GET http://localhost:18890/v1/transactions \
  -H "LINDA-PRO-API-KEY: your-api-key-here"
```

### JWT Authentication

JWT tokens can be used for additional security:

```bash
curl -X POST http://localhost:18890/wallet/getaccount \
  -H "Authorization: Bearer your-jwt-token" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -d '{"address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr"}'
```

## Rate Limiting

The gateway implements configurable rate limiting:

| User Type | Default QPS | Daily Limit |
|-----------|------------|-------------|
| Authenticated (with API key) | 15 QPS | 100,000 |
| Anonymous | 5 QPS | Strict |

Rate limit headers are included in responses:
- `X-RateLimit-Limit` - Maximum requests per second
- `X-RateLimit-Remaining` - Remaining requests in current window
- `X-RateLimit-Reset` - Time when the limit resets

## Monitoring

### Prometheus Metrics

Metrics are available at `http://localhost:2112/metrics`:

```
# Request metrics
http_requests_total{method="GET", endpoint="/v1/transactions"} 1245
http_request_duration_seconds{quantile="0.95"} 0.023

# Cache metrics
cache_hits_total{cache="account"} 892
cache_misses_total{cache="account"} 234

# Rate limit metrics
rate_limit_exceeded_total{user_type="anonymous"} 45
```

### Health Check

```bash
# Basic health check
curl http://localhost:18890/health

# Response
{"status":"ok","time":1645564800}
```

## Troubleshooting

### Common Issues and Solutions

#### Database Connection Failed

```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Verify database credentials
PGPASSWORD=your_password psql -h localhost -U lindascan -d lindascan -c "SELECT 1"
```

#### Redis Connection Refused

```bash
# Check Redis status
sudo systemctl status redis-server

# Test Redis connection
redis-cli ping
```

#### gRPC Connection Issues

```bash
# Verify Linda FullNode is running
nc -zv localhost 50051

# Check gRPC endpoint in config
grep fullnode_endpoint internal/config/config.yaml
```

#### Rate Limit Exceeded

```bash
# Check rate limit headers
curl -I http://localhost:18890/v1/transactions \
  -H "LINDA-PRO-API-KEY: your-key"

# Response headers should include rate limit information
```

### Logs

```bash
# View gateway logs
tail -f /var/log/lindagrid/gateway.log

# View indexer logs
tail -f /var/log/lindagrid/indexer.log

# Structured logs (JSON format)
journalctl -u lindagrid-gateway -o json-pretty
```

## Performance Tuning

### Database Optimization

```sql
-- Create indexes for common queries
CREATE INDEX CONCURRENTLY idx_transactions_from_to ON transactions(from_address, to_address);
CREATE INDEX CONCURRENTLY idx_events_contract_time ON events(contract_address, block_timestamp DESC);

-- Vacuum analyze for query planner
VACUUM ANALYZE;
```

### Cache Configuration

Adjust Redis cache TTLs in `config.yaml`:

```yaml
cache:
  default_ttl: 300  # 5 minutes
  account_ttl: 300
  block_ttl: 600    # 10 minutes
  transaction_ttl: 600
```

### gRPC Connection Pool

```yaml
linda:
  grpc_timeout: 30s
  max_msg_size: 10485760  # 10MB
  # Connection pool settings
  max_connections: 100
  idle_timeout: 60s
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

### Code Style

- Follow standard Go formatting: `gofmt -s -w .`
- Run linter: `golangci-lint run`
- Ensure all tests pass: `go test ./...`

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [https://docs.lindagrid.lindacoin.org](https://docs.lindagrid.lindacoin.org)
- **GitHub Issues**: [https://github.com/lindaprotocol/grpc-api-gateway/issues](https://github.com/lindaprotocol/grpc-api-gateway/issues)
- **Discord**: [Linda Protocol Discord](https://discord.gg/lindacoin)

---
