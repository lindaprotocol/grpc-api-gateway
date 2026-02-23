// internal/models/event.go
package models

import (
	"time"
)

// Event represents a blockchain event database model
type Event struct {
	ID              uint      `gorm:"primarykey" json:"-"`
	BlockNumber     int64     `gorm:"index" json:"block_number"`
	BlockTimestamp  int64     `gorm:"index" json:"block_timestamp"`
	ContractAddress string    `gorm:"index;type:varchar(42)" json:"contract_address"`
	EventIndex      string    `json:"event_index"`
	EventName       string    `gorm:"index;type:varchar(100)" json:"event_name"`
	EventSignature  string    `gorm:"type:varchar(256)" json:"event_signature"`
	TransactionID   string    `gorm:"index;type:varchar(64)" json:"transaction_id"`
	Result          JSON      `gorm:"type:jsonb" json:"result"`
	ResultType      JSON      `gorm:"type:jsonb" json:"result_type"`
	Unconfirmed     bool      `gorm:"default:false" json:"unconfirmed"`
	CreatedAt       time.Time `json:"created_at"`
}