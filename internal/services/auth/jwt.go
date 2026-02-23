package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JWTKey struct {
	ID          string    `gorm:"primaryKey;type:uuid"`
	UserID      string    `gorm:"index;not null"`
	KeyID       string    `gorm:"uniqueIndex;not null"`
	PublicKey   string    `gorm:"type:text;not null"`
	Name        string    `gorm:"not null"`
	Fingerprint string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
	ExpiresAt   *time.Time
	IsActive    bool      `gorm:"default:true"`
}

func (JWTKey) TableName() string {
	return "jwt_keys"
}

type JWTService struct {
	db        *gorm.DB
	secretKey []byte
}

func NewJWTService(db *gorm.DB, secretKey string) *JWTService {
	return &JWTService{
		db:        db,
		secretKey: []byte(secretKey),
	}
}

// GenerateToken creates a new JWT token
func (s *JWTService) GenerateToken(userID string, keyID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"kid": keyID,
		"aud": "lindagrid.io",
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID

	// Sign with private key - you'd need to load the private key
	// For now, we'll use HMAC for simplicity
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secretKey)
}

// ValidateToken validates a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	// Extract key ID from header
	if kid, ok := token.Header["kid"].(string); ok {
		// Validate that the key exists and is active
		var key JWTKey
		if err := s.db.Where("key_id = ? AND is_active = ?", kid, true).First(&key).Error; err != nil {
			return nil, jwt.ErrTokenInvalidId
		}
	}

	return token, nil
}

// RegisterPublicKey registers a new public key for JWT
func (s *JWTService) RegisterPublicKey(userID string, publicKey string, name string) (*JWTKey, error) {
	// Generate key ID
	keyID := uuid.New().String()

	// Calculate fingerprint (simplified)
	fingerprint := uuid.New().String()[:8]

	key := &JWTKey{
		ID:          uuid.New().String(),
		UserID:      userID,
		KeyID:       keyID,
		PublicKey:   publicKey,
		Name:        name,
		Fingerprint: fingerprint,
		CreatedAt:   time.Now(),
		IsActive:    true,
	}

	if err := s.db.Create(key).Error; err != nil {
		return nil, err
	}

	return key, nil
}

// RevokeKey deactivates a JWT key
func (s *JWTService) RevokeKey(keyID string) error {
	return s.db.Model(&JWTKey{}).Where("key_id = ?", keyID).Update("is_active", false).Error
}

// GetUserKeys returns all JWT keys for a user
func (s *JWTService) GetUserKeys(userID string) ([]JWTKey, error) {
	var keys []JWTKey
	err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&keys).Error
	return keys, err
}