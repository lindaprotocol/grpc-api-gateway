// internal/models/block.go
package models

import (
	"time"
)

// Block represents a blockchain block database model
type Block struct {
	Number           int64     `gorm:"primaryKey" json:"number"`
	Hash             string    `gorm:"uniqueIndex;type:varchar(64)" json:"hash"`
	ParentHash       string    `gorm:"type:varchar(64)" json:"parent_hash"`
	Timestamp        int64     `json:"timestamp"`
	WitnessAddress   string    `gorm:"index;type:varchar(42)" json:"witness_address"`
	WitnessID        int       `json:"witness_id"`
	TxTrieRoot       string    `gorm:"type:varchar(64)" json:"tx_trie_root"`
	TransactionCount int       `json:"transaction_count"`
	Size             int       `json:"size"`
	Version          int       `json:"version"`
	CreatedAt        time.Time `json:"created_at"`
}