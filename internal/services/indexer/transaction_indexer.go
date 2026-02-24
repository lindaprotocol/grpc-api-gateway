// internal/services/indexer/transaction_indexer.go
package indexer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

// TransactionIndexer struct: Indexer for transaction operations
type TransactionIndexer struct {
	indexer *Indexer
}

// NewTransactionIndexer creates a new transaction indexer
func NewTransactionIndexer(indexer *Indexer) *TransactionIndexer {
	return &TransactionIndexer{
		indexer: indexer,
	}
}

// IndexTransaction indexes a single transaction
func (ti *TransactionIndexer) IndexTransaction(ctx context.Context, tx *lindapb.Transaction, blockNum int64) error {
	// Convert signature to JSON
	sigJSON, err := json.Marshal(tx.Signature)
	if err != nil {
		return err
	}

	txModel := &models.Transaction{
		Hash:           string(tx.TxID),
		BlockNumber:    blockNum,
		BlockTimestamp: 0, // Will be set from block
		Signature:      models.JSON(sigJSON),
		CreatedAt:      time.Now(),
	}

	// Parse contract data
	if tx.RawData != nil && len(tx.RawData.Contract) > 0 {
		contract := tx.RawData.Contract[0]
		txModel.ContractType = int(contract.Type)

		// Parse parameter based on contract type
		if contract.Parameter != nil {
			var param map[string]interface{}
			if err := json.Unmarshal(contract.Parameter.Value, &param); err == nil {
				if owner, ok := param["owner_address"]; ok {
					if ownerBytes, ok := owner.([]byte); ok {
						txModel.FromAddress = utils.MustHexToBase58(string(ownerBytes))
					} else if ownerStr, ok := owner.(string); ok {
						txModel.FromAddress = utils.MustHexToBase58(ownerStr)
					}
				}
				if to, ok := param["to_address"]; ok {
					if toBytes, ok := to.([]byte); ok {
						txModel.ToAddress = utils.MustHexToBase58(string(toBytes))
					} else if toStr, ok := to.(string); ok {
						txModel.ToAddress = utils.MustHexToBase58(toStr)
					}
				}
				if amount, ok := param["amount"]; ok {
					if amt, ok := amount.(float64); ok {
						txModel.Amount = int64(amt)
					}
				}
				if contractAddr, ok := param["contract_address"]; ok {
					if addr, ok := contractAddr.(string); ok {
						txModel.ContractAddress = utils.MustHexToBase58(addr)
					}
				}
			}
		}
	}

	return ti.indexer.txRepo.SaveTransaction(txModel)
}

// IndexTransactionInfo indexes transaction info data
func (ti *TransactionIndexer) IndexTransactionInfo(info *lindapb.TransactionInfo) error {
	// Update transaction with info
	return ti.indexer.txRepo.UpdateTransactionWithInfo(info)
}

// ExtractInternalTransactions extracts internal transactions from transaction info
func (ti *TransactionIndexer) ExtractInternalTransactions(info *lindapb.TransactionInfo) []*models.InternalTransaction {
	if info == nil || len(info.InternalTransactions) == 0 {
		return nil
	}

	internalTxs := make([]*models.InternalTransaction, 0, len(info.InternalTransactions))
	for _, itx := range info.InternalTransactions {
		// Extract call value info
		callValueInfo := make([]models.CallValueInfo, 0)
		if itx.Data != nil && itx.Data.Extra != nil {
			// Try to extract call value from extra data if available
			// This is a simplified approach - in production you'd need proper parsing
			if val, ok := itx.Data.Extra["callValue"]; ok {
				if callValue, ok := val.(float64); ok {
					callValueInfo = append(callValueInfo, models.CallValueInfo{
						CallValue: int64(callValue),
					})
				}
			}
			if tokenID, ok := itx.Data.Extra["tokenId"]; ok {
				if len(callValueInfo) > 0 {
					callValueInfo[0].TokenID = tokenID
				}
			}
		}

		// Handle extra field for voting information
		var extraJSON models.JSON
		if itx.Data != nil && itx.Data.Extra != nil {
			extraBytes, err := json.Marshal(itx.Data.Extra)
			if err == nil {
				extraJSON = models.JSON(extraBytes)
			}
		}

		// Convert note to hex if it's not already
		noteHex := itx.Data.Note
		if noteHex != "" && !isHex(noteHex) {
			noteHex = hex.EncodeToString([]byte(noteHex))
		}

		internalTxs = append(internalTxs, &models.InternalTransaction{
			Hash:              string(itx.InternalTxId),
			CallerAddress:     utils.MustHexToBase58(string(itx.FromAddress)),
			TransferToAddress: utils.MustHexToBase58(string(itx.ToAddress)),
			CallValueInfo:     callValueInfo,
			Note:              noteHex,
			Rejected:          itx.Data.Rejected,
			Extra:             extraJSON,
			BlockTimestamp:    info.BlockTimeStamp,
			TransactionID:     string(info.Id),
			CreatedAt:         time.Now(),
		})
	}
	return internalTxs
}

// isHex checks if a string is a valid hex string
func isHex(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// IndexTransactionsBatch indexes a batch of transactions
func (ti *TransactionIndexer) IndexTransactionsBatch(ctx context.Context, txs []*lindapb.Transaction, blockNum int64, blockTimestamp int64) error {
	for _, tx := range txs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := ti.IndexTransaction(ctx, tx, blockNum); err != nil {
				ti.indexer.logger.WithError(err).WithField("tx", string(tx.TxID)).Error("Failed to index transaction")
			}
		}
	}
	return nil
}

// GetTransactionByHash retrieves a transaction by hash
func (ti *TransactionIndexer) GetTransactionByHash(hash string) (*models.Transaction, error) {
	return ti.indexer.txRepo.GetByHash(hash)
}

// GetTransactionsByAddress retrieves transactions for an address
func (ti *TransactionIndexer) GetTransactionsByAddress(address string, offset, limit int) ([]*models.Transaction, int64, error) {
	return ti.indexer.txRepo.GetTransactionsByAddress(address, offset, limit, "-timestamp")
}