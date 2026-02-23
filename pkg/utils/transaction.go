// pkg/utils/transaction.go
package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"google.golang.org/protobuf/proto"
)

// GetTransactionRawHex returns the hex-encoded raw data of a transaction
func GetTransactionRawHex(tx *lindapb.Transaction) string {
	if tx == nil || tx.RawData == nil {
		return ""
	}

	// Marshal the RawData to protobuf bytes
	data, err := proto.Marshal(tx.RawData)
	if err != nil {
		// Log the error if needed, but return empty string
		return ""
	}

	return hex.EncodeToString(data)
}

// GetTransactionRawFromHex decodes a hex string back to TransactionRaw
func GetTransactionRawFromHex(hexStr string) (*lindapb.TransactionRaw, error) {
	if hexStr == "" {
		return nil, fmt.Errorf("empty hex string")
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %w", err)
	}

	var raw lindapb.TransactionRaw
	if err := proto.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return &raw, nil
}

// GetTransactionSize returns the approximate size of a transaction
func GetTransactionSize(tx *lindapb.Transaction) int {
	if tx == nil {
		return 0
	}

	size := len(tx.TxID)

	if tx.RawData != nil {
		// Add size of raw data fields
		size += len(tx.RawData.RefBlockBytes)
		size += 8 // RefBlockNum is int64
		size += len(tx.RawData.RefBlockHash)
		size += 8 // Expiration is int64
		size += len(tx.RawData.Auths) * 50 // Approximate size per auth
		size += len(tx.RawData.Data)
		size += len(tx.RawData.Scripts)
		size += 8 // Timestamp is int64
		size += 8 // FeeLimit is int64

		// Add size of contracts
		for _, contract := range tx.RawData.Contract {
			size += 4 // Contract type enum
			size += len(contract.Provider)
			size += len(contract.ContractName)
			size += 8 // PermissionId is int64
			if contract.Parameter != nil {
				size += len(contract.Parameter.Value)
			}
		}
	}

	// Add signature sizes
	for _, sig := range tx.Signature {
		size += len(sig)
	}

	// Add result sizes
	for _, ret := range tx.Ret {
		size += 8 // Fee is int64
		size += 4 // Ret code enum
		// If you need to access ret fields, do it here
		_ = ret // This tells Go you're intentionally not using it
	}

	return size
}

func GetBlockSize(block *lindapb.Block) int {
	if block == nil {
		return 0
	}

	size := len(block.BlockID)
	
	if block.BlockHeader != nil {
		size += len(block.BlockHeader.WitnessSignature)
		if block.BlockHeader.RawData != nil {
			size += len(block.BlockHeader.RawData.TxTrieRoot)
			size += len(block.BlockHeader.RawData.ParentHash)
			size += 8 // Timestamp
			size += 8 // Number
			size += 8 // WitnessId
			size += len(block.BlockHeader.RawData.WitnessAddress)
			size += 4 // Version
			size += len(block.BlockHeader.RawData.AccountStateRoot)
		}
	}

	for _, tx := range block.Transactions {
		size += GetTransactionSize(tx)
	}

	return size
}