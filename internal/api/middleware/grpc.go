package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

// CustomHeaderMatcher extracts custom headers for gRPC metadata
func CustomHeaderMatcher(key string) (string, bool) {
	key = strings.ToLower(key)
	switch key {
	case "linda-pro-api-key", "linda-pro-api-key-bin":
		return key, true
	case "authorization", "authorization-bin":
		return key, true
	case "x-request-id", "x-request-id-bin":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

// AddRequestMetadata adds request metadata to outgoing context
func AddRequestMetadata(ctx context.Context, r *http.Request) metadata.MD {
	md := make(map[string]string)

	if apiKey := r.Header.Get("LINDA-PRO-API-KEY"); apiKey != "" {
		md["linda-pro-api-key"] = apiKey
	}
	if auth := r.Header.Get("Authorization"); auth != "" {
		md["authorization"] = auth
	}
	if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
		md["x-request-id"] = reqID
	}

	return metadata.New(md)
}

// CustomErrorHandler handles gRPC errors and converts them to HTTP responses
func CustomErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}
