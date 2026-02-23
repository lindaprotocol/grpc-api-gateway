package repository

import (
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// SaveAccount saves or updates an account
func (r *AccountRepository) SaveAccount(account *models.AccountResponse) error {
	return r.db.Save(account).Error
}

// GetByAddress retrieves an account by address
func (r *AccountRepository) GetByAddress(address string) (*models.AccountResponse, error) {
	var account models.AccountResponse
	err := r.db.Where("address = ?", address).First(&account).Error
	return &account, err
}

// GetList retrieves paginated list of accounts
func (r *AccountRepository) GetList(offset, limit int, sort string) ([]*models.AccountResponse, int64, error) {
	var accounts []*models.AccountResponse
	var total int64

	query := r.db.Model(&models.AccountResponse{})

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
		query = query.Order("balance DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// UpdateBalance updates account balance
func (r *AccountRepository) UpdateBalance(address string, balance int64) error {
	return r.db.Model(&models.AccountResponse{}).
		Where("address = ?", address).
		Update("balance", balance).Error
}

// GetTopAccounts retrieves top accounts by balance
func (r *AccountRepository) GetTopAccounts(limit int) ([]*models.AccountResponse, error) {
	var accounts []*models.AccountResponse
	err := r.db.Order("balance DESC").Limit(limit).Find(&accounts).Error
	return accounts, err
}