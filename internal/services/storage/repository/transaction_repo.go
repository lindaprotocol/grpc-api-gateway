package repository

import (
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// SaveTransaction saves a transaction
func (r *TransactionRepository) SaveTransaction(tx *models.Transaction) error {
	return r.db.Save(tx).Error
}

// GetByHash retrieves a transaction by hash
func (r *TransactionRepository) GetByHash(hash string) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.Where("hash = ?", hash).First(&tx).Error
	return &tx, err
}

// GetTransactions retrieves paginated transactions
func (r *TransactionRepository) GetTransactions(blockNumber int64, offset, limit int, sort string) ([]*models.Transaction, int64, error) {
	var txs []*models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{})

	if blockNumber > 0 {
		query = query.Where("block_number = ?", blockNumber)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sort != "" {
		if sort[0] == '-' {
			query = query.Order(sort[1:] + " DESC")
		} else {
			query = query.Order(sort + " ASC")
		}
	} else {
		query = query.Order("block_timestamp DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

// GetTransactionsByAddress retrieves transactions for an address
func (r *TransactionRepository) GetTransactionsByAddress(address string, offset, limit int, sort string) ([]*models.Transaction, int64, error) {
	var txs []*models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).
		Where("from_address = ? OR to_address = ?", address, address)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sort != "" {
		if sort[0] == '-' {
			query = query.Order(sort[1:] + " DESC")
		} else {
			query = query.Order(sort + " ASC")
		}
	} else {
		query = query.Order("block_timestamp DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

// GetTransactionsByContract retrieves transactions for a contract
func (r *TransactionRepository) GetTransactionsByContract(contract string, offset, limit int, sort string) ([]*models.Transaction, int64, error) {
	var txs []*models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).
		Where("contract_address = ?", contract)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sort != "" {
		if sort[0] == '-' {
			query = query.Order(sort[1:] + " DESC")
		} else {
			query = query.Order(sort + " ASC")
		}
	} else {
		query = query.Order("block_timestamp DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

// GetInternalTransactionsByAddress retrieves internal transactions for an address
func (r *TransactionRepository) GetInternalTransactionsByAddress(address string, offset, limit int, sort string) ([]*models.InternalTransaction, int64, error) {
	var txs []*models.InternalTransaction
	var total int64

	// This would need a separate internal_transactions table
	// For now, return empty
	return txs, total, nil
}

// UpdateTransactionWithInfo updates a transaction with info data
func (r *TransactionRepository) UpdateTransactionWithInfo(info *lindapb.TransactionInfo) error {
	return r.db.Model(&models.Transaction{}).
		Where("hash = ?", string(info.Id)).
		Updates(map[string]interface{}{
			"fee":          info.Fee,
			"energy_used":  info.Receipt.EnergyUsage,
			"net_usage":    info.Receipt.NetUsage,
			"result":       info.Result,
		}).Error
}