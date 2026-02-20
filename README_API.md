# Linda gRPC API Gateway

A high-performance HTTP gateway for the Linda blockchain, providing RESTful access to Linda gRPC services. This gateway converts gRPC calls to JSON/HTTP endpoints, making it easy to interact with the Linda blockchain from any programming language.

## Features

- 🔄 **Protocol Buffers to JSON** - Automatic conversion with hex-encoded binary data
- 🚀 **High Performance** - Lightweight proxy with minimal overhead
- 📚 **Comprehensive APIs** - Wallet operations, blockchain queries, token management, and explorer endpoints
- 🔍 **Search Capabilities** - Unified search across blocks, transactions, addresses, and tokens
- 📊 **Market Data** - Exchange rates and trading pair information
- 🔌 **WebSocket Support** - Real-time updates for events and transactions

## Quick Start

### Using Docker

```bash
# Build the Docker image
docker build -t linda-gateway .

# Run with Docker Compose
docker-compose up -d
```

### From Source

```bash
# Build the gateway
make build

# Start the gateway
./bin/gateway -grpc-server-endpoint=localhost:50051 -http-port=:18890
```

### Test the API

```bash
# Get current block
curl http://localhost:18890/wallet/getnowblock

# Get account info
curl -X POST http://localhost:18890/wallet/getaccount \
  -H "Content-Type: application/json" \
  -d '{"address": "MDWULEJI6D8BqFZfRSsambEvRNl7"}'
```

## Configuration

The gateway can be configured via command-line flags:

| Flag | Description | Default |
|------|-------------|---------|
| `-grpc-server-endpoint` | Linda gRPC server address | `localhost:50051` |
| `-http-port` | HTTP server port | `:18890` |

## API Documentation

The gateway provides two main API categories:

### Core Wallet APIs

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/wallet/getaccount` | POST/GET | Get account information |
| `/wallet/createtransaction` | POST/GET | Create a transfer transaction |
| `/wallet/broadcasttransaction` | POST/GET | Broadcast a signed transaction |
| `/wallet/getnowblock` | POST/GET | Get the most recent block |
| `/wallet/getblockbynum` | POST/GET | Get block by number |
| `/wallet/getblockbyid` | POST/GET | Get block by hash |
| `/wallet/gettransactionbyid` | POST/GET | Get transaction by hash |
| `/wallet/gettransactioninfobyid` | POST/GET | Get transaction execution info |
| `/wallet/getassetissuelist` | POST/GET | List all issued assets |
| `/wallet/getassetissuebyname` | POST/GET | Get asset by name |
| `/wallet/listwitnesses` | POST/GET | List all witnesses |

### Explorer & Scan APIs

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/system/homepage-bundle` | GET | Dashboard statistics |
| `/api/nodemap` | GET | Node geolocation information |
| `/api/top10` | GET | Top accounts or witnesses |
| `/api/token` | GET | List tokens with pagination |
| `/api/token_lrc20` | GET | List LRC20 tokens |
| `/api/tokenholders` | GET | Token holder distribution |
| `/api/token_lrc20/transfers` | GET | Token transfer history |
| `/api/account/list` | GET | Paginated account list |
| `/api/account/resource` | GET | Account resource information |
| `/api/stats/overview` | GET | Blockchain statistics |
| `/api/block` | GET | Paginated block list |
| `/api/transaction` | GET | Paginated transaction list |
| `/api/search` | GET | Unified search |
| `/events` | GET | Contract events |
| `/events/transaction/{hash}` | GET | Events by transaction |
| `/events/{address}` | GET | Events by contract |
| `/api/exchange/marketPair/list` | GET | Available trading pairs |
| `/api/exchange/calc` | GET | Exchange rate calculation |

### External/Tag APIs

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/external/tag` | GET | Get tags for an address |
| `/external/tag/insert` | POST | Add a new tag to an address |

## Response Format

All successful responses return JSON objects with:
- Binary fields (hashes, signatures) as **hex strings**
- Addresses in **Base58 format**

### Success Response Example

```json
{
  "blockHeader": {
    "rawData": {
      "timestamp": "1771509351000",
      "txTrieRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "parentHash": "0x000000000000ca3c51d35c70e29d540224500cd53c838962154d8b8134d2ad56",
      "number": "51773",
      "witnessId": "0",
      "witnessAddress": "MDWULEJI6D8BqFZfRSsambEvRNl7"
    },
    "witnessSignature": "0xfafaf8a10c6c055b15fe1d2f2ed06e2723d8cdf5e88fb513c50b50b5e03839390f2ab85412ae17320172e90f52ee14b7f43fa38ec5f9391bbade33e342d5eaf0"
  },
  "transactions": []
}
```

### Error Response Format

```json
{
  "error": "error message",
  "code": 500,
  "details": []
}
```

## Error Codes

| Code | Description |
|------|-------------|
| 0 | SUCCESS |
| 1 | SIGERROR (signature error) |
| 2 | CONTRACT_VALIDATE_ERROR |
| 3 | CONTRACT_EXE_ERROR |
| 4 | BANDWITH_ERROR |
| 5 | DUP_TRANSACTION_ERROR |
| 6 | TAPOS_ERROR |
| 7 | TOO_BIG_TRANSACTION_ERROR |
| 8 | TRANSACTION_EXPIRATION_ERROR |
| 9 | SERVER_BUSY |
| 20 | OTHER_ERROR |

## WebSocket Support

Real-time updates are available via WebSocket connections:

```
ws://your-gateway:18890/ws
```

## Project Structure

```
.
├── cmd/
│   ├── gateway/          # Main HTTP gateway
│   └── scanner/          # Optional scan service
├── internal/
│   ├── config/           # Configuration
│   ├── service/          # Business logic
│   └── storage/          # Database models
├── pkg/
│   ├── api/              # Generated protobuf code
│   │   └── protocol/
│   └── middleware/       # HTTP middleware
├── proto/                # Protocol buffer definitions
│   ├── api.proto
│   └── core/
└── scripts/              # Build scripts
```

## Development

### Prerequisites

- Go 1.19 or higher
- Protocol Buffers compiler
- Make

### Generating Protocol Buffers

```bash
make proto
```

### Running Tests

```bash
make test
```

### Clean Build

```bash
make clean && make build
```

### Running Locally

```bash
# Start the gateway
go run cmd/gateway/main.go -grpc-server-endpoint=localhost:50051 -http-port=:18890

# In another terminal, test the API
curl http://localhost:18890/wallet/getnowblock
```

## Production Deployment

### Rate Limiting

The gateway does not implement rate limiting by default. For production deployments, it's recommended to use a reverse proxy like Nginx with rate limiting enabled.

### Docker Compose Example

```yaml
version: '3.8'

services:
  gateway:
    build: .
    ports:
      - "18890:18890"
    command: ["./bin/gateway", "-grpc-server-endpoint=grpc-server:50051"]
    depends_on:
      - grpc-server
    restart: unless-stopped

  grpc-server:
    image: lindacoin/grpc-server:latest
    ports:
      - "50051:50051"
    restart: unless-stopped
```

## Version History

- **v1.3.0** - Added market and exchange APIs
- **v1.2.0** - Added event query APIs
- **v1.1.0** - Added Scan Service APIs
- **v1.0.0** - Initial release with core Wallet APIs

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Support

- **GitHub Issues**: [https://github.com/lindaprotocol/grpc-api-gateway/issues](https://github.com/lindaprotocol/grpc-api-gateway/issues)
- **Linda Documentation**: [https://docs.lindacoin.org](https://docs.lindacoin.org)
- **Community Forum**: [https://forum.lindacoin.org](https://forum.lindacoin.org)

## Acknowledgments

- Linda blockchain team
- Contributors and community members
- gRPC-gateway project