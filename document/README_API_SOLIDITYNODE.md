# FullNode Solidity HTTP API

The FullNode Solidity HTTP API provides access to confirmed (finalized) blockchain data. Unlike the regular FullNode API which may return unconfirmed data, the Solidity API guarantees that all returned information has been solidified by the network. This makes it ideal for applications requiring high certainty, such as explorers, wallets displaying balances, and analytics platforms.

## Base URL

All Solidity Node HTTP API endpoints are prefixed with `/walletsolidity/`

```
https://api.lindagrid.lindacoin.org/walletsolidity/
```

## Key Characteristics

| Feature | Description |
|---------|-------------|
| **Data Finality** | All data is confirmed and irreversible |
| **Performance** | Optimized for read-heavy operations |
| **Use Cases** | Explorers, wallets, analytics, dApps |
| **Availability** | Requires a Solidity node to be running |

## Headers

| Header | Description | Required |
|--------|-------------|----------|
| `LINDA-PRO-API-KEY` | Your API key for authentication | Yes (for rate limiting) |
| `Content-Type` | `application/json` | Yes |

## Address Format

The API supports both hex and base58 address formats. Use the `visible` parameter to specify the format:

- `visible: false` (default) - Addresses in hex format (starting with "30")
- `visible: true` - Addresses in base58check format (e.g., "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD")

---

## Account APIs (Confirmed)

### Get Account

Retrieves confirmed account information including balance, resources, and permissions.

**Endpoint:** `POST /walletsolidity/getaccount`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getaccount" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "balance": 1000000000,
  "account_name": "MyAccount",
  "create_time": 1645564800000,
  "is_witness": false,
  "allowance": 50000000,
  "latest_withdraw_time": 1645564800000,
  "latest_opration_time": 1645564800000,
  "latest_consume_time": 1645564800000,
  "latest_consume_free_time": 1645564800000,
  "net_window_size": 28800,
  "net_window_optimized": true,
  "frozen": [
    {
      "frozen_balance": 10000000,
      "expire_time": 1648252800000
    }
  ],
  "delegated_frozen_balance_for_bandwidth": 5000000,
  "acquired_delegated_frozen_balance_for_bandwidth": 2000000,
  "frozenV2": [
    {
      "amount": 5000000,
      "type": "ENERGY"
    }
  ],
  "account_resource": {
    "frozen_balance_for_energy": {
      "frozen_balance": 10000000,
      "expire_time": 1648252800000
    },
    "delegated_frozen_balance_for_energy": 3000000,
    "energy_usage": 50000,
    "energy_window_size": 28800
  },
  "owner_permission": {
    "type": 0,
    "permission_name": "owner",
    "threshold": 1,
    "keys": [
      {
        "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
        "weight": 1
      }
    ]
  },
  "asset": {
    "1000001": 1000
  },
  "assetV2": {
    "1000001": 1000
  },
  "lrc20": [
    {
      "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N": "1000000000"
    }
  ],
  "votes": [
    {
      "vote_address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
      "vote_count": 100
    }
  ],
  "net_usage": 5000,
  "free_net_usage": 1000
}
```

---

### Get Account by ID

Retrieves confirmed account information using account ID.

**Endpoint:** `POST /walletsolidity/getaccountbyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `account_id` | string | Account ID |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getaccountbyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": "MyAccountID123",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "balance": 1000000000,
  "account_name": "MyAccount",
  "create_time": 1645564800000
}
```

---

### Get Account Resource

Retrieves confirmed resource information for an account.

**Endpoint:** `POST /walletsolidity/getaccountresource`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getaccountresource" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "freeNetUsed": 1000,
  "freeNetLimit": 5000,
  "NetUsed": 500,
  "NetLimit": 10000,
  "TotalNetLimit": 43200000000,
  "TotalNetWeight": 1000000000,
  "totalLindaPowerWeight": 500000000,
  "lindaPowerLimit": 10000,
  "lindaPowerUsed": 500,
  "EnergyUsed": 2000,
  "EnergyLimit": 100000,
  "TotalEnergyLimit": 180000000000,
  "TotalEnergyWeight": 900000000,
  "assetNetUsed": {
    "1000001": 200
  },
  "assetNetLimit": {
    "1000001": 5000
  }
}
```

---

### Get Account Net

Retrieves confirmed bandwidth information for an account.

**Endpoint:** `POST /walletsolidity/getaccountnet`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getaccountnet" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "freeNetUsed": 1000,
  "freeNetLimit": 5000,
  "NetUsed": 500,
  "NetLimit": 10000,
  "TotalNetLimit": 43200000000,
  "TotalNetWeight": 1000000000,
  "assetNetUsed": {
    "1000001": 200
  },
  "assetNetLimit": {
    "1000001": 5000
  }
}
```

---

## Transaction APIs (Confirmed)

### Get Transaction by ID

Retrieves a confirmed transaction by its hash.

**Endpoint:** `POST /walletsolidity/gettransactionbyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Transaction hash |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/gettransactionbyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "amount": 1000000,
            "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
            "to_address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr"
          },
          "type_url": "type.googleapis.com/protocol.TransferContract"
        },
        "type": "TransferContract"
      }
    ],
    "ref_block_bytes": "f69b",
    "ref_block_hash": "7d4a3b02495f2320",
    "expiration": 1762502739000,
    "timestamp": 1762502681856
  },
  "raw_data_hex": "0a02f69b22087d4a3b02495f232040b888e6eaa5335a66080112620a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412310a1541e9d79cc47518930bc322d9bf7cddd260a0260a8d12154195fd23d3d2221cfef64167938de5e62074719e5418c0843d70c2cbe6eaa533",
  "signature": [
    "2a3743f40d53a124c1597256b155bf286bd8874afe6997ec0a7e63405dea78cd914d9aa9adb8f84dc8d1ef4b1827dd9e8960a40a4ad11e619d06e9601e8b27c000"
  ],
  "ret": [
    {
      "contractRet": "SUCCESS",
      "fee": 100000
    }
  ]
}
```

---

### Get Transaction Info by ID

Retrieves detailed confirmed information about a transaction including fees and logs.

**Endpoint:** `POST /walletsolidity/gettransactioninfobyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Transaction hash |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/gettransactioninfobyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
  "fee": 100000,
  "blockNumber": 1000001,
  "blockTimeStamp": 1645564800000,
  "contractResult": [
    ""
  ],
  "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
  "receipt": {
    "energy_usage": 50000,
    "energy_fee": 50000,
    "origin_energy_usage": 0,
    "energy_usage_total": 50000,
    "net_usage": 350,
    "net_fee": 50000,
    "result": "SUCCESS"
  },
  "log": [
    {
      "address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
      "topics": [
        "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        "000000000000000000000041f0cc5a2a84cd0f68ed1667070934542d673acbd8",
        "00000000000000000000004195fd23d3d2221cfef64167938de5e62074719e54"
      ],
      "data": "000000000000000000000000000000000000000000000000000000003b9aca00"
    }
  ],
  "result": 0,
  "internal_transactions": []
}
```

---

### Get Transaction Count by Block Number

Retrieves the number of confirmed transactions in a specific block.

**Endpoint:** `POST /walletsolidity/gettransactioncountbyblocknum`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `num` | integer | Block number |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/gettransactioncountbyblocknum" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "num": 1000000
  }'
```

**Example Response:**

```json
{
  "num": 15
}
```

---

### Get Transaction Info by Block Number

Retrieves confirmed transaction information for all transactions in a block.

**Endpoint:** `POST /walletsolidity/gettransactioninfobyblocknum`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `num` | integer | Block number |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/gettransactioninfobyblocknum" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "num": 1000000,
    "visible": true
  }'
```

**Example Response:**

```json
[
  {
    "id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
    "fee": 100000,
    "blockNumber": 1000000,
    "blockTimeStamp": 1645564740000,
    "contractResult": [""],
    "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
    "receipt": {
      "energy_usage": 50000,
      "energy_fee": 50000,
      "net_usage": 350,
      "net_fee": 50000,
      "result": "SUCCESS"
    },
    "result": 0
  }
]
```

---

## Block APIs (Confirmed)

### Get Now Block

Retrieves the most recent confirmed block.

**Endpoint:** `POST /walletsolidity/getnowblock`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getnowblock" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "visible": true
  }'
```

**Example Response:**

```json
{
  "blockID": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
  "block_header": {
    "raw_data": {
      "timestamp": 1645564800000,
      "txTrieRoot": "b1144687a9f5e9a1f1c3d8a9b2f4e5d6c7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2",
      "parentHash": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
      "number": 1000001,
      "witness_id": 1,
      "witness_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "version": 1
    },
    "witness_signature": "b1144687a9f5e9a1f1c3d8a9b2f4e5d6c7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0"
  },
  "transactions": []
}
```

---

### Get Block by Number

Retrieves a specific confirmed block by its number.

**Endpoint:** `POST /walletsolidity/getblockbynum`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `num` | integer | Block number |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getblockbynum" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "num": 1000000,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "blockID": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
  "block_header": {
    "raw_data": {
      "timestamp": 1645564740000,
      "number": 1000000,
      "witness_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD"
    }
  },
  "transactions": []
}
```

---

### Get Block by ID

Retrieves a specific confirmed block by its hash.

**Endpoint:** `POST /walletsolidity/getblockbyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Block hash |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getblockbyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "blockID": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
  "block_header": {
    "raw_data": {
      "number": 1000000
    }
  }
}
```

---

### Get Block by Limit Next

Retrieves confirmed blocks within a range.

**Endpoint:** `POST /walletsolidity/getblockbylimitnext`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `startNum` | integer | Start block number (inclusive) |
| `endNum` | integer | End block number (exclusive) |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getblockbylimitnext" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "startNum": 1000000,
    "endNum": 1000002,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "block": [
    {
      "blockID": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
      "block_header": {
        "raw_data": {
          "number": 1000000
        }
      }
    },
    {
      "blockID": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
      "block_header": {
        "raw_data": {
          "number": 1000001
        }
      }
    }
  ]
}
```

---

### Get Block by Latest Number

Retrieves the latest N confirmed blocks.

**Endpoint:** `POST /walletsolidity/getblockbylatestnum`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `num` | integer | Number of blocks to retrieve |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getblockbylatestnum" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "num": 2,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "block": [
    {
      "blockID": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
      "block_header": {
        "raw_data": {
          "number": 1000001
        }
      }
    },
    {
      "blockID": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
      "block_header": {
        "raw_data": {
          "number": 1000000
        }
      }
    }
  ]
}
```

---

### Get Block

Retrieves a confirmed block by ID or number (flexible endpoint).

**Endpoint:** `POST /walletsolidity/getblock`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id_or_num` | string | Block hash or number |
| `detail` | boolean | Include transaction details (default: false) |
| `visible` | boolean | Address format (default: false) |

**Example Request (by number with details):**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getblock" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "id_or_num": "1000000",
    "detail": true,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "blockID": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
  "block_header": {
    "raw_data": {
      "number": 1000000
    }
  },
  "transactions": [
    {
      "txID": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc"
    }
  ]
}
```

---

## Node Info

### Get Node Info

Retrieves confirmed information about the current node.

**Endpoint:** `POST /walletsolidity/getnodeinfo`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getnodeinfo" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Example Response:**

```json
{
  "beginSyncNum": 1000000,
  "block": "1000001",
  "solidityBlock": "1000000",
  "currentConnectCount": 5,
  "activeConnectCount": 3,
  "passiveConnectCount": 2,
  "totalFlow": 1024000,
  "peerInfoList": [
    {
      "host": "192.168.1.100",
      "port": 18888,
      "lastBlockTime": 1645564800000,
      "score": 100,
      "syncToPeer": true,
      "syncFromPeer": false
    }
  ],
  "configNodeInfo": {
    "codeVersion": "v4.5.0",
    "versionNum": "4.5.0",
    "p2pVersion": "1",
    "listenPort": 18888,
    "discoverEnable": true,
    "activeNodeCount": 3,
    "passiveNodeCount": 2,
    "maxConnectCount": 30,
    "dbVersion": "2.0.0"
  },
  "machineInfo": {
    "threadCount": 8,
    "cpuCount": 4,
    "totalMemory": 16777216,
    "freeMemory": 8388608,
    "cpuRate": 25.5,
    "javaCpuRate": 15.2,
    "processCpuRate": 10.3
  }
}
```

---

## Asset APIs (LRC-10 - Confirmed)

### Get Asset Issue by ID

Retrieves confirmed LRC-10 token information by ID.

**Endpoint:** `POST /walletsolidity/getassetissuebyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | integer | Token ID |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getassetissuebyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": 1000001,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "id": "1000001",
  "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "name": "MyToken",
  "abbr": "MTK",
  "total_supply": 1000000000,
  "lind_num": 1000,
  "num": 1,
  "precision": 6,
  "start_time": 1645564800000,
  "end_time": 1648252800000,
  "description": "My Token Description",
  "url": "https://example.com"
}
```

---

### Get Asset Issue by Name

Retrieves confirmed LRC-10 token information by name.

**Endpoint:** `POST /walletsolidity/getassetissuebyname`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Token name |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getassetissuebyname" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "MyToken",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "id": "1000001",
  "name": "MyToken",
  "abbr": "MTK",
  "total_supply": 1000000000
}
```

---

### Get Asset Issue List

Retrieves all confirmed LRC-10 tokens.

**Endpoint:** `GET /walletsolidity/getassetissuelist`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/walletsolidity/getassetissuelist?visible=true" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "assetIssue": [
    {
      "id": "1000001",
      "name": "MyToken",
      "abbr": "MTK"
    },
    {
      "id": "1000002",
      "name": "AnotherToken",
      "abbr": "ATK"
    }
  ]
}
```

---

### Get Asset Issue List by Name

Retrieves all confirmed LRC-10 tokens with a given name.

**Endpoint:** `POST /walletsolidity/getassetissuelistbyname`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Token name |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getassetissuelistbyname" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "MyToken",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "assetIssue": [
    {
      "id": "1000001",
      "name": "MyToken"
    },
    {
      "id": "1000003",
      "name": "MyToken"
    }
  ]
}
```

---

### Get Paginated Asset Issue List

Retrieves paginated list of confirmed LRC-10 tokens.

**Endpoint:** `POST /walletsolidity/getpaginatedassetissuelist`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `offset` | integer | Pagination offset |
| `limit` | integer | Number of items per page |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getpaginatedassetissuelist" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "offset": 0,
    "limit": 10,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "assetIssue": [
    {
      "id": "1000001",
      "name": "MyToken"
    }
  ]
}
```

---

## Exchange APIs (Confirmed)

### Get Exchange by ID

Retrieves confirmed exchange information by ID.

**Endpoint:** `POST /walletsolidity/getexchangebyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `exchange_id` | integer | Exchange ID |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getexchangebyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "exchange_id": 1,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "exchange_id": 1,
  "creator_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "create_time": 1645564800000,
  "first_token_id": "1000001",
  "first_token_balance": 1000000,
  "second_token_id": "_",
  "second_token_balance": 1000000000
}
```

---

### List Exchanges

Retrieves all confirmed exchanges.

**Endpoint:** `GET /walletsolidity/listexchanges`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/walletsolidity/listexchanges?visible=true" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "exchanges": [
    {
      "exchange_id": 1,
      "creator_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD"
    }
  ]
}
```

---

## Market APIs (Confirmed)

### Get Market Order by Account

Retrieves confirmed market orders for an account.

**Endpoint:** `POST /walletsolidity/getmarketorderbyaccount`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getmarketorderbyaccount" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "orders": [
    {
      "order_id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "create_time": 1645564800000,
      "sell_token_id": "1000001",
      "sell_token_value": 1000,
      "buy_token_id": "_",
      "buy_token_value": 1000000,
      "order_status": "COMPLETED"
    }
  ]
}
```

---

### Get Market Order by ID

Retrieves a confirmed market order by its ID.

**Endpoint:** `POST /walletsolidity/getmarketorderbyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `order_id` | string | Order ID |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getmarketorderbyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "order_id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
  "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "create_time": 1645564800000,
  "sell_token_id": "1000001",
  "sell_token_value": 1000,
  "buy_token_id": "_",
  "buy_token_value": 1000000,
  "order_status": "COMPLETED"
}
```

---

### Get Market Price by Pair

Retrieves confirmed market price for a trading pair.

**Endpoint:** `POST /walletsolidity/getmarketpricebypair`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `sell_token_id` | string | Sell token ID |
| `buy_token_id` | string | Buy token ID |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getmarketpricebypair" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "sell_token_id": "1000001",
    "buy_token_id": "_",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "prices": [
    {
      "sell_token_id": "1000001",
      "buy_token_id": "_",
      "price": "1000"
    }
  ]
}
```

---

### Get Market Order List by Pair

Retrieves confirmed market orders for a trading pair.

**Endpoint:** `POST /walletsolidity/getmarketorderlistbypair`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `sell_token_id` | string | Sell token ID |
| `buy_token_id` | string | Buy token ID |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getmarketorderlistbypair" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "sell_token_id": "1000001",
    "buy_token_id": "_",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "orders": [
    {
      "order_id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "sell_token_id": "1000001",
      "sell_token_value": 1000
    }
  ]
}
```

---

### Get Market Pair List

Retrieves all confirmed market pairs.

**Endpoint:** `GET /walletsolidity/getmarketpairlist`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/walletsolidity/getmarketpairlist?visible=true" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "pairs": [
    {
      "sell_token_id": "1000001",
      "buy_token_id": "_"
    }
  ]
}
```

---

## Smart Contract APIs (Confirmed)

### Trigger Constant Contract

Triggers a constant/simulated contract call (does not create a transaction).

**Endpoint:** `POST /walletsolidity/triggerconstantcontract`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Caller address |
| `contract_address` | string | Contract address |
| `function_selector` | string | Function signature (e.g., "balanceOf(address)") |
| `parameter` | string | ABI-encoded parameters (hex) |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/triggerconstantcontract" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
    "function_selector": "balanceOf(address)",
    "parameter": "000000000000000000000041f0cc5a2a84cd0f68ed1667070934542d673acbd8",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "result": {
    "result": true
  },
  "energy_used": 1000,
  "constant_result": [
    "000000000000000000000000000000000000000000000000000000003b9aca00"
  ],
  "transaction": {
    "txID": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc"
  }
}
```

---

### Estimate Energy

Estimates energy required for a contract call.

**Endpoint:** `POST /walletsolidity/estimateenergy`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Caller address |
| `contract_address` | string | Contract address |
| `function_selector` | string | Function signature |
| `parameter` | string | ABI-encoded parameters (hex) |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/estimateenergy" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
    "function_selector": "transfer(address,uint256)",
    "parameter": "00000000000000000000004195fd23d3d2221cfef64167938de5e62074719e54000000000000000000000000000000000000000000000000000000003b9aca00",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "result": {
    "result": true
  },
  "energy_required": 50000
}
```

---

### Get Contract

Retrieves confirmed contract information.

**Endpoint:** `POST /walletsolidity/getcontract`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Contract address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getcontract" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "origin_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
  "abi": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"type\":\"function\"}]",
  "bytecode": "608060405234801561001057600080fd5b5060de8061001f6000396000f3...",
  "name": "MyToken",
  "consume_user_resource_percent": 100,
  "origin_energy_limit": 10000000,
  "code_hash": "b1144687a9f5e9a1f1c3d8a9b2f4e5d6c7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2"
}
```

---

### Get Contract Info

Retrieves detailed confirmed contract information including runtime code.

**Endpoint:** `POST /walletsolidity/getcontractinfo`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Contract address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getcontractinfo" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "runtimecode": "6080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde03811461005b575b600080fd5b34801561006757600080fd5b50610070610086565b604080516020808252918101919091520190565b60606040518060400160405280600781526020017f4d79546f6b656e000000000000000000000000000000000000000000000000815250905090565b600080fd00a165627a7a72305820a2fb39541e90eda9a2f5f9e7905ef98e66e60dd4b38e00b05de418da3154e7570029",
  "smart_contract": {
    "origin_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N",
    "abi": "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"type\":\"function\"}]",
    "bytecode": "608060405234801561001057600080fd5b5060de8061001f6000396000f3..."
  },
  "contract_state": {
    "energy_usage": 1000,
    "energy_factor": 100,
    "update_cycle": 1000000
  }
}
```

---

## Witness APIs (Confirmed)

### List Witnesses

Retrieves all confirmed witnesses (Super Representatives).

**Endpoint:** `GET /walletsolidity/listwitnesses`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/walletsolidity/listwitnesses?visible=true" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "witnesses": [
    {
      "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "voteCount": 10000000,
      "url": "https://example.com",
      "totalProduced": 10000,
      "totalMissed": 10,
      "latestBlockNum": 1000000,
      "latestSlotNum": 100,
      "isJobs": true
    }
  ]
}
```

---

### Get Brokerage

Retrieves confirmed brokerage ratio for a Super Representative.

**Endpoint:** `POST /walletsolidity/getBrokerage`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | SR address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getBrokerage" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "brokerage": 20
}
```

---

### Get Reward

Retrieves confirmed withdrawable rewards for an account.

**Endpoint:** `POST /walletsolidity/getReward`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getReward" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "reward": 5000000
}
```

---

### Get Paginated Now Witness List

Retrieves paginated list of confirmed witnesses with current vote counts.

**Endpoint:** `GET /soliditywallet/getpaginatednowwitnesslist`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `offset` | integer | Pagination offset |
| `limit` | integer | Number of items per page |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/soliditywallet/getpaginatednowwitnesslist?offset=0&limit=10&visible=true" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "witnesses": [
    {
      "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "voteCount": 10000000
    }
  ]
}
```

---

## Resource Delegation APIs (Confirmed)

### Get Delegated Resource

Retrieves confirmed resource delegations (Stake 1.0).

**Endpoint:** `POST /walletsolidity/getdelegatedresource`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `fromAddress` | string | Delegator address |
| `toAddress` | string | Delegatee address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getdelegatedresource" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "fromAddress": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "toAddress": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "delegatedResource": [
    {
      "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
      "frozen_balance_for_bandwidth": 10000000,
      "frozen_balance_for_energy": 5000000,
      "expire_time_for_bandwidth": 1648252800000,
      "expire_time_for_energy": 1648252800000
    }
  ]
}
```

---

### Get Delegated Resource V2

Retrieves confirmed resource delegations (Stake 2.0).

**Endpoint:** `POST /walletsolidity/getdelegatedresourcev2`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `fromAddress` | string | Delegator address |
| `toAddress` | string | Delegatee address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getdelegatedresourcev2" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "fromAddress": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "toAddress": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "delegatedResource": [
    {
      "from": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "to": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
      "frozen_balance_for_bandwidth": 10000000,
      "frozen_balance_for_energy": 5000000,
      "expire_time_for_bandwidth": 1648252800000,
      "expire_time_for_energy": 1648252800000
    }
  ]
}
```

---

### Get Delegated Resource Account Index

Retrieves confirmed delegation index for an account (Stake 1.0).

**Endpoint:** `POST /walletsolidity/getdelegatedresourceaccountindex`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getdelegatedresourceaccountindex" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "account": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "fromAccounts": [
    "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr"
  ],
  "toAccounts": [
    "LdaqF7MRaNe2iGy5d3n3qajvuqPLrYWxfo"
  ]
}
```

---

### Get Delegated Resource Account Index V2

Retrieves confirmed delegation index for an account (Stake 2.0).

**Endpoint:** `POST /walletsolidity/getdelegatedresourceaccountindexv2`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getdelegatedresourceaccountindexv2" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "account": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "fromAccounts": [
    "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr"
  ],
  "toAccounts": [
    "LdaqF7MRaNe2iGy5d3n3qajvuqPLrYWxfo"
  ]
}
```

---

### Get Can Delegated Max Size

Retrieves maximum delegatable amount for an account (Stake 2.0).

**Endpoint:** `POST /walletsolidity/getcandelegatedmaxsize`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `type` | integer | Resource type (0: bandwidth, 1: energy) |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getcandelegatedmaxsize" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "type": 1,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "max_size": 10000000
}
```

---

### Get Available Unfreeze Count

Retrieves remaining unstake operations count (Stake 2.0).

**Endpoint:** `POST /walletsolidity/getavailableunfreezecount`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getavailableunfreezecount" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "count": 32
}
```

---

### Get Can Withdraw Unfreeze Amount

Retrieves withdrawable amount at a specific timestamp (Stake 2.0).

**Endpoint:** `POST /walletsolidity/getcanwithdrawunfreezeamount`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `timestamp` | integer | Timestamp in milliseconds |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getcanwithdrawunfreezeamount" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "timestamp": 1648252800000,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "amount": 5000000
}
```

---

## Shielded Transaction APIs (Confirmed)

### Scan Note by Ivk

Scans for notes using an incoming viewing key.

**Endpoint:** `POST /walletsolidity/scannotebyivk`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ivk` | string | Incoming viewing key |
| `start_block_index` | string | Start block |
| `end_block_index` | string | End block |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/scannotebyivk" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "ivk": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "start_block_index": "1000000",
    "end_block_index": "1001000"
  }'
```

**Example Response:**

```json
{
  "notes": ["note1", "note2"],
  "block_number": 1000000,
  "block_timestamp": 1645564800000
}
```

---

### Scan and Mark Note by Ivk

Scans and marks notes using an incoming viewing key.

**Endpoint:** `POST /walletsolidity/scanandmarknotebyivk`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ivk` | string | Incoming viewing key |
| `start_block_index` | string | Start block |
| `end_block_index` | string | End block |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/scanandmarknotebyivk" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "ivk": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "start_block_index": "1000000",
    "end_block_index": "1001000"
  }'
```

**Example Response:**

```json
{
  "notes": ["note1", "note2"],
  "block_number": 1000000,
  "block_timestamp": 1645564800000
}
```

---

### Scan Note by Ovk

Scans for notes using an outgoing viewing key.

**Endpoint:** `POST /walletsolidity/scannotebyovk`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `ovk` | string | Outgoing viewing key |
| `start_block_index` | string | Start block |
| `end_block_index` | string | End block |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/scannotebyovk" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "ovk": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "start_block_index": "1000000",
    "end_block_index": "1001000"
  }'
```

**Example Response:**

```json
{
  "notes": ["note1", "note2"],
  "block_number": 1000000,
  "block_timestamp": 1645564800000
}
```

---

### Get Merkle Tree Voucher Info

Retrieves Merkle tree voucher information.

**Endpoint:** `POST /walletsolidity/getmerkletreevoucherinfo`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `notes` | array | List of notes |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/getmerkletreevoucherinfo" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": ["note1", "note2"]
  }'
```

**Example Response:**

```json
{
  "vouchers": ["voucher1", "voucher2"]
}
```

---

### Is Spend

Checks if a note has been spent.

**Endpoint:** `POST /walletsolidity/isspend`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `nullifier` | string | Nullifier |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/walletsolidity/isspend" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "nullifier": "1234567890abcdef1234567890abcdef"
  }'
```

**Example Response:**

```json
{
  "is_spend": false
}
```

---

## Burn LIND

### Get Burn LIND

Retrieves the amount of LIND burned from transaction fees.

**Endpoint:** `GET /walletsolidity/getburnlind`

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/walletsolidity/getburnlind" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "burnTrxAmount": 1000000000
}
```

---

## Error Responses

All endpoints return standardized error responses:

```json
{
  "Error": "Error message describing what went wrong"
}
```

Or with more detail:

```json
{
  "result": false,
  "code": "ERROR_CODE",
  "message": "Error message"
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| `SUCCESS` | Operation successful |
| `SIGERROR` | Signature error |
| `CONTRACT_VALIDATE_ERROR` | Contract validation failed |
| `CONTRACT_EXE_ERROR` | Contract execution failed |
| `BANDWITH_ERROR` | Insufficient bandwidth |
| `DUP_TRANSACTION_ERROR` | Duplicate transaction |
| `TAPOS_ERROR` | TAPOS validation failed |
| `TOO_BIG_TRANSACTION_ERROR` | Transaction too large |
| `TRANSACTION_EXPIRATION_ERROR` | Transaction expired |
| `SERVER_BUSY` | Server is busy |

---

## Best Practices

1. **Use for confirmed data only** - This API is optimized for finalized data
2. **Cache responses** - Data is immutable once confirmed
3. **Use visible=true** for base58 addresses to make responses human-readable
4. **Monitor block progression** to know when data is confirmed
5. **Combine with FullNode API** for real-time + confirmed data needs
6. **Use pagination** for large result sets
7. **Handle errors gracefully** with appropriate retry logic