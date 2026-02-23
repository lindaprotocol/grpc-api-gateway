package models

import (
	"time"
)

// TokenInfo represents basic token information
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

// LRC20TokenInfo represents LRC20 token information
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
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// TokenHolder represents a token holder
type TokenHolder struct {
	ID              uint      `gorm:"primarykey" json:"-"`
	ContractAddress string    `gorm:"index;type:varchar(42)" json:"-"`
	Address         string    `gorm:"index;type:varchar(42)" json:"address"`
	Balance         string    `gorm:"type:varchar(100)" json:"balance"`
	Percentage      float64   `json:"percentage"`
	UpdatedAt       time.Time `json:"-"`
}

// TokenHolderResponse is the API response for token holders
type TokenHolderResponse struct {
	Address    string  `json:"address"`
	Balance    string  `json:"balance"`
	Percentage float64 `json:"percentage"`
	Rank       int64   `json:"rank,omitempty"`
}

// TokenTransferResponse represents a token transfer
type TokenTransferResponse struct {
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
	CreatedAt      time.Time `json:"-"`
}

// LRC20TokenListResponse is the API response for LRC20 token list
type LRC20TokenListResponse struct {
	Tokens []LRC20TokenInfo `json:"tokens"`
	Total  int64            `json:"total"`
}

// TokenHoldersResponse is the API response for token holders list
type TokenHoldersResponse struct {
	Holders []TokenHolderResponse `json:"holders"`
	Total   int64                  `json:"total"`
}

// TokenTransfersResponse is the API response for token transfers list
type TokenTransfersResponse struct {
	Transfers []TokenTransferResponse `json:"transfers"`
	Total     int64                   `json:"total"`
}

// TokenPriceResponse represents token price information
type TokenPriceResponse struct {
	Price     float64 `json:"price"`
	Change24h float64 `json:"change_24h"`
	Volume24h float64 `json:"volume_24h"`
	MarketCap float64 `json:"market_cap"`
}

// WinkFundResponse represents WINK fund information
type WinkFundResponse struct {
	Total       int64 `json:"total"`
	Burned      int64 `json:"burned"`
	Circulating int64 `json:"circulating"`
}

// JSTFundResponse represents JST fund information
type JSTFundResponse struct {
	Total       int64 `json:"total"`
	Burned      int64 `json:"burned"`
	Circulating int64 `json:"circulating"`
}

// AssetIssueResponse represents LRC-10 asset issue information
type AssetIssueResponse struct {
	ID                       string         `json:"id"`
	OwnerAddress             string         `json:"owner_address"`
	Name                     string         `json:"name"`
	Abbr                     string         `json:"abbr"`
	TotalSupply              int64          `json:"total_supply"`
	FrozenSupply             []FrozenSupply `json:"frozen_supply,omitempty"`
	LindNum                   int32          `json:"lind_num"`
	Num                      int32          `json:"num"`
	Precision                int32          `json:"precision,omitempty"`
	StartTime                int64          `json:"start_time"`
	EndTime                  int64          `json:"end_time"`
	VoteScore                int32          `json:"vote_score"`
	Description              string         `json:"description"`
	URL                      string         `json:"url"`
	FreeAssetNetLimit        int64          `json:"free_asset_net_limit,omitempty"`
	PublicFreeAssetNetLimit  int64          `json:"public_free_asset_net_limit,omitempty"`
	PublicFreeAssetNetUsage  int64          `json:"public_free_asset_net_usage,omitempty"`
	PublicLatestFreeNetTime  int64          `json:"public_latest_free_net_time,omitempty"`
}

// FrozenSupply represents frozen supply for LRC-10 token
type FrozenSupply struct {
	FrozenAmount int64 `json:"frozen_amount"`
	FrozenDays   int64 `json:"frozen_days"`
}

// AssetIssueListResponse is the API response for multiple assets
type AssetIssueListResponse struct {
	AssetIssue []AssetIssueResponse `json:"assetIssue"`
}