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
  -d '{"address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko"}'
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

## API Examples

### Core Wallet APIs

#### Get Account Information

```bash
# Using POST
curl -X POST http://localhost:18890/wallet/getaccount \
  -H "Content-Type: application/json" \
  -d '{"address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko"}'

# Using GET
curl "http://localhost:18890/wallet/getaccount?address=LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko"
```

**Response:**
```json
{
  "address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
  "balance": 1000000000,
  "create_time": 1771509351000,
  "allowance": 0,
  "latest_withdraw_time": 0
}
```

#### Create a Transaction

```bash
curl -X POST http://localhost:18890/wallet/createtransaction \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
    "to_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "amount": 1000000
  }'
```

**Response:**
```json
{
  "txID": "0x9fc5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5",
  "raw_data": {
    "contract": [{
      "parameter": {
        "value": {
          "amount": 1000000,
          "owner_address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
          "to_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD"
        },
        "type_url": "type.googleapis.com/protocol.TransferContract"
      },
      "type": "TransferContract"
    }],
    "ref_block_bytes": "0x1234",
    "ref_block_hash": "0x5678",
    "expiration": 1771509354000,
    "timestamp": 1771509351000
  },
  "raw_data_hex": "0x0a..."
}
```

#### Broadcast a Signed Transaction

```bash
curl -X POST http://localhost:18890/wallet/broadcasttransaction \
  -H "Content-Type: application/json" \
  -d '{
    "txID": "0x9fc5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5",
    "raw_data": {...},
    "signature": ["0x1234567890abcdef..."],
    "raw_data_hex": "0x0a..."
  }'
```

**Response:**
```json
{
  "result": true,
  "code": 0,
  "message": ""
}
```

#### Get Current Block

```bash
curl http://localhost:18890/wallet/getnowblock
```

**Response:**
```json
{
  "blockID": "0x000000000000ca3c51d35c70e29d540224500cd53c838962154d8b8134d2ad56",
  "block_header": {
    "raw_data": {
      "number": 51773,
      "txTrieRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "witness_address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
      "parentHash": "0x000000000000ca3c51d35c70e29d540224500cd53c838962154d8b8134d2ad55",
      "timestamp": 1771509351000
    },
    "witness_signature": "0xfafaf8a10c6c055b15fe1d2f2ed06e2723d8cdf5e88fb513c50b50b5e0383939"
  },
  "transactions": []
}
```

#### Get Block by Number

```bash
# Using POST
curl -X POST http://localhost:18890/wallet/getblockbynum \
  -H "Content-Type: application/json" \
  -d '{"num": 51773}'

# Using GET
curl http://localhost:18890/wallet/getblockbynum/51773
```

#### Get Transaction by ID

```bash
# Using POST
curl -X POST http://localhost:18890/wallet/gettransactionbyid \
  -H "Content-Type: application/json" \
  -d '{"value": "0x9fc5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5"}'

# Using GET
curl http://localhost:18890/wallet/gettransactionbyid/0x9fc5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5
```

#### List Witnesses

```bash
curl http://localhost:18890/wallet/listwitnesses
```

**Response:**
```json
{
  "witnesses": [
    {
      "address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
      "voteCount": 1000000,
      "url": "https://lindawitness.com",
      "totalProduced": 5000,
      "totalMissed": 10,
      "latestBlockNum": 51773
    }
  ]
}
```

### Explorer & Scan APIs

#### Get Homepage Bundle (Dashboard)

```bash
curl http://localhost:18890/api/system/homepage-bundle
```

**Response:**
```json
{
  "totalBlocks": 51773,
  "totalTransactions": 1254321,
  "totalAccounts": 89234,
  "totalContracts": 1567,
  "totalTokens": 892,
  "priceUSD": 0.0123,
  "marketCap": 12345678,
  "volume24h": 123456,
  "recentBlocks": [
    {
      "number": 51773,
      "hash": "0x000000000000ca3c51d35c70e29d540224500cd53c838962154d8b8134d2ad56",
      "timestamp": 1771509351000,
      "transactions": 0,
      "witness": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko"
    }
  ],
  "recentTransactions": [
    {
      "hash": "0x9fc5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5",
      "timestamp": 1771509351000,
      "from": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
      "to": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "amount": 1000000,
      "block": 51773
    }
  ]
}
```

#### Get Top 10 Accounts

```bash
curl "http://localhost:18890/api/top10?type=accounts"
```

**Response:**
```json
[
  {
    "rank": 1,
    "address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
    "balance": 10000000000,
    "percentage": 5.2
  },
  {
    "rank": 2,
    "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "balance": 8500000000,
    "percentage": 4.4
  }
]
```

#### List Tokens with Pagination

```bash
curl "http://localhost:18890/api/token?limit=10&start=1&sort=issue_time"
```

**Response:**
```json
{
  "total": 892,
  "page": 1,
  "limit": 10,
  "tokens": [
    {
      "name": "LindaToken",
      "symbol": "LDT",
      "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "totalSupply": 1000000000,
      "decimals": 6,
      "holders": 1250,
      "transfers": 50000
    }
  ]
}
```

#### Get Token Holders

```bash
curl "http://localhost:18890/api/tokenholders?contract=LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD&limit=20"
```

**Response:**
```json
{
  "contract": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "totalHolders": 1250,
  "holders": [
    {
      "address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
      "balance": 500000000,
      "percentage": 50.0
    }
  ]
}
```

#### Get Account Resource Information

```bash
curl "http://localhost:18890/api/account/resource?address=LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko"
```

**Response:**
```json
{
  "address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
  "freeNetUsed": 100,
  "freeNetLimit": 5000,
  "energyUsed": 0,
  "energyLimit": 0,
  "netUsed": 0,
  "netLimit": 0
}
```

#### Get Blocks with Pagination

```bash
curl "http://localhost:18890/api/block?limit=10&start=1&sort=desc"
```

**Response:**
```json
{
  "total": 51773,
  "blocks": [
    {
      "number": 51773,
      "hash": "0x000000000000ca3c51d35c70e29d540224500cd53c838962154d8b8134d2ad56",
      "timestamp": 1771509351000,
      "transactions": 0,
      "size": 512,
      "witness": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko"
    }
  ]
}
```

#### Search

```bash
curl "http://localhost:18890/api/search?query=LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko"
```

**Response:**
```json
{
  "results": [
    {
      "type": "address",
      "id": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
      "name": "Account",
      "url": "#/address/LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
      "description": "Balance: 1,000,000,000 LINDA"
    },
    {
      "type": "transaction",
      "id": "0x9fc5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5",
      "name": "Transaction",
      "url": "#/transaction/0x9fc5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5",
      "description": "Value: 1,000,000 LINDA"
    }
  ]
}
```

#### Get Contract Events

```bash
curl "http://localhost:18890/events?limit=10&sort=-timeStamp"
```

**Response:**
```json
{
  "events": [
    {
      "transaction_id": "0x9fc5b2d5c5e5e5b2d5c5e5e5b2d5c5e5e5",
      "contract_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "event_name": "Transfer",
      "block_number": 51773,
      "timestamp": 1771509351000,
      "result": {
        "from": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
        "to": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "value": "1000000"
      }
    }
  ]
}
```

#### Get Events by Contract

```bash
curl "http://localhost:18890/events/LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD?limit=20"
```

#### Get Market Pairs

```bash
curl "http://localhost:18890/api/exchange/marketPair/list?sortType=volume"
```

**Response:**
```json
{
  "pairs": [
    {
      "pair": "LINDA/USDT",
      "lastPrice": 0.0123,
      "volume24h": 1234567,
      "change24h": 5.2,
      "high24h": 0.0125,
      "low24h": 0.0118,
      "hot": true
    }
  ]
}
```

#### Exchange Rate Calculation

```bash
curl "http://localhost:18890/api/exchange/calc?pair=LINDA/USDT&amount=1000&type=sell"
```

**Response:**
```json
{
  "pair": "LINDA/USDT",
  "amount": 1000,
  "type": "sell",
  "rate": 0.0123,
  "result": 12.3,
  "fee": 0.0123,
  "netResult": 12.2877
}
```

### External/Tag APIs

#### Get Tags for an Address

```bash
curl "http://localhost:18890/external/tag?address=LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko&limit=10"
```

**Response:**
```json
{
  "address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
  "total": 2,
  "tags": [
    {
      "id": "tag_123",
      "tag": "exchange",
      "description": "Centralized exchange wallet",
      "owner": "TAGGING_SERVICE",
      "verified": true,
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "tag_456",
      "tag": "hot_wallet",
      "description": "Exchange hot wallet",
      "owner": "COMMUNITY",
      "verified": false,
      "created_at": "2024-01-02T00:00:00Z"
    }
  ]
}
```

#### Insert a New Tag

```bash
curl -X POST http://localhost:18890/external/tag/insert \
  -H "Content-Type: application/json" \
  -d '{
    "address": "LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko",
    "tag": "exchange",
    "description": "Centralized exchange wallet",
    "owner": "MY_SERVICE",
    "signature": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Tag inserted successfully",
  "tag_id": "tag_789"
}
```

## Response Format

All successful responses return JSON objects with:
- Binary fields (hashes, signatures) as **hex strings**
- Addresses in **Base58 format**

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

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:18890/ws');

// Subscribe to new blocks
ws.send(JSON.stringify({
  type: 'subscribe',
  channel: 'blocks'
}));

// Subscribe to transactions for a specific address
ws.send(JSON.stringify({
  type: 'subscribe',
  channel: 'transactions',
  address: 'LXBzYhTSEWYaVH8cdtMAbPhgrfcMN7jZko'
}));

// Listen for messages
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('New event:', data);
};
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
    environment:
      - GIN_MODE=release
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:18890/wallet/getnowblock"]
      interval: 30s
      timeout: 10s
      retries: 3

  grpc-server:
    image: lindacoin/grpc-server:latest
    ports:
      - "50051:50051"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "grpc_health_probe", "-addr=localhost:50051"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Nginx Configuration Example

```nginx
server {
    listen 80;
    server_name api.lindacoin.org;
    
    location / {
        proxy_pass http://localhost:18890;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Rate limiting
        limit_req zone=api burst=20 nodelay;
        limit_req_status 429;
    }
    
    # WebSocket support
    location /ws {
        proxy_pass http://localhost:18890/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
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

### Development Guidelines

- Write tests for new features
- Update documentation for API changes
- Follow Go best practices and formatting
- Use meaningful commit messages

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Support

- **GitHub Issues**: [https://github.com/lindaprotocol/grpc-api-gateway/issues](https://github.com/lindaprotocol/grpc-api-gateway/issues)
- **Linda Documentation**: [https://docs.lindacoin.org](https://docs.lindacoin.org)
- **Community Forum**: [https://forum.lindacoin.org](https://forum.lindacoin.org)
- **Discord**: [https://discord.gg/lindacoin](https://discord.gg/lindacoin)

## Acknowledgments

- Linda blockchain team
- Contributors and community members
- gRPC-gateway project
- Protocol Buffers team