package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Allowlist struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	APIKeyID  string    `gorm:"index;not null"`
	UserID    string    `gorm:"index;not null"`
	Type      string    `gorm:"not null"` // user_agent, origin, contract_address, api_method
	Value     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (Allowlist) TableName() string {
	return "allowlists"
}

type AllowlistService struct {
	db *gorm.DB
}

func NewAllowlistService(db *gorm.DB) *AllowlistService {
	return &AllowlistService{db: db}
}

// AddToAllowlist adds an entry to the allowlist
func (s *AllowlistService) AddToAllowlist(apiKeyID, userID, allowlistType, value string) error {
	entry := &Allowlist{
		ID:        uuid.New().String(),
		APIKeyID:  apiKeyID,
		UserID:    userID,
		Type:      allowlistType,
		Value:     value,
		CreatedAt: time.Now(),
	}
	return s.db.Create(entry).Error
}

// RemoveFromAllowlist removes an entry from the allowlist
func (s *AllowlistService) RemoveFromAllowlist(entryID string) error {
	return s.db.Delete(&Allowlist{}, "id = ?", entryID).Error
}

// GetAllowlist returns all allowlist entries for an API key
func (s *AllowlistService) GetAllowlist(apiKeyID string) (map[string][]string, error) {
	var entries []Allowlist
	if err := s.db.Where("api_key_id = ?", apiKeyID).Find(&entries).Error; err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for _, entry := range entries {
		result[entry.Type] = append(result[entry.Type], entry.Value)
	}
	return result, nil
}

// CheckUserAgent validates if a user agent is allowed
func (s *AllowlistService) CheckUserAgent(apiKeyID, userAgent string) bool {
	if apiKeyID == "" {
		return true
	}

	var count int64
	s.db.Model(&Allowlist{}).
		Where("api_key_id = ? AND type = ?", apiKeyID, "user_agent").
		Count(&count)

	if count == 0 {
		return true // No restrictions
	}

	var allowed bool
	s.db.Raw(`
		SELECT EXISTS(
			SELECT 1 FROM allowlists 
			WHERE api_key_id = ? 
			AND type = 'user_agent' 
			AND ? LIKE '%' || value || '%'
		)`, apiKeyID, userAgent).Scan(&allowed)

	return allowed
}

// CheckOrigin validates if an origin is allowed
func (s *AllowlistService) CheckOrigin(apiKeyID, origin string) bool {
	if apiKeyID == "" {
		return true
	}

	var count int64
	s.db.Model(&Allowlist{}).
		Where("api_key_id = ? AND type = ?", apiKeyID, "origin").
		Count(&count)

	if count == 0 {
		return true
	}

	var allowed bool
	s.db.Raw(`
		SELECT EXISTS(
			SELECT 1 FROM allowlists 
			WHERE api_key_id = ? 
			AND type = 'origin' 
			AND ? LIKE value
		)`, apiKeyID, origin).Scan(&allowed)

	// Check wildcard patterns
	if !allowed {
		s.db.Raw(`
			SELECT EXISTS(
				SELECT 1 FROM allowlists 
				WHERE api_key_id = ? 
				AND type = 'origin' 
				AND ? LIKE REPLACE(value, '*', '%')
			)`, apiKeyID, origin).Scan(&allowed)
	}

	return allowed
}

// CheckContractAddress validates if a contract address is allowed
func (s *AllowlistService) CheckContractAddress(apiKeyID, address string) bool {
	if apiKeyID == "" {
		return true
	}

	var count int64
	s.db.Model(&Allowlist{}).
		Where("api_key_id = ? AND type = ?", apiKeyID, "contract_address").
		Count(&count)

	if count == 0 {
		return true
	}

	var allowed bool
	s.db.Raw(`
		SELECT EXISTS(
			SELECT 1 FROM allowlists 
			WHERE api_key_id = ? 
			AND type = 'contract_address' 
			AND value = ?
		)`, apiKeyID, address).Scan(&allowed)

	return allowed
}

// CheckAPIMethod validates if an API method is allowed
func (s *AllowlistService) CheckAPIMethod(apiKeyID, method string) bool {
	if apiKeyID == "" {
		return true
	}

	var count int64
	s.db.Model(&Allowlist{}).
		Where("api_key_id = ? AND type = ?", apiKeyID, "api_method").
		Count(&count)

	if count == 0 {
		return true
	}

	var allowed bool
	s.db.Raw(`
		SELECT EXISTS(
			SELECT 1 FROM allowlists 
			WHERE api_key_id = ? 
			AND type = 'api_method' 
			AND value = ?
		)`, apiKeyID, method).Scan(&allowed)

	return allowed
}