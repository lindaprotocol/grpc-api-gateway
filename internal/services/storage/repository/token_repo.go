// internal/services/storage/repository/token_repo.go
package repository

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"gorm.io/gorm"
)

// TokenRepository struct: Repository for token operations
type TokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository function: Creates a new token repository
func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// SaveLRC20Token function: Saves or updates an LRC20 token
func (r *TokenRepository) SaveLRC20Token(token *models.LRC20TokenInfo) error {
	return r.db.Save(token).Error
}

// GetLRC20TokenByContract function: Retrieves an LRC20 token by contract address
func (r *TokenRepository) GetLRC20TokenByContract(contract string) (*models.LRC20TokenInfo, error) {
	var token models.LRC20TokenInfo
	err := r.db.Where("contract = ?", contract).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetLRC20Tokens function: Retrieves paginated LRC20 tokens
func (r *TokenRepository) GetLRC20Tokens(offset, limit int, sort string) ([]*models.LRC20TokenInfo, int64, error) {
	var tokens []*models.LRC20TokenInfo
	var total int64

	query := r.db.Model(&models.LRC20TokenInfo{})

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
		query = query.Order("issue_time DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&tokens).Error; err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

// GetTokenHolders function: Retrieves holders for a token
func (r *TokenRepository) GetTokenHolders(contract string, offset, limit int, sort string) ([]*models.TokenHolderResponse, int64, error) {
	var holders []*models.TokenHolderResponse
	var total int64

	query := r.db.Model(&models.TokenHolder{}).
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
		query = query.Order("balance DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&holders).Error; err != nil {
		return nil, 0, err
	}

	return holders, total, nil
}

// SaveTokenTransfer function: Saves a token transfer
func (r *TokenRepository) SaveTokenTransfer(transfer *models.TokenTransferResponse) error {
	return r.db.Save(transfer).Error
}

// UpdateHolderBalance function: Updates the balance of a token holder
func (r *TokenRepository) UpdateHolderBalance(contractAddr, address string, delta *big.Int) error {
	var holder models.TokenHolder
	err := r.db.Where("contract_address = ? AND address = ?", contractAddr, address).First(&holder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new holder
			holder = models.TokenHolder{
				ContractAddress: contractAddr,
				Address:         address,
				Balance:         delta.String(),
				Percentage:      0,
			}
			return r.db.Create(&holder).Error
		}
		return err
	}

	currentBalance, ok := new(big.Int).SetString(holder.Balance, 10)
	if !ok {
		return errors.New("invalid balance format")
	}
	currentBalance.Add(currentBalance, delta)
	holder.Balance = currentBalance.String()
	return r.db.Save(&holder).Error
}

// GetHolderCount function: Gets the number of token holders
func (r *TokenRepository) GetHolderCount(contractAddr string) (int64, error) {
	var count int64
	err := r.db.Model(&models.TokenHolder{}).Where("contract_address = ?", contractAddr).Count(&count).Error
	return count, err
}

// UpdateHolderCount function: Updates the number of token holders
func (r *TokenRepository) UpdateHolderCount(contractAddr string, count int64) error {
	return r.db.Model(&models.LRC20TokenInfo{}).
		Where("contract = ?", contractAddr).
		Update("holders", count).Error
}

// UpdateHolderPercentage function: Updates the percentage of a token holder
func (r *TokenRepository) UpdateHolderPercentage(contractAddr, address string, percentage float64) error {
	return r.db.Model(&models.TokenHolder{}).
		Where("contract_address = ? AND address = ?", contractAddr, address).
		Update("percentage", percentage).Error
}

// SearchTokens function: Searches for tokens by name or symbol
func (r *TokenRepository) SearchTokens(query string, limit int) ([]models.LRC20TokenInfo, error) {
	var tokens []models.LRC20TokenInfo
	err := r.db.Where("name ILIKE ? OR symbol ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(limit).
		Find(&tokens).Error
	return tokens, err
}

// SaveLRC10Token function: Saves a LRC10 token
func (r *TokenRepository) SaveLRC10Token(token *models.TokenInfo) error {
	return r.db.Save(token).Error
}

// GetTokenTransfers function: Retrieves transfers for a token
func (r *TokenRepository) GetTokenTransfers(contract, from, to string, offset, limit int, sort string) ([]*models.TokenTransferResponse, int64, error) {
	var transfers []*models.TokenTransferResponse
	var total int64

	query := r.db.Model(&models.TokenTransferResponse{})

	if contract != "" {
		query = query.Where("token_address = ?", contract)
	}
	if from != "" {
		query = query.Where("from = ?", from)
	}
	if to != "" {
		query = query.Where("to = ?", to)
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
	if err := query.Offset(offset).Limit(limit).Find(&transfers).Error; err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}

// GetTokensOverview function: Retrieves token overview statistics
func (r *TokenRepository) GetTokensOverview(filter string, offset, limit int, sort string) ([]json.RawMessage, int64, error) {
	// This would aggregate data from multiple tables
	var results []json.RawMessage
	var total int64 = 0
	return results, total, nil
}

// GetTokenPrice function: Retrieves token price information
func (r *TokenRepository) GetTokenPrice() (*models.TokenPriceResponse, error) {
	// This would fetch from cache or external API
	return &models.TokenPriceResponse{
		Price:     0.01,
		Change24h: 5.2,
		Volume24h: 1000000,
		MarketCap: 50000000,
	}, nil
}

// GetParticipations function: Retrieves token participations
func (r *TokenRepository) GetParticipations(offset, limit int) ([]json.RawMessage, int64, error) {
	var participations []json.RawMessage
	var total int64 = 0
	return participations, total, nil
}

// TokenPosition represents a single token position
type TokenPosition struct {
	Address    string  `json:"address"`
	Balance    string  `json:"balance"`
	Percentage float64 `json:"percentage"`
	Rank       int     `json:"rank"`
}

// GetTokenPositionDistribution function: Retrieves token holder distribution
func (r *TokenRepository) GetTokenPositionDistribution(contract string, limit int) ([]TokenPosition, error) {
	var positions []TokenPosition
	err := r.db.Raw(`
		SELECT address, balance, percentage, 
		ROW_NUMBER() OVER (ORDER BY CAST(balance AS NUMERIC) DESC) as rank
		FROM token_holders
		WHERE contract_address = ?
		ORDER BY CAST(balance AS NUMERIC) DESC
		LIMIT ?
	`, contract, limit).Scan(&positions).Error
	if err != nil {
		return nil, err
	}
	return positions, nil
}

// GetWinkFund function: Retrieves WINK fund information
func (r *TokenRepository) GetWinkFund() (*models.WinkFundResponse, error) {
	var fund models.WinkFundResponse
	// Query from database - this is a placeholder
	// In real implementation, you would query from a table
	err := r.db.Raw("SELECT total, burned, circulating FROM wink_fund ORDER BY timestamp DESC LIMIT 1").Scan(&fund).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.WinkFundResponse{
				Total:       0,
				Burned:      0,
				Circulating: 0,
			}, nil
		}
		return nil, err
	}
	return &fund, nil
}

// GetWinkGraphic function: Retrieves WINK graphic data
func (r *TokenRepository) GetWinkGraphic() (*models.GraphicResponse, error) {
	var graphic models.GraphicResponse
	// Query from database - this is a placeholder
	// In real implementation, you would query time series data
	err := r.db.Raw(`
		SELECT array_agg(date) as labels, array_agg(value) as data 
		FROM wink_daily_stats 
		ORDER BY date DESC LIMIT 30
	`).Scan(&graphic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.GraphicResponse{
				Labels: []string{},
				Data:   []float64{},
			}, nil
		}
		return nil, err
	}
	return &graphic, nil
}

// GetJSTFund function: Retrieves JST fund information
func (r *TokenRepository) GetJSTFund() (*models.JSTFundResponse, error) {
	var fund models.JSTFundResponse
	// Query from database - this is a placeholder
	err := r.db.Raw("SELECT total, burned, circulating FROM jst_fund ORDER BY timestamp DESC LIMIT 1").Scan(&fund).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.JSTFundResponse{
				Total:       0,
				Burned:      0,
				Circulating: 0,
			}, nil
		}
		return nil, err
	}
	return &fund, nil
}

// GetJSTGraphic function: Retrieves JST graphic data
func (r *TokenRepository) GetJSTGraphic() (*models.GraphicResponse, error) {
	var graphic models.GraphicResponse
	// Query from database - this is a placeholder
	err := r.db.Raw(`
		SELECT array_agg(date) as labels, array_agg(value) as data 
		FROM jst_daily_stats 
		ORDER BY date DESC LIMIT 30
	`).Scan(&graphic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.GraphicResponse{
				Labels: []string{},
				Data:   []float64{},
			}, nil
		}
		return nil, err
	}
	return &graphic, nil
}

// GetBitTorrentGraphic function: Retrieves BitTorrent graphic data
func (r *TokenRepository) GetBitTorrentGraphic() (*models.GraphicResponse, error) {
	var graphic models.GraphicResponse
	// Query from database - this is a placeholder
	err := r.db.Raw(`
		SELECT array_agg(date) as labels, array_agg(value) as data 
		FROM btt_daily_stats 
		ORDER BY date DESC LIMIT 30
	`).Scan(&graphic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.GraphicResponse{
				Labels: []string{},
				Data:   []float64{},
			}, nil
		}
		return nil, err
	}
	return &graphic, nil
}

// GetAssetTransfers function: Retrieves asset transfers
func (r *TokenRepository) GetAssetTransfers(assetName string, offset, limit int, sort string) ([]*models.TokenTransferResponse, int64, error) {
	var transfers []*models.TokenTransferResponse
	var total int64

	query := r.db.Model(&models.TokenTransferResponse{})
	if assetName != "" {
		query = query.Where("token_symbol = ? OR token_address = ?", assetName, assetName)
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
	if err := query.Offset(offset).Limit(limit).Find(&transfers).Error; err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}