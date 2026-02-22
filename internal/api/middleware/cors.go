package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CorsMiddleware handles CORS headers
func CorsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		if len(allowedOrigins) > 0 {
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
				// Check wildcard subdomain patterns
				if strings.Contains(o, "*") {
					pattern := strings.ReplaceAll(o, "*", "")
					if strings.HasSuffix(origin, pattern) {
						allowed = true
						break
					}
				}
			}
		} else {
			allowed = true
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, LINDA-PRO-API-KEY, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}