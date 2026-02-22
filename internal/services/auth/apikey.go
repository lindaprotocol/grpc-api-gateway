package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
	
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/models"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIKey struct {
	ID          string     `gorm:"primaryKey;type:uuid"`
	UserID      string     `gorm:"index;not null"`
	Key         string     `gorm:"uniqueIndex;not null"`
	KeyHash     string     `gorm:"not null"`
	Name        string     `gorm:"not null"`
	CreatedAt   time.Time  `gorm:"not null"`
	ExpiresAt   *time.Time `gorm:"index"`
	LastUsedAt  *time.Time `gorm:"index"`
	IsActive    bool       `gorm:"default:true"`
	DailyLimit  int64      `gorm:"default:100000"`
	RateLimitQPS int       `gorm:"default:15"`
	BlockedUntil *time.Time
	Metadata    JSON       `gorm:"type:jsonb"`
}

func (APIKey) TableName() string {
	return "api_keys"
}

type APIKeyService struct {
	db *gorm.DB
}

func NewAPIKeyService(db *gorm.DB) *APIKeyService {
	return &APIKeyService{db: db}
}

// GenerateAPIKey creates a new API key
func (s *APIKeyService) GenerateAPIKey(userID, name string, dailyLimit int64, rateLimitQPS int) (*APIKey, string, error) {
	// Generate random key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, "", err
	}
	plainKey := hex.EncodeToString(keyBytes)

	// Hash the key for storage
	keyHash, err := bcrypt.GenerateFromPassword([]byte(plainKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	apiKey := &APIKey{
		ID:          uuid.New().String(),
		UserID:      userID,
		Key:         plainKey, // Store hash instead in production
		KeyHash:     string(keyHash),
		Name:        name,
		CreatedAt:   time.Now(),
		IsActive:    true,
		DailyLimit:  dailyLimit,
		RateLimitQPS: rateLimitQPS,
		Metadata:    make(JSON),
	}

	if err := s.db.Create(apiKey).Error; err != nil {
		return nil, "", err
	}

	return apiKey, plainKey, nil
}

// ValidateAPIKey checks if an API key is valid
func (s *APIKeyService) ValidateAPIKey(apiKey string) (*APIKey, error) {
	var key APIKey
	if err := s.db.Where("key = ?", apiKey).First(&key).Error; err != nil {
		return nil, err
	}

	// Check if active
	if !key.IsActive {
		return nil, ErrKeyInactive
	}

	// Check expiration
	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return nil, ErrKeyExpired
	}

	// Check if blocked
	if key.BlockedUntil != nil && key.BlockedUntil.After(time.Now()) {
		return nil, ErrKeyBlocked
	}

	// Update last used
	now := time.Now()
	key.LastUsedAt = &now
	s.db.Model(&key).Update("last_used_at", now)

	return &key, nil
}

// RevokeAPIKey deactivates an API key
func (s *APIKeyService) RevokeAPIKey(keyID string) error {
	return s.db.Model(&APIKey{}).Where("id = ?", keyID).Update("is_active", false).Error
}

// GetAPIKeys returns all keys for a user
func (s *APIKeyService) GetAPIKeys(userID string) ([]APIKey, error) {
	var keys []APIKey
	err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&keys).Error
	return keys, err
}

// BlockAPIKey temporarily blocks a key
func (s *APIKeyService) BlockAPIKey(keyID string, duration time.Duration) error {
	blockedUntil := time.Now().Add(duration)
	return s.db.Model(&APIKey{}).Where("id = ?", keyID).Update("blocked_until", blockedUntil).Error
}

// UpdateKeyLimits updates rate limits for a key
func (s *APIKeyService) UpdateKeyLimits(keyID string, dailyLimit int64, rateLimitQPS int) error {
	return s.db.Model(&APIKey{}).Where("id = ?", keyID).Updates(map[string]interface{}{
		"daily_limit":   dailyLimit,
		"rate_limit_qps": rateLimitQPS,
	}).Error
}

// TrackRequest records an API request for rate limiting
func (s *APIKeyService) TrackRequest(keyID string) error {
	// Implement request tracking for daily limits
	// This would typically use Redis for counters
	return nil
}