// internal/models/stats.go
package models

import (
	"encoding/json"
)

// Statistic represents a stored statistic database model
type Statistic struct {
	ID        uint            `gorm:"primarykey" json:"-"`
	Type      string          `gorm:"index;type:varchar(50)" json:"type"`
	Value     json.RawMessage `gorm:"type:jsonb" json:"value"`
	Timestamp int64           `gorm:"index" json:"timestamp"`
}