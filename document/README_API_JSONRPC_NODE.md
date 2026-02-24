# Full Node JSON-RPC API

The Full Node JSON-RPC API provides an Ethereum-compatible interface for interacting with the Linda blockchain. This API follows the JSON-RPC 2.0 specification and allows developers familiar with Ethereum to easily integrate with Linda. It supports common Ethereum JSON-RPC methods while also providing Linda-specific extensions for transaction building.

## Base URL

All JSON-RPC API endpoints are served from:

```
https://api.lindagrid.lindacoin.org/jsonrpc
```

## Protocol Specification

The API follows the JSON-RPC 2.0 specification with the following requirements:

| Requirement | Description |
|-------------|-------------|
| **Content-Type** | `application/json` |
| **Request Format** | JSON object with `jsonrpc`, `method`, `params`, and `id` fields |
| **Response Format** | JSON object with `jsonrpc`, `result` or `error`, and `id` fields |
| **Batch Requests** | Supported - send multiple requests in an array |
| **Notification** | Requests without `id` are treated as notifications (no response) |

### Request Format

```json
{
  "jsonrpc": "2.0",
  "method": "method_name",
  "params": [param1, param2, ...],
  "id": 1
}
```

### Success Response

```json
{
  "jsonrpc": "2.0",
  "result": "0x...",
  "id": 1
}
```

### Error Response

```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32000,
    "message": "Error message"
  },
  "id": 1
}
```

## Data Encoding

### Quantity Encoding
- Integers are encoded as hex with a `0x` prefix
- Example: `0x1f4` (500 in decimal)
- Zero is represented as `0x0`
- No leading zeros (except for zero itself)

### Unformatted Data Encoding
- Byte arrays, hashes, addresses are hex encoded with `0x` prefix
- Must have an even number of hex digits
- Example: `0x30f0cc5a2a84cd0f68ed1667070934542d673acbd8`

---

## Ethereum-Compatible Methods

### eth_blockNumber

Returns the number of the most recent block.

**Method:** `eth_blockNumber`

**Parameters:** None

**Returns:** `QUANTITY` - The latest block number

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_blockNumber",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x20e0cf0"
}
```

---

### eth_getBalance

Returns the balance of an address.

**Method:** `eth_getBalance`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA`, 20 Bytes | Address to check |
| 2 | `QUANTITY|TAG` | Block number or "latest" |

**Returns:** `QUANTITY` - Current balance in SUN

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0x30f0cc5a2a84cd0f68ed1667070934542d673acbd8", "latest"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x492780"
}
```

---

### eth_getTransactionByHash

Returns information about a transaction by hash.

**Method:** `eth_getTransactionByHash`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA`, 32 Bytes | Transaction hash |

**Returns:** `Object` - Transaction object, or `null` if not found

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionByHash",
    "params": ["0xc9af231ad59bcd7e8dcf827afd45020a02112704dce74ec5f72cb090aa07eef0"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "blockHash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
    "blockNumber": "0x20ef11c",
    "from": "0x6eced5214d62c3bc9eaa742e2f86d5c516785e14",
    "gas": "0x0",
    "gasPrice": "0x8c",
    "hash": "0xc9af231ad59bcd7e8dcf827afd45020a02112704dce74ec5f72cb090aa07eef0",
    "input": "0x",
    "nonce": null,
    "r": "0x433eaf0a7df3a08c8828a2180987146d39d44de4ac327c4447d0eeda42230ea8",
    "s": "0x6f91f63b37f4d1cd9342f570205beefaa5b5ba18d616fec643107f8c1ae1339d",
    "to": "0x0697250b9d73b460a9d2bbfd8c4cacebb05dd1f1",
    "transactionIndex": "0x6",
    "type": "0x0",
    "v": "0x1b",
    "value": "0x1cb2310"
  }
}
```

---

### eth_getTransactionReceipt

Returns the receipt of a transaction.

**Method:** `eth_getTransactionReceipt`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA`, 32 Bytes | Transaction hash |

**Returns:** `Object` - Transaction receipt, or `null` if not found

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionReceipt",
    "params": ["0xc9af231ad59bcd7e8dcf827afd45020a02112704dce74ec5f72cb090aa07eef0"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "blockHash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
    "blockNumber": "0x20ef11c",
    "contractAddress": null,
    "cumulativeGasUsed": "0x646e2",
    "effectiveGasPrice": "0x8c",
    "from": "0x6eced5214d62c3bc9eaa742e2f86d5c516785e14",
    "gasUsed": "0x0",
    "logs": [],
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "status": "0x1",
    "to": "0x0697250b9d73b460a9d2bbfd8c4cacebb05dd1f1",
    "transactionHash": "0xc9af231ad59bcd7e8dcf827afd45020a02112704dce74ec5f72cb090aa07eef0",
    "transactionIndex": "0x6",
    "type": "0x0"
  }
}
```

---

### eth_getBlockByNumber

Returns information about a block by block number.

**Method:** `eth_getBlockByNumber`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `QUANTITY|TAG` | Block number or "latest"/"earliest" |
| 2 | `Boolean` | If true, returns full transaction objects; if false, returns transaction hashes |

**Returns:** `Object` - Block object, or `null` if not found

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockByNumber",
    "params": ["0xF9CC56", true],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "number": "0xf9cc56",
    "hash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
    "parentHash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc59",
    "nonce": "0x0000000000000000",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "logsBloom": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "transactionsRoot": "0xb1144687a9f5e9a1f1c3d8a9b2f4e5d6c7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2",
    "stateRoot": "0xc2d4e6f8a0b2c4d6e8f0a2b4c6d8e0f2a4b6c8d0e2f4a6b8c0d2e4f6a8b0c2d4",
    "receiptsRoot": "0xd3e5f7a9b1c3d5e7f9a1b3c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5a7b9c1d3",
    "miner": "0x30f0cc5a2a84cd0f68ed1667070934542d673acbd8",
    "difficulty": "0x0",
    "totalDifficulty": "0x0",
    "extraData": "0x",
    "size": "0x3e8",
    "gasLimit": "0x0",
    "gasUsed": "0x0",
    "timestamp": "0x5e8a4b3c",
    "transactions": [],
    "uncles": []
  }
}
```

---

### eth_getBlockByHash

Returns information about a block by block hash.

**Method:** `eth_getBlockByHash`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA`, 32 Bytes | Block hash |
| 2 | `Boolean` | If true, returns full transaction objects; if false, returns transaction hashes |

**Returns:** `Object` - Block object, or `null` if not found

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockByHash",
    "params": ["0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a", false],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "number": "0xf9cc56",
    "hash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5a",
    "parentHash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc59",
    "transactions": ["0x7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc"]
  }
}
```

---

### eth_getBlockTransactionCountByNumber

Returns the number of transactions in a block by block number.

**Method:** `eth_getBlockTransactionCountByNumber`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `QUANTITY|TAG` | Block number or "latest" |

**Returns:** `QUANTITY` - Number of transactions in the block

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockTransactionCountByNumber",
    "params": ["0xF96B0F"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x23"
}
```

---

### eth_getCode

Returns the runtime code of a smart contract.

**Method:** `eth_getCode`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA`, 20 Bytes | Contract address |
| 2 | `QUANTITY|TAG` | Block number or "latest" |

**Returns:** `DATA` - Runtime bytecode of the contract

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getCode",
    "params": ["0x3070082243784DCDF3042034E7B044D6D342A91360", "latest"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x6080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde03811461005b575b600080fd5b34801561006757600080fd5b50610070610086565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100ae578082015181840152602081019050610093565b50505050905090810190601f1680156100db5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040805190810160405280600781526020017f4d79546f6b656e0000000000000000000000000000000000000000000000000081525090509056fea165627a7a72305820a2fb39541e90eda9a2f5f9e7905ef98e66e60dd4b38e00b05de418da3154e7570029"
}
```

---

### eth_getStorageAt

Returns the value from a storage position at a given address.

**Method:** `eth_getStorageAt`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA`, 20 Bytes | Contract address |
| 2 | `QUANTITY` | Storage position (integer) |
| 3 | `QUANTITY|TAG` | Block number or "latest" |

**Returns:** `DATA` - Value at the storage position

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getStorageAt",
    "params": ["0xE94EAD5F4CA072A25B2E5500934709F1AEE3C64B", "0x0", "latest"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x0000000000000000000000000000000000000000000000000000000000000000"
}
```

---

### eth_call

Executes a message call immediately without creating a transaction.

**Method:** `eth_call`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `Object` | Call object |
| 2 | `QUANTITY|TAG` | Block number or "latest" |

**Call Object Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `from` | `DATA`, 20 Bytes | Caller address |
| `to` | `DATA`, 20 Bytes | Contract address |
| `gas` | `QUANTITY` | Not used (set to 0x0) |
| `gasPrice` | `QUANTITY` | Not used (set to 0x0) |
| `value` | `QUANTITY` | Not used (set to 0x0) |
| `data` | `DATA` | Function selector + ABI-encoded parameters |

**Returns:** `DATA` - Result of the contract call

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_call",
    "params": [{
      "from": "0xF0CC5A2A84CD0F68ED1667070934542D673ACBD8",
      "to": "0x70082243784DCDF3042034E7B044D6D342A91360",
      "gas": "0x0",
      "gasPrice": "0x0",
      "value": "0x0",
      "data": "0x70a08231000000000000000000000041f0cc5a2a84cd0f68ed1667070934542d673acbd8"
    }, "latest"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x000000000000000000000000000000000000000000000000000000003b9aca00"
}
```

---

### eth_estimateGas

Estimates the gas (energy) required for a transaction.

**Method:** `eth_estimateGas`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `Object` | Call object (same as eth_call) |

**Returns:** `QUANTITY` - Estimated gas required

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_estimateGas",
    "params": [{
      "from": "0x30F0CC5A2A84CD0F68ED1667070934542D673ACBD8",
      "to": "0x3070082243784DCDF3042034E7B044D6D342A91360",
      "value": "0x1",
      "data": "0x70a08231000000000000000000000041f0cc5a2a84cd0f68ed1667070934542d673acbd8"
    }],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0xc350"
}
```

---

### eth_gasPrice

Returns the current energy price.

**Method:** `eth_gasPrice`

**Parameters:** None

**Returns:** `QUANTITY` - Current energy price in SUN

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_gasPrice",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x8c"
}
```

---

### eth_chainId

Returns the chain ID of the Linda network.

**Method:** `eth_chainId`

**Parameters:** None

**Returns:** `DATA` - Chain ID (last 4 bytes of genesis block hash)

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_chainId",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x2b6653dc"
}
```

---

### eth_getLogs

Returns logs matching the filter criteria.

**Method:** `eth_getLogs`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `Object` | Filter object |

**Filter Object Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `fromBlock` | `QUANTITY|TAG` | Start block (optional) |
| `toBlock` | `QUANTITY|TAG` | End block (optional) |
| `address` | `DATA|Array` | Contract address(es) to filter (optional) |
| `topics` | `Array` | Topic filters (optional) |
| `blockHash` | `DATA`, 32 Bytes | Block hash (optional - exclusive with from/toBlock) |

**Returns:** `Array` - Array of log objects

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getLogs",
    "params": [{
      "address": "0xE518C608A37E2A262050E10BE0C9D03C7A0877F3",
      "fromBlock": "0x989680",
      "toBlock": "0x9959d0",
      "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]
    }],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": [
    {
      "address": "0xe518c608a37e2a262050e10be0c9d03c7a0877f3",
      "topics": [
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
        "0x000000000000000000000041f0cc5a2a84cd0f68ed1667070934542d673acbd8",
        "0x00000000000000000000004195fd23d3d2221cfef64167938de5e62074719e54"
      ],
      "data": "0x000000000000000000000000000000000000000000000000000000003b9aca00",
      "blockNumber": "0x9926a8",
      "transactionHash": "0x7c2d4206c03a883dd9066d620335dc1be272a8dc733cfa3f6d10308faa37facc",
      "transactionIndex": "0x0",
      "blockHash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
      "logIndex": "0x0",
      "removed": false
    }
  ]
}
```

---

### eth_getTransactionByBlockHashAndIndex

Returns a transaction by block hash and transaction index.

**Method:** `eth_getTransactionByBlockHashAndIndex`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA`, 32 Bytes | Block hash |
| 2 | `QUANTITY` | Transaction index position |

**Returns:** `Object` - Transaction object, or `null` if not found

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionByBlockHashAndIndex",
    "params": ["0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b", "0x0"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "blockHash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
    "blockNumber": "0x20ef11c",
    "from": "0xb4f1b6e3a1461266b01c2c4ff9237191d5c3d5ce",
    "gas": "0x0",
    "gasPrice": "0x8c",
    "hash": "0x8dd26d1772231569f022adb42f7d7161dee88b97b4b35eeef6ce73fcd6613bc2",
    "input": "0x",
    "nonce": null,
    "r": "0x6212a53b962345fb8ab02215879a2de05f32e822c54e257498f0b70d33825cc5",
    "s": "0x6e04221f5311cf2b70d3aacfc444e43a5cf14d0bf31d9227218efaabd9b5a812",
    "to": "0x047d4a0a1b7a9d495d6503536e2a49bb5cc72cfe",
    "transactionIndex": "0x0",
    "type": "0x0",
    "v": "0x1b",
    "value": "0x203226"
  }
}
```

---

### eth_getTransactionByBlockNumberAndIndex

Returns a transaction by block number and transaction index.

**Method:** `eth_getTransactionByBlockNumberAndIndex`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `QUANTITY|TAG` | Block number or tag |
| 2 | `QUANTITY` | Transaction index position |

**Returns:** `Object` - Transaction object, or `null` if not found

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionByBlockNumberAndIndex",
    "params": ["0xfb82f0", "0x0"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": null
}
```

---

### eth_getBlockReceipts

Returns all transaction receipts in a block.

**Method:** `eth_getBlockReceipts`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA|TAG` | Block identifier (hash, number, or tag) |

**Returns:** `Array` - Array of transaction receipt objects

**Example Request (by block number):**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockReceipts",
    "params": ["0x377a8a2"],
    "id": 1
  }'
```

**Example Request (by block hash):**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockReceipts",
    "params": ["0x00000000049e470616a96ca7af19fc46a473e9733796960d840697dd70ac14ad"],
    "id": 1
  }'
```

**Example Request (by tag):**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockReceipts",
    "params": ["latest"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": [{
    "blockHash": "0x00000000020ef11c87517739090601aa0a7be1de6faebf35ddb14e7ab7d1cc5b",
    "blockNumber": "0x20ef11c",
    "contractAddress": null,
    "cumulativeGasUsed": "0x646e2",
    "effectiveGasPrice": "0x8c",
    "from": "0x6eced5214d62c3bc9eaa742e2f86d5c516785e14",
    "gasUsed": "0x0",
    "logs": [],
    "logsBloom": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "status": "0x1",
    "to": "0x0697250b9d73b460a9d2bbfd8c4cacebb05dd1f1",
    "transactionHash": "0xc9af231ad59bcd7e8dcf827afd45020a02112704dce74ec5f72cb090aa07eef0",
    "transactionIndex": "0x6",
    "type": "0x0"
  }]
}
```

---

### eth_accounts

Returns a list of addresses owned by the client.

**Method:** `eth_accounts`

**Parameters:** None

**Returns:** `Array` - Empty array (in LINDA, this always returns empty)

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_accounts",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": []
}
```

---

### eth_coinbase

Returns the Super Representative address of the current node.

**Method:** `eth_coinbase`

**Parameters:** None

**Returns:** `DATA` - SR address, or error if not configured

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_coinbase",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "etherbase must be explicitly specified",
    "data": "{}"
  }
}
```

---

### eth_protocolVersion

Returns the current Linda protocol version.

**Method:** `eth_protocolVersion`

**Parameters:** None

**Returns:** `String` - Protocol version

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_protocolVersion",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x16"
}
```

---

### eth_syncing

Returns information about the sync status of the node.

**Method:** `eth_syncing`

**Parameters:** None

**Returns:** `Object|Boolean` - Sync status object or false if not syncing

**Example Request (when syncing):**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_syncing",
    "params": [],
    "id": 1
  }'
```

**Example Response (when syncing):**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "startingBlock": "0x20e76cc",
    "currentBlock": "0x20e76df",
    "highestBlock": "0x20e76e0"
  }
}
```

**Example Response (when not syncing):**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": false
}
```

---

### eth_getWork

Returns the hash of the current block.

**Method:** `eth_getWork`

**Parameters:** None

**Returns:** `Array` - Array with block hash

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getWork",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": ["0x00000000020e73915413df0c816e327dc4b9d17069887aef1fff0e854f8d9ad0", null, null]
}
```

---

## Filter Methods

### eth_newFilter

Creates a filter object to monitor for logs.

**Method:** `eth_newFilter`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `Object` | Filter options (same as eth_getLogs) |

**Returns:** `QUANTITY` - Filter ID

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_newFilter",
    "params": [{
      "address": ["0xcc2e32f2388f0096fae9b055acffd76d4b3e5532", "0xE518C608A37E2A262050E10BE0C9D03C7A0877F3"],
      "fromBlock": "0x989680",
      "toBlock": "0x9959d0",
      "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]
    }],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x2bab51aee6345d2748e0a4a3f4569d80"
}
```

---

### eth_newBlockFilter

Creates a filter to monitor for new blocks.

**Method:** `eth_newBlockFilter`

**Parameters:** None

**Returns:** `QUANTITY` - Filter ID

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_newBlockFilter",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0xc11a84d5e906ecb9f5c1eb65ee940b154ad37dce8f5ac29c80764508b901d996"
}
```

---

### eth_getFilterChanges

Polls a filter for changes since the last poll.

**Method:** `eth_getFilterChanges`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `QUANTITY` | Filter ID |

**Returns:** `Array` - Array of log objects or block hashes

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getFilterChanges",
    "params": ["0xc11a84d5e906ecb9f5c1eb65ee940b154ad37dce8f5ac29c80764508b901d996"],
    "id": 1
  }'
```

**Example Response (if filter not found):**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "filter not found",
    "data": "{}"
  }
}
```

---

### eth_getFilterLogs

Returns all logs matching a filter.

**Method:** `eth_getFilterLogs`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `QUANTITY` | Filter ID |

**Returns:** `Array` - Array of log objects

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getFilterLogs",
    "params": ["0xc11a84d5e906ecb9f5c1eb65ee940b154ad37dce8f5ac29c80764508b901d996"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "filter not found",
    "data": "{}"
  }
}
```

---

### eth_uninstallFilter

Uninstalls a filter.

**Method:** `eth_uninstallFilter`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `QUANTITY` | Filter ID |

**Returns:** `Boolean` - True if successfully uninstalled

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_uninstallFilter",
    "params": ["0xc11a84d5e906ecb9f5c1eb65ee940b154ad37dce8f5ac29c80764508b901d996"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": true
}
```

---

## Network Methods

### net_version

Returns the chain ID (genesis block hash).

**Method:** `net_version`

**Parameters:** None

**Returns:** `String` - Chain ID

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "net_version",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x2b6653dc"
}
```

---

### net_listening

Returns true if the client is listening for network connections.

**Method:** `net_listening`

**Parameters:** None

**Returns:** `Boolean` - Listening status

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "net_listening",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": true
}
```

---

### net_peerCount

Returns the number of connected peers.

**Method:** `net_peerCount`

**Parameters:** None

**Returns:** `QUANTITY` - Number of connected peers

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "net_peerCount",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x9"
}
```

---

## Web3 Methods

### web3_clientVersion

Returns the current client version.

**Method:** `web3_clientVersion`

**Parameters:** None

**Returns:** `String` - Client version

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "web3_clientVersion",
    "params": [],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "LINDA/v4.3.0/Linux/Java1.8/GreatVoyage-v4.2.2.1-281-gc1d9dfd6c"
}
```

---

### web3_sha3

Returns Keccak-256 hash of the given data.

**Method:** `web3_sha3`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `DATA` | Data to hash |

**Returns:** `DATA` - Keccak-256 hash

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "web3_sha3",
    "params": ["0x68656c6c6f20776f726c64"],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
}
```

---

## Linda-Specific Transaction Building Methods

### buildTransaction (TransferContract)

Creates a LIND transfer transaction.

**Method:** `buildTransaction`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `Object` | Transaction parameters |

**Transaction Parameters (TransferContract):**

| Field | Type | Description |
|-------|------|-------------|
| `from` | `DATA`, 20 Bytes | Sender address |
| `to` | `DATA`, 20 Bytes | Recipient address |
| `value` | `DATA` | Amount to transfer (hex) |

**Returns:** `Object` - Unsigned transaction object

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "buildTransaction",
    "params": [{
      "from": "0xC4DB2C9DFBCB6AA344793F1DDA7BD656598A06D8",
      "to": "0x95FD23D3D2221CFEF64167938DE5E62074719E54",
      "value": "0x1f4"
    }],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "transaction": {
      "visible": false,
      "txID": "ae02a80abd985a6f05478b9bbf04706f00cdbf71e38c77d21ed77e44c634cef9",
      "raw_data": {
        "contract": [{
          "parameter": {
            "value": {
              "amount": 500,
              "owner_address": "30c4db2c9dfbcb6aa344793f1dda7bd656598a06d8",
              "to_address": "3095fd23d3d2221cfef64167938de5e62074719e54"
            },
            "type_url": "type.googleapis.com/protocol.TransferContract"
          },
          "type": "TransferContract"
        }],
        "ref_block_bytes": "957e",
        "ref_block_hash": "3922d8c0d28b5283",
        "expiration": 1684469286000,
        "timestamp": 1684469226841
      },
      "raw_data_hex": "0a02957e22083922d8c0d28b528340f088c69183315a66080112620a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412310a1541c4db2c9dfbcb6aa344793f1dda7bd656598a06d812154195fd23d3d2221cfef64167938de5e62074719e5418f40370d9bac2918331"
    }
  }
}
```

---

### buildTransaction (TransferAssetContract)

Creates an LRC-10 token transfer transaction.

**Method:** `buildTransaction`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `Object` | Transaction parameters |

**Transaction Parameters (TransferAssetContract):**

| Field | Type | Description |
|-------|------|-------------|
| `from` | `DATA`, 20 Bytes | Sender address |
| `to` | `DATA`, 20 Bytes | Recipient address |
| `tokenId` | `QUANTITY` | Token ID |
| `tokenValue` | `QUANTITY` | Amount to transfer |

**Returns:** `Object` - Unsigned transaction object

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "buildTransaction",
    "params": [{
      "from": "0xC4DB2C9DFBCB6AA344793F1DDA7BD656598A06D8",
      "to": "0x95FD23D3D2221CFEF64167938DE5E62074719E54",
      "tokenId": 1000016,
      "tokenValue": 20
    }],
    "id": 1
  }'
```

**Example Response (if error):**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32600,
    "message": "assetBalance must be greater than 0.",
    "data": "{}"
  }
}
```

---

### buildTransaction (CreateSmartContract)

Creates a smart contract deployment transaction.

**Method:** `buildTransaction`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `Object` | Transaction parameters |

**Transaction Parameters (CreateSmartContract):**

| Field | Type | Description |
|-------|------|-------------|
| `from` | `DATA`, 20 Bytes | Deployer address |
| `name` | `DATA` | Contract name |
| `gas` | `DATA` | Fee limit |
| `abi` | `DATA` | Contract ABI JSON string |
| `data` | `DATA` | Contract bytecode |
| `consumeUserResourcePercent` | `QUANTITY` | User resource consumption percentage |
| `originEnergyLimit` | `QUANTITY` | Origin energy limit |
| `value` | `DATA` | Call value |
| `tokenId` | `QUANTITY` | Token ID |
| `tokenValue` | `QUANTITY` | Token value |

**Returns:** `Object` - Unsigned transaction object

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "buildTransaction",
    "params": [{
      "from": "0xC4DB2C9DFBCB6AA344793F1DDA7BD656598A06D8",
      "name": "transferTokenContract",
      "gas": "0x245498",
      "abi": "[{\"constant\":false,\"inputs\":[],\"name\":\"getResultInCon\",\"outputs\":[{\"name\":\"\",\"type\":\"lrcToken\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"}]",
      "data": "6080604052d3600055d2600155346002556101418061001f6000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166305c24200811461005b5780633be9ece71461008157806371dc08ce146100aa575b600080fd5b6100636100b2565b60408051938452602084019290925282820152519081900360600190f35b6100a873ffffffffffffffffffffffffffffffffffffffff600435166024356044356100c0565b005b61006361010d565b600054600154600254909192565b60405173ffffffffffffffffffffffffffffffffffffffff84169082156108fc029083908590600081818185878a8ad0945050505050158015610107573d6000803e3d6000fd5b50505050565bd3d2349091925600a165627a7a72305820a2fb39541e90eda9a2f5f9e7905ef98e66e60dd4b38e00b05de418da3154e7570029",
      "consumeUserResourcePercent": 100,
      "originEnergyLimit": 11111111111111,
      "value": "0x1f4",
      "tokenId": 1000033,
      "tokenValue": 100000
    }],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "transaction": {
      "visible": false,
      "txID": "598d8aafbf9340e92c8f72a38389ce9661b643ff37dd2a609f393336a76025b9",
      "contract_address": "30dfd93697c0a978db343fe7a92333e11eeb2f967d",
      "raw_data": {
        "contract": [{
          "parameter": {
            "value": {
              "token_id": 1000033,
              "owner_address": "30c4db2c9dfbcb6aa344793f1dda7bd656598a06d8",
              "call_token_value": 100000,
              "new_contract": {
                "bytecode": "6080604052d3600055d2600155346002556101418061001f6000396000f300...",
                "consume_user_resource_percent": 100,
                "name": "transferTokenContract",
                "origin_address": "30c4db2c9dfbcb6aa344793f1dda7bd656598a06d8",
                "abi": {"entrys": [...]},
                "origin_energy_limit": 11111111111111,
                "call_value": 500
              }
            },
            "type_url": "type.googleapis.com/protocol.CreateSmartContract"
          },
          "type": "CreateSmartContract"
        }],
        "ref_block_bytes": "80be",
        "ref_block_hash": "ac7c3d59c55ac92c",
        "expiration": 1634030190000,
        "fee_limit": 333333280,
        "timestamp": 1634030131693
      },
      "raw_data_hex": "0a0280be2208ac7c3d59c55ac92c40b0fba79ec72f5ad805081e12d3050a30747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e437265617465536d617274436f6e7472616374129e050a1541c4db2c9dfbcb6aa344793f1dda7bd656598a06d812fc040a1541c4db2c9dfbcb6aa344793f1dda7bd656598a06d81adb010a381a0e676574526573756c74496e436f6e2a0a1a08747263546f6b656e2a091a0775696e743235362a091a0775696e743235363002380140040a501a0f5472616e73666572546f6b656e546f22141209746f416464726573731a0761646472657373220e120269641a08747263546f6b656e22111206616d6f756e741a0775696e743235363002380140040a451a1b6d7367546f6b656e56616c7565416e64546f6b656e4964546573742a0a1a08747263546f6b656e2a091a0775696e743235362a091a0775696e743235363002380140040a0630013801400422e0026080604052d3600055d2600155346002556101418061001f6000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166305c24200811461005b5780633be9ece71461008157806371dc08ce146100aa575b600080fd5b6100636100b2565b60408051938452602084019290925282820152519081900360600190f35b6100a873ffffffffffffffffffffffffffffffffffffffff600435166024356044356100c0565b005b61006361010d565b600054600154600254909192565b60405173ffffffffffffffffffffffffffffffffffffffff84169082156108fc029083908590600081818185878a8ad0945050505050158015610107573d6000803e3d6000fd5b50505050565bd3d2349091925600a165627a7a72305820a2fb39541e90eda9a2f5f9e7905ef98e66e60dd4b38e00b05de418da3154e757002928f40330643a157472616e73666572546f6b656e436f6e747261637440c7e3d28eb0c30218a08d0620e1843d70edb3a49ec72f9001a086f99e01"
    }
  }
}
```

---

### buildTransaction (TriggerSmartContract)

Creates a smart contract call transaction.

**Method:** `buildTransaction`

**Parameters:**

| Position | Type | Description |
|----------|------|-------------|
| 1 | `Object` | Transaction parameters |

**Transaction Parameters (TriggerSmartContract):**

| Field | Type | Description |
|-------|------|-------------|
| `from` | `DATA`, 20 Bytes | Caller address |
| `to` | `DATA`, 20 Bytes | Contract address |
| `data` | `DATA` | Function selector + encoded parameters |
| `gas` | `DATA` | Fee limit |
| `value` | `DATA` | Call value |
| `tokenId` | `QUANTITY` | Token ID |
| `tokenValue` | `QUANTITY` | Token value |

**Returns:** `Object` - Unsigned transaction object

**Example Request:**

```bash
curl -X POST "https://api.lindagrid.lindacoin.org/jsonrpc" \
  -H "LINDA-PRO-API-KEY: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "buildTransaction",
    "params": [{
      "from": "0xC4DB2C9DFBCB6AA344793F1DDA7BD656598A06D8",
      "to": "0xf859b5c93f789f4bcffbe7cc95a71e28e5e6a5bd",
      "data": "0x3be9ece7000000000000000000000000ba8e28bdb6e49fbb3f5cd82a9f5ce8363587f1f600000000000000000000000000000000000000000000000000000000000f42630000000000000000000000000000000000000000000000000000000000000001",
      "gas": "0x245498",
      "value": "0xA",
      "tokenId": 1000035,
      "tokenValue": 20
    }],
    "id": 1
  }'
```

**Example Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "transaction": {
      "visible": false,
      "txID": "c3c746beb86ffc366ec0ff8bf6c9504c88f8714e47bc0009e4f7e2b1d49eb967",
      "raw_data": {
        "contract": [{
          "parameter": {
            "value": {
              "amount": 10,
              "owner_address": "30c4db2c9dfbcb6aa344793f1dda7bd656598a06d8",
              "to_address": "30f859b5c93f789f4bcffbe7cc95a71e28e5e6a5bd"
            },
            "type_url": "type.googleapis.com/protocol.TransferContract"
          },
          "type": "TransferContract"
        }],
        "ref_block_bytes": "958c",
        "ref_block_hash": "9d8c6bae734a2281",
        "expiration": 1684469328000,
        "timestamp": 1684469270364
      },
      "raw_data_hex": "0a02958c22089d8c6bae734a22814080d1c89183315a65080112610a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412300a1541c4db2c9dfbcb6aa344793f1dda7bd656598a06d8121541f859b5c93f789f4bcffbe7cc95a71e28e5e6a5bd180a70dc8ec5918331"
    }
  }
}
```

---

## Error Codes

| Code | Description |
|------|-------------|
| `-32700` | Parse error - Invalid JSON |
| `-32600` | Invalid request |
| `-32601` | Method not found |
| `-32602` | Invalid parameters |
| `-32603` | Internal error |
| `-32000` | Server error - Implementation-specific error |

---

## Best Practices

1. **Use JSON-RPC 2.0** - Always include `jsonrpc: "2.0"` in requests
2. **Include unique IDs** - Use incrementing integers for request IDs
3. **Handle errors gracefully** - Check for `error` field in responses
4. **Batch requests** - Combine multiple requests for efficiency
5. **Use hex encoding** - All numeric values must be hex-encoded with `0x` prefix
6. **Address format** - Use hex format with `30` prefix for Linda addresses
7. **Gas estimation** - Use `eth_estimateGas` before sending transactions
8. **Poll responsibly** - Use appropriate intervals for filter polling