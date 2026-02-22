package middleware

import (
	"net/http"
)

// ResponseInterceptor wraps the next handler for response modification (e.g. address conversion).
// For gRPC-gateway, proto marshaling is done by the runtime; this is a pass-through for now.
func ResponseInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
