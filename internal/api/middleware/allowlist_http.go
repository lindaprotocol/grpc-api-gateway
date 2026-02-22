package middleware

import (
	"net/http"

	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

// Allowlist validates requests against the user's allowlist (User-Agent, Origin, etc.)
func Allowlist(authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(AuthUserKey).(*auth.User)
			if !ok || user == nil || user.IsAnonymous {
				next.ServeHTTP(w, r)
				return
			}

			// Check User-Agent allowlist if configured
			if len(user.Allowlist.UserAgents) > 0 {
				ua := r.Header.Get("User-Agent")
				allowed := false
				for _, a := range user.Allowlist.UserAgents {
					if a == ua || a == "*" {
						allowed = true
						break
					}
				}
				if !allowed {
					utils.RespondWithErrorHTTP(w, http.StatusForbidden, "User-Agent not in allowlist")
					return
				}
			}

			// Check Origin allowlist if configured
			if len(user.Allowlist.Origins) > 0 {
				origin := r.Header.Get("Origin")
				allowed := false
				for _, o := range user.Allowlist.Origins {
					if o == origin || o == "*" {
						allowed = true
						break
					}
				}
				if !allowed && origin != "" {
					utils.RespondWithErrorHTTP(w, http.StatusForbidden, "Origin not in allowlist")
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
