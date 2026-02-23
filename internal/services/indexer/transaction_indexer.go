// internal/services/indexer/transaction_indexer.go
package indexer

import (
	"context"
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

	internalTxs := make([]*models.InternalTransaction, len(info.InternalTransactions))
	for i, itx := range info.InternalTransactions {
		internalTxs[i] = &models.InternalTransaction{
			InternalTxID:    string(itx.InternalTxId),
			Data:            convertInternalTransactionData(itx.Data),
			BlockTimestamp:  info.BlockTimeStamp,
			ToAddress:       utils.MustHexToBase58(string(itx.ToAddress)),
			TxID:            string(info.Id),
			FromAddress:     utils.MustHexToBase58(string(itx.FromAddress)),
		}
	}
	return internalTxs
}

// convertInternalTransactionData converts protobuf internal transaction data to model
func convertInternalTransactionData(data *lindapb.InternalTransactionData) *models.InternalTransactionData {
	if data == nil {
		return nil
	}
	return &models.InternalTransactionData{
		Note:     data.Note,
		Rejected: data.Rejected,
	}
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