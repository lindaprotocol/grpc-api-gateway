package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

// Recovery recovers from panics and returns 500
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic recovered: %v\n%s", err, debug.Stack())
					utils.RespondWithErrorHTTP(w, http.StatusInternalServerError, "Internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
