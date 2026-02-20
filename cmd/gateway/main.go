package main

import (
    "bytes"
    "context"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "flag"
    "math/big"
    "net/http"
    "strings"

    "github.com/golang/glog"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    "github.com/lindaprotocol/grpc-api-gateway/pkg/api/protocol"
)

var (
    grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
    httpPort           = flag.String("http-port", ":18890", "HTTP port")
)

// Base58 alphabet (Bitcoin/Linda style)
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// encodeBase58Check converts bytes to base58check string (with version byte and checksum)
func encodeBase58Check(input []byte) string {
    if len(input) == 0 {
        return ""
    }

    // Calculate checksum (first 4 bytes of double SHA256)
    firstSHA := sha256.Sum256(input)
    secondSHA := sha256.Sum256(firstSHA[:])
    checksum := secondSHA[:4]

    // Append checksum to input
    payload := append(input, checksum...)

    // Count leading zeros
    zeros := 0
    for zeros < len(payload) && payload[zeros] == 0 {
        zeros++
    }

    // Convert to big integer
    num := new(big.Int).SetBytes(payload)

    // Convert to base58
    result := make([]byte, 0)
    base := big.NewInt(58)
    zero := big.NewInt(0)
    rem := new(big.Int)

    for num.Cmp(zero) > 0 {
        num.DivMod(num, base, rem)
        result = append([]byte{base58Alphabet[rem.Int64()]}, result...)
    }

    // Add leading zeros as '1's
    for i := 0; i < zeros; i++ {
        result = append([]byte{base58Alphabet[0]}, result...)
    }

    return string(result)
}

// hexEncodeBytes converts bytes to hex string
func hexEncodeBytes(data []byte) string {
    if data == nil || len(data) == 0 {
        return ""
    }
    return hex.EncodeToString(data)
}

// responseWriter captures the response for modification
type responseWriter struct {
    http.ResponseWriter
    body       *bytes.Buffer
    statusCode int
}

func (w *responseWriter) Write(b []byte) (int, error) {
    return w.body.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
}

// isAddressField checks if a field name indicates it contains an address
func isAddressField(key string) bool {
    keyLower := strings.ToLower(key)
    return strings.Contains(keyLower, "address") || 
           keyLower == "witness" ||
           keyLower == "owneraddress" ||
           keyLower == "toaddress" ||
           keyLower == "fromaddress"
}

// needsHexConversion checks if a field should be converted to hex
func needsHexConversion(key string) bool {
    keyLower := strings.ToLower(key)
    
    // Fields that should be hex
    hexFields := []string{
        "parenthash", "txtrieroot", "witnesssignature", "signature",
        "refblockbytes", "refblockhash", "txid", "blockid", "hash",
        "data", "script", "proof", "merkle",
    }
    
    for _, field := range hexFields {
        if strings.Contains(keyLower, field) {
            return true
        }
    }
    return false
}

// decodeBase64String attempts to decode a string from various base64 formats
func decodeBase64String(s string) ([]byte, error) {
    // Try standard base64
    bytes, err := base64.StdEncoding.DecodeString(s)
    if err == nil {
        return bytes, nil
    }
    // Try URL-safe base64
    bytes, err = base64.URLEncoding.DecodeString(s)
    if err == nil {
        return bytes, nil
    }
    // Try raw base64 (without padding)
    bytes, err = base64.RawStdEncoding.DecodeString(s)
    if err == nil {
        return bytes, nil
    }
    return nil, err
}

// processJSON recursively processes JSON to convert bytes appropriately
func processJSON(obj interface{}) interface{} {
    switch v := obj.(type) {
    case map[string]interface{}:
        result := make(map[string]interface{})
        for key, val := range v {
            // Handle signature array specially
            if key == "signature" {
                // Check if it's an array of signatures
                if sigArray, ok := val.([]interface{}); ok {
                    convertedSigs := make([]interface{}, len(sigArray))
                    for i, sig := range sigArray {
                        if sigStr, ok := sig.(string); ok {
                            // Decode base64 to bytes, then encode as hex
                            if bytes, err := decodeBase64String(sigStr); err == nil {
                                convertedSigs[i] = hexEncodeBytes(bytes)
                            } else {
                                convertedSigs[i] = sig
                            }
                        } else {
                            convertedSigs[i] = processJSON(sig)
                        }
                    }
                    result[key] = convertedSigs
                    continue
                }
            }
            
            // Handle string values that need conversion
            if str, ok := val.(string); ok && len(str) > 0 {
                // Check if this is an address field
                if isAddressField(key) {
                    // Decode from base64 to bytes, then encode as base58check
                    if bytes, err := decodeBase64String(str); err == nil {
                        result[key] = encodeBase58Check(bytes)
                    } else {
                        result[key] = val
                    }
                
                // Check if this needs hex conversion
                } else if needsHexConversion(key) {
                    // Decode from base64 to bytes, then encode as hex
                    if bytes, err := decodeBase64String(str); err == nil {
                        result[key] = hexEncodeBytes(bytes)
                    } else {
                        result[key] = val
                    }
                
                // Handle all other fields
                } else {
                    result[key] = val
                }
            } else {
                // Recursively process non-string values
                result[key] = processJSON(val)
            }
        }
        return result
        
    case []interface{}:
        result := make([]interface{}, len(v))
        for i, val := range v {
            result[i] = processJSON(val)
        }
        return result
        
    default:
        return v
    }
}

// responseInterceptor intercepts and modifies responses
func responseInterceptor(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Create a custom response writer
        rw := &responseWriter{
            ResponseWriter: w,
            body:           &bytes.Buffer{},
            statusCode:     http.StatusOK,
        }

        // Serve the request
        h.ServeHTTP(rw, r)

        // Parse and modify the response if it's JSON
        contentType := rw.Header().Get("Content-Type")
        if strings.Contains(contentType, "application/json") || rw.body.Len() > 0 {
            var data interface{}
            if err := json.Unmarshal(rw.body.Bytes(), &data); err == nil {
                // Process the JSON
                processed := processJSON(data)
                
                // Marshal back with indentation
                if modified, err := json.MarshalIndent(processed, "", "  "); err == nil {
                    w.Header().Set("Content-Type", "application/json")
                    w.WriteHeader(rw.statusCode)
                    w.Write(modified)
                    return
                }
            }
        }

        // If anything fails, write the original response
        w.Header().Set("Content-Type", rw.Header().Get("Content-Type"))
        w.WriteHeader(rw.statusCode)
        w.Write(rw.body.Bytes())
    })
}

func corsMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        h.ServeHTTP(w, r)
    })
}

func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

    // Register all services
    services := []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
        protocol.RegisterWalletHandlerFromEndpoint,
        protocol.RegisterWalletSolidityHandlerFromEndpoint,
        protocol.RegisterWalletExtensionHandlerFromEndpoint,
        protocol.RegisterScanServiceHandlerFromEndpoint,
    }

    for _, register := range services {
        if err := register(ctx, mux, *grpcServerEndpoint, opts); err != nil {
            return err
        }
    }

    // Wrap with CORS and response interceptor
    handler := corsMiddleware(mux)
    handler = responseInterceptor(handler)

    glog.Infof("HTTP server listening on %s", *httpPort)
    return http.ListenAndServe(*httpPort, handler)
}

func main() {
    flag.Parse()
    defer glog.Flush()

    if err := run(); err != nil {
        glog.Fatal(err)
    }
}