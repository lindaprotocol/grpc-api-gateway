// internal/services/auth/service.go
package auth

import (
    "errors"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/lindaprotocol/grpc-api-gateway/internal/config"
    "github.com/lindaprotocol/grpc-api-gateway/internal/services/cache"
    "gorm.io/gorm"
)

var (
    ErrKeyInactive = errors.New("API key is inactive")
    ErrKeyExpired  = errors.New("API key has expired")
    ErrKeyBlocked  = errors.New("API key is temporarily blocked")
)

// Service struct: Service for authentication
type Service struct {
    db            *gorm.DB
    cache         *cache.RedisClient
    config        config.AuthConfig
    apiKeyService *APIKeyService
    jwtService    *JWTService
    allowlist     *AllowlistService
}

type User struct {
    ID           string
    IsAnonymous  bool
    APIKeyID     string
    DailyLimit   int64
    RateLimit    int
    BlockedUntil *time.Time
    Allowlist    AllowlistData
}

type AllowlistData struct {
    UserAgents        []string
    Origins           []string
    ContractAddresses []string
    APIMethods        []string
}

func NewService(cfg config.AuthConfig, db *gorm.DB, cache *cache.RedisClient) *Service {
    return &Service{
        db:            db,
        cache:         cache,
        config:        cfg,
        apiKeyService: NewAPIKeyService(db),
        jwtService:    NewJWTService(db, cfg.JWTSecret),
        allowlist:     NewAllowlistService(db),
    }
}

// ValidateAPIKey function: Validates an API key and returns the associated user
func (s *Service) ValidateAPIKey(apiKey string) (*User, error) {
    // Check cache first
    var user User
    if err := s.cache.Get("apikey:"+apiKey, &user); err == nil {
        return &user, nil
    }

    // Validate in database
    key, err := s.apiKeyService.ValidateAPIKey(apiKey)
    if err != nil {
        return nil, err
    }

    // Get allowlist
    allowlist, err := s.allowlist.GetAllowlist(key.ID)
    if err != nil {
        return nil, err
    }

    user = User{
        ID:           key.UserID,
        APIKeyID:     key.ID,
        DailyLimit:   key.DailyLimit,
        RateLimit:    key.RateLimitQPS,
        BlockedUntil: key.BlockedUntil,
        Allowlist: AllowlistData{
            UserAgents:        allowlist["user_agent"],
            Origins:           allowlist["origin"],
            ContractAddresses: allowlist["contract_address"],
            APIMethods:        allowlist["api_method"],
        },
    }

    // Cache for 5 minutes
    s.cache.Set("apikey:"+apiKey, user, 5*time.Minute)

    return &user, nil
}

// ValidateJWT function: Validates a JWT token
func (s *Service) ValidateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
    token, err := s.jwtService.ValidateToken(tokenString)
    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        // Convert to RegisteredClaims
        rc := &jwt.RegisteredClaims{
            Subject: claims["sub"].(string),
        }
        return rc, nil
    }

    return nil, jwt.ErrTokenInvalidClaims
}

// GetUserByID function: Retrieves a user by ID
func (s *Service) GetUserByID(userID string) (*User, error) {
    // For now, return a basic user
    return &User{
        ID:         userID,
        DailyLimit: s.config.DefaultDailyLimit,
        RateLimit:  s.config.DefaultRateLimitQPS,
    }, nil
}

// TrackRequest function: Tracks an API request for rate limiting
func (s *Service) TrackRequest(user *User, r *http.Request) error {
    if user.IsAnonymous {
        // Track by IP
        ip := r.RemoteAddr
        return s.trackRateLimit("ip:"+ip, s.config.UnauthenticatedRateLimitQPS)
    }

    // Track by API key
    return s.trackRateLimit("apikey:"+user.APIKeyID, user.RateLimit)
}

func (s *Service) trackRateLimit(key string, limit int) error {
    // Implement rate limiting counters in Redis
    // This would increment counters and check against limits
    return nil
}

// CheckViolation function: Checks if a user has violated rate limits
func (s *Service) CheckViolation(user *User) bool {
    // Check violation count in Redis
    violationKey := "violation:" + user.APIKeyID
    count, _ := s.cache.Incr(violationKey)
    s.cache.Expire(violationKey, 24*time.Hour)

    return count >= 3
}

// BlockUser function: Temporarily blocks a user
func (s *Service) BlockUser(user *User, duration time.Duration) error {
    // Simply pass the duration to BlockAPIKey, don't create an unused variable
    return s.apiKeyService.BlockAPIKey(user.APIKeyID, duration)
}