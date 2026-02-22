package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

// ResponseInterceptor middleware for modifying responses
func ResponseInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture response
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		// Only process JSON responses
		if strings.Contains(c.Writer.Header().Get("Content-Type"), "application/json") {
			var data interface{}
			if err := json.Unmarshal(blw.body.Bytes(), &data); err == nil {
				// Process addresses in response
				processed := processJSONForAddresses(data)
				
				// Add metadata for v1 endpoints
				if strings.HasPrefix(c.Request.URL.Path, "/v1/") {
					if respMap, ok := processed.(map[string]interface{}); ok {
						if _, hasMeta := respMap["meta"]; !hasMeta {
							respMap["meta"] = gin.H{
								"at":        time.Now().UnixMilli(),
								"page_size": c.GetInt("page_size"),
							}
						}
						if _, hasSuccess := respMap["success"]; !hasSuccess {
							respMap["success"] = c.Writer.Status() < 400
						}
						processed = respMap
					}
				}

				// Marshal back
				if modified, err := json.Marshal(processed); err == nil {
					c.Writer.Header().Set("Content-Length", string(len(modified)))
					c.Writer.Write(modified)
				}
			}
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func processJSONForAddresses(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			// Check if field contains address
			if isAddressField(key) {
				if str, ok := val.(string); ok && len(str) > 0 {
					// Convert hex to base58 if needed
					if strings.HasPrefix(str, "30") && len(str) == 42 {
						if base58, err := utils.HexToBase58(str); err == nil {
							result[key] = base58
							continue
						}
					}
				}
			}
			result[key] = processJSONForAddresses(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = processJSONForAddresses(val)
		}
		return result
	default:
		return v
	}
}

func isAddressField(key string) bool {
	keyLower := strings.ToLower(key)
	addressFields := []string{
		"address", "owner_address", "to_address", "from_address",
		"vote_address", "account_address", "contract_address",
		"proposer_address", "creator_address", "witness_address",
		"receiver_address", "caller_contract_address",
	}
	for _, field := range addressFields {
		if strings.Contains(keyLower, field) {
			return true
		}
	}
	return false
}