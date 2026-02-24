# Lindagrid V1 API

The Lindagrid V1 API provides indexed blockchain data through a RESTful interface. It's designed for efficient querying of transactions, transfers, events, blocks, and contract logs with pagination support.

## Base URL

All V1 API endpoints are prefixed with `/v1/`

```
https://api.lindagrid.lindacoin.org/v1/
```

## Headers

| Header | Description | Required |
|--------|-------------|----------|
| `LINDA-PRO-API-KEY` | Your API key for authentication | Yes (for rate limiting) |
| `Content-Type` | `application/json` | For POST requests |

## Pagination

Most list endpoints support pagination with the following parameters:

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `limit` | integer | Number of items per page (max 200) | 20 |
| `start` | integer | Starting page index | 0 |
| `sort` | string | Sort field with direction (`-` for descending) | `-timestamp` |
| `fingerprint` | string | Pagination cursor from previous response | - |

Response includes a `meta` object with pagination information:

```json
{
  "data": [...],
  "success": true,
  "meta": {
    "at": 1645564800000,
    "page_size": 20,
    "fingerprint": "next_page_cursor",
    "links": {
      "next": "/v1/transactions?fingerprint=next_page_cursor&limit=20"
    }
  }
}
```

---

## Transactions

### Get Transaction List

Retrieves a paginated list of transactions.

**Endpoint:** `GET /v1/transactions`

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `start` | integer | Start page (default: 0) |
| `block` | integer | Filter by block number >= this value (default: 0) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/transactions?limit=2&sort=-timestamp&start=0&block=1000000" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "blockNumber": 1000001,
      "blockTimestamp": 1645564800000,
      "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
      "value": "1000000",
      "fee": 100000,
      "contractAddress": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N"
    },
    {
      "id": "8dd26d1772231569f022adb42f7d7161dee88b97b4b35eeef6ce73fcd6613bc2",
      "blockNumber": 1000000,
      "blockTimestamp": 1645564740000,
      "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
      "value": "500000",
      "fee": 50000
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 2,
    "fingerprint": "eyJvZmZzZXQiOjJ9",
    "links": {
      "next": "/v1/transactions?fingerprint=eyJvZmZzZXQiOjJ9&limit=2"
    }
  }
}
```

---

### Get Transaction by Hash

Retrieves a specific transaction by its hash.

**Endpoint:** `GET /v1/transactions/{hash}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `hash` | string | Transaction hash (64 characters hex) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/transactions/7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "blockNumber": 1000001,
      "blockTimestamp": 1645564800000,
      "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
      "value": "1000000",
      "fee": 100000,
      "contractAddress": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N"
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

## Transfers

### Get Transfers List

Retrieves a paginated list of token transfers (LIND and LRC10/LRC20 tokens).

**Endpoint:** `GET /v1/transfers`

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `start` | integer | Start page (default: 0) |
| `from` | string | Filter by sender address |
| `to` | string | Filter by recipient address |
| `token` | string | Filter by token name or contract address |
| `block` | integer | Filter by block number >= this value |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/transfers?token=USDT&limit=2&from=LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD&to=LeudxtcgoduEyhFTuNzquusYk6yuM73iRr" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "transaction_id": "70d655a17e04d6b6b7ee5d53e7f37655974f4e71b0edd6bcb311915a151a4700",
      "block_number": 12345678,
      "block_timestamp": 1645564800000,
      "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
      "value": "1000000000",
      "token_id": "1000221",
      "token_name": "Tether USD",
      "token_symbol": "USDT",
      "token_decimals": 6,
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N"
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Transfer by Hash

Retrieves a specific transfer by its transaction hash.

**Endpoint:** `GET /v1/transfers/{hash}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `hash` | string | Transaction hash |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/transfers/70d655a17e04d6b6b7ee5d53e7f37655974f4e71b0edd6bcb311915a151a4700" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "transaction_id": "70d655a17e04d6b6b7ee5d53e7f37655974f4e71b0edd6bcb311915a151a4700",
      "block_number": 12345678,
      "block_timestamp": 1645564800000,
      "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
      "value": "1000000000",
      "token_id": "1000221",
      "token_name": "Tether USD",
      "token_symbol": "USDT",
      "token_decimals": 6,
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N"
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

## Events

### Get Events List

Retrieves a paginated list of contract events.

**Endpoint:** `GET /v1/events`

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `since` | integer | Start timestamp (milliseconds) |
| `start` | integer | Start page (default: 0) |
| `block` | integer | Filter by block number >= this value |
| `contract` | string | Filter by contract address |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/events?limit=2&sort=-timestamp&since=1645564800000&block=1000000" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "block_number": 1000001,
      "block_timestamp": 1645564800000,
      "caller_contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "event_index": "0",
      "event_name": "Transfer",
      "event": "Transfer(address,address,uint256)",
      "transaction_id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "result": {
        "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
        "value": "1000000000"
      },
      "result_type": {
        "from": "address",
        "to": "address",
        "value": "uint256"
      }
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Events by Transaction ID

Retrieves events for a specific transaction.

**Endpoint:** `GET /v1/events/transaction/{transactionId}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `transactionId` | string | Transaction hash |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/events/transaction/cd402e64cad7e69c086649401f6427f5852239f41f51a100abfc7beaa8aa0f9c" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "block_number": 1000001,
      "block_timestamp": 1645564800000,
      "caller_contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "event_index": "0",
      "event_name": "Transfer",
      "event": "Transfer(address,address,uint256)",
      "transaction_id": "cd402e64cad7e69c086649401f6427f5852239f41f51a100abfc7beaa8aa0f9c",
      "result": {
        "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
        "value": "1000000000"
      },
      "result_type": {
        "from": "address",
        "to": "address",
        "value": "uint256"
      }
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Events by Contract Address

Retrieves events for a specific contract address.

**Endpoint:** `GET /v1/events/{contractAddress}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `contractAddress` | string | Contract address (base58 or hex) |

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `since` | integer | Start timestamp (milliseconds) |
| `block` | integer | Filter by block number >= this value |
| `start` | integer | Start page (default: 0) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/events/LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N?limit=2&sort=-timestamp" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "block_number": 1000001,
      "block_timestamp": 1645564800000,
      "caller_contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "event_index": "0",
      "event_name": "Bet",
      "event": "Bet(address,uint256)",
      "transaction_id": "cd402e64cad7e69c086649401f6427f5852239f41f51a100abfc7beaa8aa0f9c",
      "result": {
        "player": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "amount": "1000000"
      },
      "result_type": {
        "player": "address",
        "amount": "uint256"
      }
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Events by Contract Address and Event Name

Retrieves events filtered by contract address and event name.

**Endpoint:** `GET /v1/events/contract/{contractAddress}/{eventName}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `contractAddress` | string | Contract address |
| `eventName` | string | Event name (e.g., Transfer, Approval, Bet) |

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `since` | integer | Start timestamp (milliseconds) |
| `start` | integer | Start page (default: 0) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/events/contract/LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N/Bet?limit=2" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "block_number": 1000001,
      "block_timestamp": 1645564800000,
      "caller_contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "event_index": "0",
      "event_name": "Bet",
      "event": "Bet(address,uint256)",
      "transaction_id": "cd402e64cad7e69c086649401f6427f5852239f41f51a100abfc7beaa8aa0f9c",
      "result": {
        "player": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "amount": "1000000"
      },
      "result_type": {
        "player": "address",
        "amount": "uint256"
      }
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Events by Contract Address, Event Name, and Block Number

Retrieves events for a specific contract, event name, and block number.

**Endpoint:** `GET /v1/events/contract/{contractAddress}/{eventName}/{blockNumber}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `contractAddress` | string | Contract address |
| `eventName` | string | Event name |
| `blockNumber` | integer | Block number |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/events/contract/LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N/Bet/4835773" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "block_number": 4835773,
      "block_timestamp": 1645564800000,
      "caller_contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "event_index": "0",
      "event_name": "Bet",
      "event": "Bet(address,uint256)",
      "transaction_id": "cd402e64cad7e69c086649401f6427f5852239f41f51a100abfc7beaa8aa0f9c",
      "result": {
        "player": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "amount": "1000000"
      },
      "result_type": {
        "player": "address",
        "amount": "uint256"
      }
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Events by Timestamp

Retrieves events filtered by timestamp range.

**Endpoint:** `GET /v1/events/timestamp`

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `since` | integer | Start timestamp (milliseconds) |
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `start` | integer | Start page (default: 0) |
| `contract` | string | Filter by contract address |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/events/timestamp?since=1644483426749&limit=2" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "block_number": 1000001,
      "block_timestamp": 1645564800000,
      "caller_contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "event_index": "0",
      "event_name": "Transfer",
      "event": "Transfer(address,address,uint256)",
      "transaction_id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "result": {
        "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
        "value": "1000000000"
      },
      "result_type": {
        "from": "address",
        "to": "address",
        "value": "uint256"
      }
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Confirmed Events List

Retrieves a list of confirmed (finalized) events.

**Endpoint:** `GET /v1/events/confirmed`

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `since` | integer | Start timestamp (milliseconds) |
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `start` | integer | Start page (default: 0) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/events/confirmed?since=1644483426749&limit=2" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "block_number": 1000001,
      "block_timestamp": 1645564800000,
      "caller_contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "event_index": "0",
      "event_name": "Transfer",
      "event": "Transfer(address,address,uint256)",
      "transaction_id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "result": {
        "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
        "value": "1000000000"
      },
      "result_type": {
        "from": "address",
        "to": "address",
        "value": "uint256"
      }
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

## Blocks

### Get Block by Hash

Retrieves a specific block by its hash.

**Endpoint:** `GET /v1/blocks/{hash}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `hash` | string | Block hash (64 characters hex) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/blocks/000000000049c11f15d4e91e988bc950fa9f194d2cb2e04cda76675dbb349009" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "hash": "000000000049c11f15d4e91e988bc950fa9f194d2cb2e04cda76675dbb349009",
      "number": 1000000,
      "timestamp": 1645564800000,
      "parentHash": "000000000049c11f15d4e91e988bc950fa9f194d2cb2e04cda76675dbb349008",
      "witnessAddress": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "transactionCount": 15
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Block List

Retrieves a paginated list of blocks.

**Endpoint:** `GET /v1/blocks`

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `start` | integer | Start page (default: 0) |
| `block` | integer | Filter by block number >= this value |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/blocks?limit=2&sort=-timestamp" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "hash": "000000000049c11f15d4e91e988bc950fa9f194d2cb2e04cda76675dbb349009",
      "number": 1000001,
      "timestamp": 1645564800000,
      "parentHash": "000000000049c11f15d4e91e988bc950fa9f194d2cb2e04cda76675dbb349008",
      "witnessAddress": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "transactionCount": 15
    },
    {
      "hash": "000000000049c11f15d4e91e988bc950fa9f194d2cb2e04cda76675dbb349008",
      "number": 1000000,
      "timestamp": 1645564740000,
      "parentHash": "000000000049c11f15d4e91e988bc950fa9f194d2cb2e04cda76675dbb349007",
      "witnessAddress": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "transactionCount": 12
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 2,
    "fingerprint": "eyJvZmZzZXQiOjJ9",
    "links": {
      "next": "/v1/blocks?fingerprint=eyJvZmZzZXQiOjJ9&limit=2"
    }
  }
}
```

---

### Get Latest Solidified Block Number

Retrieves the number of the latest solidified (finalized) block.

**Endpoint:** `GET /v1/blocks/latestSolidifiedBlockNumber`

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/blocks/latestSolidifiedBlockNumber" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "number": 1000001
}
```

---

## Contract Logs

### Get Contract Log List

Retrieves a paginated list of contract logs.

**Endpoint:** `GET /v1/contractlogs`

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `limit` | integer | Page size (default: 20, max: 200) |
| `sort` | string | Sort by `timestamp` (default: `-timestamp`) |
| `start` | integer | Start page (default: 0) |
| `block` | integer | Filter by block number >= this value |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/contractlogs?limit=2" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "topics": [
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        "0x000000000000000000000041f0cc5a2a84cd0f68ed1667070934542d673acbd8",
        "0x00000000000000000000004195fd23d3d2221cfef64167938de5e62074719e54"
      ],
      "data": "0x000000000000000000000000000000000000000000000000000000003b9aca00",
      "blockNumber": 1000001,
      "blockTimestamp": 1645564800000,
      "transactionId": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "transactionIndex": 0,
      "logIndex": 0
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Contract Logs by Transaction ID

Retrieves contract logs for a specific transaction.

**Endpoint:** `GET /v1/contractlogs/transaction/{transactionId}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `transactionId` | string | Transaction hash |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/contractlogs/transaction/7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "topics": [
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        "0x000000000000000000000041f0cc5a2a84cd0f68ed1667070934542d673acbd8",
        "0x00000000000000000000004195fd23d3d2221cfef64167938de5e62074719e54"
      ],
      "data": "0x000000000000000000000000000000000000000000000000000000003b9aca00",
      "blockNumber": 1000001,
      "blockTimestamp": 1645564800000,
      "transactionId": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "transactionIndex": 0,
      "logIndex": 0
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Contract Logs by Contract Address

Retrieves contract logs for a specific contract address.

**Endpoint:** `GET /v1/contractlogs/contract/{contractAddress}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `contractAddress` | string | Contract address |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/v1/contractlogs/contract/LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "data": [
    {
      "address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "topics": [
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        "0x000000000000000000000041f0cc5a2a84cd0f68ed1667070934542d673acbd8",
        "0x00000000000000000000004195fd23d3d2221cfef64167938de5e62074719e54"
      ],
      "data": "0x000000000000000000000000000000000000000000000000000000003b9aca00",
      "blockNumber": 1000001,
      "blockTimestamp": 1645564800000,
      "transactionId": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "transactionIndex": 0,
      "logIndex": 0
    }
  ],
  "success": true,
  "meta": {
    "at": 1645564800123,
    "page_size": 1
  }
}
```

---

### Get Contract Logs with ABI (by Transaction ID)

Posts an ABI string to decode contract logs for a specific transaction.

**Endpoint:** `POST /v1/contract/transaction/{transactionId}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `transactionId` | string | Transaction hash |

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `abi` | string | Contract ABI JSON string |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/v1/contract/transaction/7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "abi": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
  }'
```

**Example Response:**

```json
{
  "contractAddress": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
  "abi": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]",
  "logs": [
    {
      "name": "Transfer",
      "params": {
        "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
        "value": "1000000000"
      }
    }
  ]
}
```

---

### Get Contract Logs with ABI (by Contract Address)

Posts an ABI string to decode contract logs for a specific contract address.

**Endpoint:** `POST /v1/contract/contractAddress/{contractAddress}`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `contractAddress` | string | Contract address |

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `abi` | string | Contract ABI JSON string |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/v1/contract/contractAddress/LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "abi": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"
  }'
```

**Example Response:**

```json
{
  "contractAddress": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
  "abi": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]",
  "logs": [
    {
      "name": "Transfer",
      "params": {
        "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
        "value": "1000000000"
      }
    }
  ]
}
```

---

## Error Responses

All endpoints return standardized error responses:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

### Common HTTP Status Codes

| Code | Description |
|------|-------------|
| `200` | Success |
| `400` | Bad Request - Invalid parameters |
| `401` | Unauthorized - Missing or invalid API key |
| `403` | Forbidden - Rate limit exceeded |
| `404` | Not Found - Resource not found |
| `429` | Too Many Requests - Rate limit exceeded |
| `500` | Internal Server Error |

---

## Rate Limiting

The V1 API implements rate limiting to ensure fair usage:

| User Type | Rate Limit |
|-----------|------------|
| With API Key | 15 requests per second |
| Without API Key | 5 requests per second (stricter) |

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Maximum requests per second
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Time when the limit resets (Unix timestamp)

---

## Best Practices

1. **Always include your API key** in the `LINDA-PRO-API-KEY` header for higher rate limits
2. **Use pagination** with `limit` and `fingerprint` for large result sets
3. **Cache responses** when appropriate to reduce API calls
4. **Handle rate limiting** gracefully by implementing retry logic with exponential backoff
5. **Use specific filters** (like `contract`, `from`, `to`) to narrow down results
6. **Monitor rate limit headers** to stay within quotas