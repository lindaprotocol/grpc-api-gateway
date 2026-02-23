// internal/models/transaction.go
package models

import (
	"time"
)

// Transaction represents a blockchain transaction database model
type Transaction struct {
	ID              uint      `gorm:"primarykey" json:"-"`
	Hash            string    `gorm:"uniqueIndex;type:varchar(64)" json:"hash"`
	BlockNumber     int64     `gorm:"index" json:"block_number"`
	BlockTimestamp  int64     `gorm:"index" json:"block_timestamp"`
	FromAddress     string    `gorm:"index;type:varchar(42)" json:"from_address"`
	ToAddress       string    `gorm:"index;type:varchar(42)" json:"to_address"`
	ContractAddress string    `gorm:"index;type:varchar(42)" json:"contract_address,omitempty"`
	Amount          int64     `json:"amount,omitempty"`
	Fee             int64     `json:"fee,omitempty"`
	EnergyUsed      int64     `json:"energy_used,omitempty"`
	EnergyFee       int64     `json:"energy_fee,omitempty"`
	NetUsage        int64     `json:"net_usage,omitempty"`
	NetFee          int64     `json:"net_fee,omitempty"`
	Result          int       `json:"result,omitempty"` // 0: success, 1: failed
	ContractType    int       `json:"contract_type,omitempty"`
	Data            string    `gorm:"type:text" json:"data,omitempty"`
	RawData         string    `gorm:"type:text" json:"raw_data,omitempty"` // This will store the hex string
	RawDataHex      string    `gorm:"type:text" json:"raw_data_hex,omitempty"` // This field for hex
	Signature       JSON      `gorm:"type:jsonb" json:"signature,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}