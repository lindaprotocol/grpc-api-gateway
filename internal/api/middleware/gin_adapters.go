package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/cache"
)

// AuthMiddleware is the Gin adapter for the Auth HTTP middleware
func AuthMiddleware(authService *auth.Service, cfg config.AuthConfig) gin.HandlerFunc {
	httpMiddleware := Auth(authService, cfg)
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Propagate request with auth context and copy user to Gin context
			c.Request = r
			if user := r.Context().Value(AuthUserKey); user != nil {
				c.Set("auth_user", user)
			}
			c.Next()
		})
		httpMiddleware(next).ServeHTTP(c.Writer, c.Request)
	}
}

// RateLimitMiddleware is the Gin adapter for the RateLimit HTTP middleware
func RateLimitMiddleware(cacheClient *cache.RedisClient, cfg config.RateLimitConfig) gin.HandlerFunc {
	httpMiddleware := RateLimit(cacheClient.Client(), RateLimitConfig{
		Enabled:      cfg.Enabled,
		DefaultQPS:   cfg.DefaultQPS,
		DefaultBurst: cfg.DefaultBurst,
		Strategy:     cfg.Strategy,
		Store:        cfg.Store,
	})
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			c.Next()
		})
		httpMiddleware(next).ServeHTTP(c.Writer, c.Request)
	}
}

// AllowlistMiddleware is the Gin adapter for the Allowlist HTTP middleware
func AllowlistMiddleware(authService *auth.Service) gin.HandlerFunc {
	httpMiddleware := Allowlist(authService)
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			c.Next()
		})
		httpMiddleware(next).ServeHTTP(c.Writer, c.Request)
	}
}
