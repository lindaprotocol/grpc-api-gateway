package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type contextKey string

const (
	AuthUserKey    contextKey = "auth_user"
	AuthAPIKeyKey  contextKey = "api_key"
	AuthJWTKey     contextKey = "jwt"
	AuthPermissions contextKey = "permissions"
)

func Auth(authService *auth.Service, cfg config.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract API key from header
			apiKey := r.Header.Get("LINDA-PRO-API-KEY")
			if apiKey == "" {
				// Try query parameter
				apiKey = r.URL.Query().Get("api_key")
			}

			// Extract JWT from Authorization header
			jwtToken := extractJWT(r)

			// If neither API key nor JWT provided, use anonymous user
			if apiKey == "" && jwtToken == "" {
				if !cfg.AllowAnonymous {
					utils.RespondWithErrorHTTP(w, http.StatusUnauthorized, "API key or JWT required")
					return
				}
				// Anonymous user with strict rate limits
				ctx := context.WithValue(r.Context(), AuthUserKey, &auth.User{
					IsAnonymous: true,
					RateLimit:   cfg.UnauthenticatedRateLimitQPS,
					DailyLimit:  cfg.UnauthenticatedDailyLimit,
				})
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			var user *auth.User
			var err error

			// Validate API key if provided
			if apiKey != "" {
				user, err = authService.ValidateAPIKey(apiKey)
				if err != nil {
					utils.RespondWithErrorHTTP(w, http.StatusUnauthorized, "Invalid API key")
					return
				}
			}

			// Validate JWT if provided
			if jwtToken != "" {
				// If both API key and JWT, ensure they match
				if user != nil {
					claims, err := authService.ValidateJWT(jwtToken)
					if err != nil || claims.Subject != user.ID {
						utils.RespondWithErrorHTTP(w, http.StatusUnauthorized, "JWT does not match API key")
						return
					}
				} else {
					claims, err := authService.ValidateJWT(jwtToken)
					if err != nil {
						utils.RespondWithErrorHTTP(w, http.StatusUnauthorized, "Invalid JWT")
						return
					}
					user, err = authService.GetUserByID(claims.Subject)
					if err != nil {
						utils.RespondWithErrorHTTP(w, http.StatusUnauthorized, "User not found")
						return
					}
				}
			}

			// Check if user is blocked
			if user.BlockedUntil != nil && user.BlockedUntil.After(time.Now()) {
				utils.RespondWithErrorHTTP(w, http.StatusForbidden, "Your API key is temporarily blocked due to rate limit violation")
				return
			}

			// Update rate limit counters
			if err := authService.TrackRequest(user, r); err != nil {
				utils.RespondWithErrorHTTP(w, http.StatusTooManyRequests, "Rate limit exceeded")
				return
			}

			// Store user in context
			ctx := context.WithValue(r.Context(), AuthUserKey, user)
			ctx = context.WithValue(ctx, AuthAPIKeyKey, apiKey)
			ctx = context.WithValue(ctx, AuthJWTKey, jwtToken)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractJWT(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}