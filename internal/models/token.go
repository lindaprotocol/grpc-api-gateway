// internal/models/token.go
package models

import (
	"time"
)

// TokenInfo represents LRC-10 token database model
type TokenInfo struct {
	ID          string    `gorm:"primaryKey;type:varchar(100)" json:"id"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Symbol      string    `gorm:"type:varchar(20)" json:"symbol"`
	TotalSupply int64     `json:"total_supply"`
	Owner       string    `gorm:"type:varchar(42)" json:"owner"`
	Decimals    int       `json:"decimals"`
	StartTime   int64     `json:"start_time"`
	EndTime     int64     `json:"end_time"`
	URL         string    `gorm:"type:text" json:"url"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// LRC20TokenInfo represents LRC20 token database model
type LRC20TokenInfo struct {
	ID          uint      `gorm:"primarykey" json:"-"`
	Contract    string    `gorm:"uniqueIndex;type:varchar(42)" json:"contract"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Symbol      string    `gorm:"type:varchar(20)" json:"symbol"`
	Decimals    int32     `json:"decimals"`
	TotalSupply string    `gorm:"type:varchar(100)" json:"total_supply"`
	Owner       string    `gorm:"type:varchar(42)" json:"owner"`
	IssueTime   int64     `json:"issue_time"`
	Holders     int64     `json:"holders"`
	Transfers   int64     `json:"transfers"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TokenHolder represents a token holder database model
type TokenHolder struct {
	ID              uint      `gorm:"primarykey" json:"-"`
	ContractAddress string    `gorm:"index;type:varchar(42)" json:"contract_address"`
	Address         string    `gorm:"index;type:varchar(42)" json:"address"`
	Balance         string    `gorm:"type:varchar(100)" json:"balance"`
	Percentage      float64   `json:"percentage"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TokenTransfer represents a token transfer database model
type TokenTransfer struct {
	ID             uint      `gorm:"primarykey" json:"-"`
	TransactionID  string    `gorm:"index;type:varchar(64)" json:"transaction_id"`
	BlockNumber    int64     `gorm:"index" json:"block_number"`
	BlockTimestamp int64     `gorm:"index" json:"block_timestamp"`
	From           string    `gorm:"index;type:varchar(42)" json:"from"`
	To             string    `gorm:"index;type:varchar(42)" json:"to"`
	Value          string    `gorm:"type:varchar(100)" json:"value"`
	TokenAddress   string    `gorm:"index;type:varchar(42)" json:"token_address"`
	TokenSymbol    string    `gorm:"type:varchar(20)" json:"token_symbol"`
	TokenDecimals  int32     `json:"token_decimals"`
	CreatedAt      time.Time `json:"created_at"`
}