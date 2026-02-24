# FullNode HTTP API

The FullNode HTTP API provides real-time interaction with the Linda blockchain. These endpoints allow you to query account information, create and broadcast transactions, interact with smart contracts, and access live blockchain data directly from a Linda FullNode.

## Base URL

All FullNode HTTP API endpoints are prefixed with `/wallet/`

```
https://api.lindagrid.lindacoin.org/wallet/
```

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

## Account APIs

### Get Account

Retrieves account information including balance, resources, and permissions.

**Endpoint:** `POST /wallet/getaccount`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getaccount" \
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

### Get Account Balance

Retrieves the historical balance of an account at a specific block.

**Endpoint:** `POST /wallet/getaccountbalance`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `account_identifier.address` | string | Account address |
| `block_identifier.hash` | string | Block hash |
| `block_identifier.number` | integer | Block number |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getaccountbalance" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "account_identifier": {
      "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD"
    },
    "block_identifier": {
      "hash": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
      "number": 1000000
    },
    "visible": true
  }'
```

**Example Response:**

```json
{
  "balance": 950000000,
  "block_identifier": {
    "hash": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
    "number": 1000000
  }
}
```

---

### Get Account Resource

Retrieves resource information (bandwidth, energy, etc.) for an account.

**Endpoint:** `POST /wallet/getaccountresource`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getaccountresource" \
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

Retrieves bandwidth information for an account.

**Endpoint:** `POST /wallet/getaccountnet`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getaccountnet" \
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

### Create Account

Activates a new account on the blockchain.

**Endpoint:** `POST /wallet/createaccount`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Transaction initiator address |
| `account_address` | string | Address of the account to activate |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/createaccount" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "account_address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
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
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "account_address": "301320a6fb4dcd4ff8e91392a8cb98378633cf7dd8",
            "type": 0
          },
          "type_url": "type.googleapis.com/protocol.AccountCreateContract"
        },
        "type": "AccountCreateContract"
      }
    ],
    "ref_block_bytes": "f69b",
    "ref_block_hash": "7d4a3b02495f2320",
    "expiration": 1762502739000,
    "timestamp": 1762502681856
  },
  "raw_data_hex": "0a02f69b22087d4a3b02495f232040b888e6eaa5335a6f080d126b0a32747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e4163636f756e74437265617465436f6e747261637412350a1541e9d79cc47518930bc322d9bf7cddd260a0260a8d1215411320a6fb4dcd4ff8e91392a8cb98378633cf7dd8180070c2cbe6eaa533",
  "signature": []
}
```

---

### Update Account

Updates the name of an account.

**Endpoint:** `POST /wallet/updateaccount`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `account_name` | string | New account name |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/updateaccount" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "account_name": "NewAccountName",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "8dd26d1772231569f022adb42f7d7161dee88b97b4b35eeef6ce73fcd6613bc2",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "account_name": "4e65774163636f756e744e616d65"
          },
          "type_url": "type.googleapis.com/protocol.AccountUpdateContract"
        },
        "type": "AccountUpdateContract"
      }
    ],
    "ref_block_bytes": "f69c",
    "ref_block_hash": "7d4a3b02495f2321",
    "expiration": 1762502739000,
    "timestamp": 1762502681856
  },
  "signature": []
}
```

---

### Account Permission Update

Updates the permissions of an account (multi-signature setup).

**Endpoint:** `POST /wallet/accountpermissionupdate`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `owner` | object | Owner permission |
| `witness` | object | Witness permission (for SRs) |
| `actives` | array | List of active permissions |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/accountpermissionupdate" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "owner": {
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
    "actives": [
      {
        "type": 2,
        "permission_name": "active0",
        "threshold": 2,
        "operations": "7fff1fc0037e0000000000000000000000000000000000000000000000000000",
        "keys": [
          {
            "address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
            "weight": 1
          },
          {
            "address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
            "weight": 1
          }
        ]
      }
    ],
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "9dd37f288374a5b7d1c9a5b7e8f9a0b1c2d3e4f5g6h7i8j9k0l1m2n3o4p5q6r",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "owner": {
              "type": 0,
              "permission_name": "owner",
              "threshold": 1,
              "keys": [
                {
                  "address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
                  "weight": 1
                }
              ]
            },
            "actives": [
              {
                "type": 2,
                "permission_name": "active0",
                "threshold": 2,
                "operations": "7fff1fc0037e0000000000000000000000000000000000000000000000000000",
                "keys": [
                  {
                    "address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
                    "weight": 1
                  },
                  {
                    "address": "3095fd23d3d2221cfef64167938de5e62074719e54",
                    "weight": 1
                  }
                ]
              }
            ]
          },
          "type_url": "type.googleapis.com/protocol.AccountPermissionUpdateContract"
        },
        "type": "AccountPermissionUpdateContract"
      }
    ],
    "ref_block_bytes": "f69d",
    "ref_block_hash": "7d4a3b02495f2322",
    "expiration": 1762502739000,
    "timestamp": 1762502681856
  },
  "signature": []
}
```

---

## Transaction APIs

### Create Transaction

Creates a LIND transfer transaction.

**Endpoint:** `POST /wallet/createtransaction`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Sender address |
| `to_address` | string | Recipient address |
| `amount` | integer | Amount in SUN (1 LIND = 1,000,000 SUN) |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |
| `extra_data` | string | Optional memo (hex format) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/createtransaction" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "to_address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
    "amount": 1000000,
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
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "to_address": "3095fd23d3d2221cfef64167938de5e62074719e54"
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
  "signature": []
}
```

---

### Broadcast Transaction

Broadcasts a signed transaction to the network.

**Endpoint:** `POST /wallet/broadcasttransaction`

**Request Body:**

The complete signed transaction object (including signatures).

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/broadcasttransaction" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "txID": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
    "raw_data": {
      "contract": [
        {
          "parameter": {
            "value": {
              "amount": 1000000,
              "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
              "to_address": "3095fd23d3d2221cfef64167938de5e62074719e54"
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
    ]
  }'
```

**Example Response:**

```json
{
  "result": true,
  "code": "SUCCESS",
  "txid": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
  "message": ""
}
```

---

### Broadcast Hex

Broadcasts a signed transaction from a hex string.

**Endpoint:** `POST /wallet/broadcasthex`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `transaction` | string | Signed transaction in hex format |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/broadcasthex" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "transaction": "0a8a010a0202db2208c89d4811359a28004098a4e0a6b52d5a730802126f0a32747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e736665724173736574436f6e747261637412390a07313030303030311215415a523b449890854c8fc460ab602df9f31fe4293f1a15416b0580da195542ddabe288fec436c7d5af769d24206412418bf3f2e492ed443607910ea9ef0a7ef79728daaaac0ee2ba6cb87da38366df9ac4ade54b2912c1deb0ee6666b86a07a6c7df68f1f9da171eee6a370b3ca9cbbb00"
  }'
```

**Example Response:**

```json
{
  "result": true,
  "code": "SUCCESS",
  "txid": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
  "message": ""
}
```

---

### Get Transaction by ID

Retrieves a transaction by its hash.

**Endpoint:** `POST /wallet/gettransactionbyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Transaction hash |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/gettransactionbyid" \
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

Retrieves detailed information about a transaction including fees and logs.

**Endpoint:** `POST /wallet/gettransactioninfobyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Transaction hash |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/gettransactioninfobyid" \
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

### Get Transaction Receipt by ID

Retrieves the receipt of a transaction.

**Endpoint:** `POST /wallet/gettransactionreceiptbyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Transaction hash |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/gettransactionreceiptbyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc"
  }'
```

**Example Response:**

```json
{
  "id": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
  "fee": 100000,
  "blockNumber": 1000001,
  "blockTimeStamp": 1645564800000,
  "receipt": {
    "energy_usage": 50000,
    "energy_fee": 50000,
    "origin_energy_usage": 0,
    "energy_usage_total": 50000,
    "net_usage": 350,
    "net_fee": 50000,
    "result": "SUCCESS"
  },
  "contract_address": "LRaZ9mjCUjbT6EBFTKXV8xb8peMCb87k1N"
}
```

---

### Get Transaction Count by Block Number

Retrieves the number of transactions in a specific block.

**Endpoint:** `POST /wallet/gettransactioncountbyblocknum`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `num` | integer | Block number |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/gettransactioncountbyblocknum" \
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

Retrieves transaction information for all transactions in a block.

**Endpoint:** `POST /wallet/gettransactioninfobyblocknum`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `num` | integer | Block number |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/gettransactioninfobyblocknum" \
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

## Block APIs

### Get Now Block

Retrieves the most recent block.

**Endpoint:** `POST /wallet/getnowblock`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getnowblock" \
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

Retrieves a specific block by its number.

**Endpoint:** `POST /wallet/getblockbynum`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `num` | integer | Block number |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getblockbynum" \
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
      "txTrieRoot": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
      "parentHash": "00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc59",
      "number": 1000000,
      "witness_id": 1,
      "witness_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "version": 1
    },
    "witness_signature": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c"
  },
  "transactions": [
    {
      "txID": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "raw_data": {
        "contract": [...]
      }
    }
  ]
}
```

---

### Get Block by ID

Retrieves a specific block by its hash.

**Endpoint:** `POST /wallet/getblockbyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Block hash |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getblockbyid" \
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
      "timestamp": 1645564740000,
      "number": 1000000
    }
  }
}
```

---

### Get Block by Limit Next

Retrieves blocks within a range.

**Endpoint:** `POST /wallet/getblockbylimitnext`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `startNum` | integer | Start block number (inclusive) |
| `endNum` | integer | End block number (exclusive) |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getblockbylimitnext" \
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

Retrieves the latest N blocks.

**Endpoint:** `POST /wallet/getblockbylatestnum`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `num` | integer | Number of blocks to retrieve |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getblockbylatestnum" \
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

Retrieves a block by ID or number (flexible endpoint).

**Endpoint:** `POST /wallet/getblock`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id_or_num` | string | Block hash or number |
| `detail` | boolean | Include transaction details (default: false) |
| `visible` | boolean | Address format (default: false) |

**Example Request (by number with details):**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getblock" \
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
      "txID": "7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "raw_data": {
        "contract": [...]
      }
    }
  ]
}
```

---

## Node APIs

### List Nodes

Retrieves the list of peers connected to the node.

**Endpoint:** `POST /wallet/listnodes`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/listnodes" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Example Response:**

```json
{
  "nodes": [
    {
      "address": {
        "host": "192.168.1.100",
        "port": 18888
      }
    },
    {
      "address": {
        "host": "192.168.1.101",
        "port": 18888
      }
    }
  ]
}
```

---

### Get Node Info

Retrieves information about the current node.

**Endpoint:** `POST /wallet/getnodeinfo`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getnodeinfo" \
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

## Asset (LRC-10) APIs

### Get Asset Issue by Account

Retrieves LRC-10 tokens issued by an account.

**Endpoint:** `POST /wallet/getassetissuebyaccount`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `address` | string | Account address |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getassetissuebyaccount" \
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
  "assetIssue": [
    {
      "id": "1000001",
      "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
      "name": "MyToken",
      "abbr": "MTK",
      "total_supply": 1000000000,
      "frozen_supply": [
        {
          "frozen_amount": 100000000,
          "frozen_days": 30
        }
      ],
      "lind_num": 1000,
      "num": 1,
      "precision": 6,
      "start_time": 1645564800000,
      "end_time": 1648252800000,
      "vote_score": 0,
      "description": "My Token Description",
      "url": "https://example.com",
      "free_asset_net_limit": 10000,
      "public_free_asset_net_limit": 50000
    }
  ]
}
```

---

### Get Asset Issue by ID

Retrieves an LRC-10 token by its ID.

**Endpoint:** `POST /wallet/getassetissuebyid`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Token ID |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getassetissuebyid" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "value": "1000001",
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

Retrieves an LRC-10 token by its name (may return multiple if names are duplicated).

**Endpoint:** `POST /wallet/getassetissuebyname`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `value` | string | Token name |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/getassetissuebyname" \
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
  "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
  "name": "MyToken",
  "abbr": "MTK",
  "total_supply": 1000000000,
  "lind_num": 1000,
  "num": 1,
  "start_time": 1645564800000,
  "end_time": 1648252800000,
  "description": "My Token Description",
  "url": "https://example.com"
}
```

---

### Get Asset Issue List

Retrieves all LRC-10 tokens.

**Endpoint:** `GET /wallet/getassetissuelist`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X GET "https://api.lindagrid.lindacoin.org/wallet/getassetissuelist?visible=true" \
  -H "LINDA-PRO-API-KEY: your-api-key"
```

**Example Response:**

```json
{
  "assetIssue": [
    {
      "id": "1000001",
      "name": "MyToken",
      "symbol": "MTK"
    },
    {
      "id": "1000002",
      "name": "AnotherToken",
      "symbol": "ATK"
    }
  ]
}
```

---

### Create Asset Issue

Creates a new LRC-10 token.

**Endpoint:** `POST /wallet/createassetissue`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Issuer address |
| `name` | string | Token name |
| `abbr` | string | Token symbol |
| `total_supply` | integer | Total supply |
| `lind_num` | integer | Price ratio numerator |
| `num` | integer | Price ratio denominator |
| `start_time` | integer | ICO start time (ms) |
| `end_time` | integer | ICO end time (ms) |
| `description` | string | Token description |
| `url` | string | Project URL |
| `free_asset_net_limit` | integer | Free bandwidth limit per account |
| `public_free_asset_net_limit` | integer | Total free bandwidth limit |
| `frozen_supply` | array | Frozen supply configuration |
| `precision` | integer | Token decimals |
| `visible` | boolean | Address format (default: false) |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/createassetissue" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "name": "MyNewToken",
    "abbr": "MNT",
    "total_supply": 1000000000,
    "lind_num": 1000,
    "num": 1,
    "start_time": 1645564800000,
    "end_time": 1648252800000,
    "description": "My New Token Description",
    "url": "https://example.com",
    "free_asset_net_limit": 10000,
    "public_free_asset_net_limit": 50000,
    "precision": 6,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "ae02a80abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634cef9",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "name": "4d794e6577546f6b656e",
            "abbr": "4d4e54",
            "total_supply": 1000000000,
            "lind_num": 1000,
            "num": 1,
            "start_time": 1645564800000,
            "end_time": 1648252800000,
            "description": "4d79204e657720546f6b656e204465736372697074696f6e",
            "url": "68747470733a2f2f6578616d706c652e636f6d",
            "free_asset_net_limit": 10000,
            "public_free_asset_net_limit": 50000,
            "precision": 6
          },
          "type_url": "type.googleapis.com/protocol.AssetIssueContract"
        },
        "type": "AssetIssueContract"
      }
    ]
  }
}
```

---

## Resource Staking APIs (Stake 1.0)

### Freeze Balance

Stakes LIND to obtain bandwidth or energy (Stake 1.0).

**Endpoint:** `POST /wallet/freezebalance`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `frozen_balance` | integer | Amount to stake (SUN) |
| `frozen_duration` | integer | Lock-up duration (days, default: 3) |
| `resource` | string | Resource type: "BANDWIDTH" or "ENERGY" |
| `receiver_address` | string | Resource receiver address (optional) |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/freezebalance" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "frozen_balance": 1000000000,
    "frozen_duration": 3,
    "resource": "ENERGY",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "bf03a80abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634cefa",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "frozen_balance": 1000000000,
            "frozen_duration": 3,
            "resource": "ENERGY"
          },
          "type_url": "type.googleapis.com/protocol.FreezeBalanceContract"
        },
        "type": "FreezeBalanceContract"
      }
    ]
  }
}
```

---

### Unfreeze Balance

Unstakes LIND to release bandwidth or energy (Stake 1.0).

**Endpoint:** `POST /wallet/unfreezebalance`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `resource` | string | Resource type: "BANDWIDTH" or "ENERGY" |
| `receiver_address` | string | Resource receiver address (optional) |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/unfreezebalance" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "resource": "ENERGY",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "cg14b80abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634cefb",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "resource": "ENERGY"
          },
          "type_url": "type.googleapis.com/protocol.UnfreezeBalanceContract"
        },
        "type": "UnfreezeBalanceContract"
      }
    ]
  }
}
```

---

### Withdraw Balance

Withdraws voting or block production rewards.

**Endpoint:** `POST /wallet/withdrawbalance`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/withdrawbalance" \
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
  "txID": "dh25c80abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634cefc",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d"
          },
          "type_url": "type.googleapis.com/protocol.WithdrawBalanceContract"
        },
        "type": "WithdrawBalanceContract"
      }
    ]
  }
}
```

---

## Resource Staking APIs (Stake 2.0)

### Freeze Balance V2

Stakes LIND to obtain bandwidth or energy (Stake 2.0).

**Endpoint:** `POST /wallet/freezebalancev2`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `frozen_balance` | integer | Amount to stake (SUN) |
| `resource` | string | Resource type: "BANDWIDTH" or "ENERGY" |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/freezebalancev2" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "frozen_balance": 1000000000,
    "resource": "ENERGY",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "ei36d90abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634cefd",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "frozen_balance": 1000000000,
            "resource": "ENERGY"
          },
          "type_url": "type.googleapis.com/protocol.FreezeBalanceV2Contract"
        },
        "type": "FreezeBalanceV2Contract"
      }
    ]
  }
}
```

---

### Unfreeze Balance V2

Unstakes LIND to release bandwidth or energy (Stake 2.0).

**Endpoint:** `POST /wallet/unfreezebalancev2`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `unfreeze_balance` | integer | Amount to unstake (SUN) |
| `resource` | string | Resource type: "BANDWIDTH" or "ENERGY" |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/unfreezebalancev2" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "unfreeze_balance": 500000000,
    "resource": "ENERGY",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "fj47e10abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634cefe",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "unfreeze_balance": 500000000,
            "resource": "ENERGY"
          },
          "type_url": "type.googleapis.com/protocol.UnfreezeBalanceV2Contract"
        },
        "type": "UnfreezeBalanceV2Contract"
      }
    ]
  }
}
```

---

### Delegate Resource

Delegates bandwidth or energy to another account (Stake 2.0).

**Endpoint:** `POST /wallet/delegateresource`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Delegator address |
| `receiver_address` | string | Receiver address |
| `balance` | integer | Amount to delegate (SUN) |
| `resource` | string | Resource type: "BANDWIDTH" or "ENERGY" |
| `lock` | boolean | Whether to lock the delegation |
| `lock_period` | integer | Lock period in blocks |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/delegateresource" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "receiver_address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
    "balance": 200000000,
    "resource": "ENERGY",
    "lock": false,
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "gk58f20abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634ceff",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "receiver_address": "3095fd23d3d2221cfef64167938de5e62074719e54",
            "balance": 200000000,
            "resource": "ENERGY",
            "lock": false
          },
          "type_url": "type.googleapis.com/protocol.DelegateResourceContract"
        },
        "type": "DelegateResourceContract"
      }
    ]
  }
}
```

---

### UnDelegate Resource

Cancels resource delegation (Stake 2.0).

**Endpoint:** `POST /wallet/undelegateresource`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Delegator address |
| `receiver_address` | string | Receiver address |
| `balance` | integer | Amount to undelegate (SUN) |
| `resource` | string | Resource type: "BANDWIDTH" or "ENERGY" |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/undelegateresource" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_address": "LXo7bagwm5CCh3EZBnGP6Twxz7e3FCr4xD",
    "receiver_address": "LeudxtcgoduEyhFTuNzquusYk6yuM73iRr",
    "balance": 200000000,
    "resource": "ENERGY",
    "visible": true
  }'
```

**Example Response:**

```json
{
  "txID": "hl69g30abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634d000",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d",
            "receiver_address": "3095fd23d3d2221cfef64167938de5e62074719e54",
            "balance": 200000000,
            "resource": "ENERGY"
          },
          "type_url": "type.googleapis.com/protocol.UnDelegateResourceContract"
        },
        "type": "UnDelegateResourceContract"
      }
    ]
  }
}
```

---

### Withdraw Expire Unfreeze

Withdraws unfrozen balance after the waiting period (Stake 2.0).

**Endpoint:** `POST /wallet/withdrawexpireunfreeze`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/withdrawexpireunfreeze" \
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
  "txID": "im70h40abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634d001",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d"
          },
          "type_url": "type.googleapis.com/protocol.WithdrawExpireUnfreezeContract"
        },
        "type": "WithdrawExpireUnfreezeContract"
      }
    ]
  }
}
```

---

### Cancel All Unfreeze V2

Cancels all unstaking operations (Stake 2.0).

**Endpoint:** `POST /wallet/cancelallunfreezev2`

**Request Body:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner_address` | string | Account address |
| `visible` | boolean | Address format (default: false) |
| `permission_id` | integer | Permission ID for multi-signature |

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/wallet/cancelallunfreezev2" \
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
  "txID": "jn81i50abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634d002",
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "owner_address": "30e9d79cc47518930bc322d9bf7cddd260a0260a8d"
          },
          "type_url": "type.googleapis.com/protocol.CancelAllUnfreezeV2Contract"
        },
        "type": "CancelAllUnfreezeV2Contract"
      }
    ]
  }
}
```

---

## Error Responses

All endpoints return standardized error responses:

```json
{
  "result": false,
  "code": "ERROR_CODE",
  "message": "Error message describing what went wrong"
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| `SIGERROR` | Signature error |
| `CONTRACT_VALIDATE_ERROR` | Contract validation failed |
| `CONTRACT_EXE_ERROR` | Contract execution failed |
| `BANDWITH_ERROR` | Insufficient bandwidth |
| `DUP_TRANSACTION_ERROR` | Duplicate transaction |
| `TAPOS_ERROR` | TAPOS validation failed |
| `TOO_BIG_TRANSACTION_ERROR` | Transaction too large |
| `TRANSACTION_EXPIRATION_ERROR` | Transaction expired |
| `SERVER_BUSY` | Server is busy |
| `OTHER_ERROR` | Other error |

---

## Best Practices

1. **Always sign transactions** after creating them using a private key (never send private keys to the API)
2. **Check transaction expiration** - transactions expire after 1 minute
3. **Use visible=true** for base58 addresses to make responses human-readable
4. **Implement retry logic** for broadcast failures with exponential backoff
5. **Monitor bandwidth and energy** usage to avoid transaction failures
6. **Cache account information** when appropriate to reduce API calls
7. **Use permission_id** for multi-signature accounts