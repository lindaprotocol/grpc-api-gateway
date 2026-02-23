package utils

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "math/big"
    "strings"

    "github.com/btcsuite/btcutil/base58"
)

const (
    AddressPrefix = byte(0x30) // Linda address prefix
)

// Base58ToHex converts a base58 address to hex format
func Base58ToHex(base58Addr string) (string, error) {
    if base58Addr == "" {
        return "", errors.New("empty address")
    }

    // Decode base58
    decoded := base58.Decode(base58Addr)
    if len(decoded) == 0 {
        return "", errors.New("invalid base58 address")
    }

    // Remove checksum (last 4 bytes)
    if len(decoded) < 4 {
        return "", errors.New("address too short")
    }
    hexBytes := decoded[:len(decoded)-4]

    // Ensure prefix is correct (30 for Linda)
    if len(hexBytes) == 21 && hexBytes[0] != AddressPrefix {
        // Add prefix if missing
        hexBytes = append([]byte{AddressPrefix}, hexBytes...)
    }

    return hex.EncodeToString(hexBytes), nil
}

// HexToBase58 converts a hex address to base58 format
func HexToBase58(hexAddr string) (string, error) {
    if hexAddr == "" {
        return "", errors.New("empty address")
    }

    // Decode hex
    bytes, err := hex.DecodeString(hexAddr)
    if err != nil {
        return "", err
    }

    // Ensure correct length (21 bytes with prefix)
    if len(bytes) != 21 {
        if len(bytes) == 20 {
            // Add prefix
            bytes = append([]byte{AddressPrefix}, bytes...)
        } else {
            return "", errors.New("invalid address length")
        }
    }

    // Calculate checksum (double SHA256, first 4 bytes)
    firstSHA := sha256.Sum256(bytes)
    secondSHA := sha256.Sum256(firstSHA[:])
    checksum := secondSHA[:4]

    // Append checksum
    payload := append(bytes, checksum...)

    // Encode to base58
    return base58.Encode(payload), nil
}

// MustHexToBase58 converts hex to base58 or returns original on error
func MustHexToBase58(hexAddr string) string {
    base58, err := HexToBase58(hexAddr)
    if err != nil {
        return hexAddr
    }
    return base58
}

// IsValidBase58Address checks if a base58 address is valid
func IsValidBase58Address(addr string) bool {
    if len(addr) < 30 || len(addr) > 40 {
        return false
    }

    decoded := base58.Decode(addr)
    if len(decoded) != 25 { // 21 bytes address + 4 bytes checksum
        return false
    }

    // Verify checksum
    payload := decoded[:21]
    checksum := decoded[21:]

    firstSHA := sha256.Sum256(payload)
    secondSHA := sha256.Sum256(firstSHA[:])
    expectedChecksum := secondSHA[:4]

    for i := 0; i < 4; i++ {
        if checksum[i] != expectedChecksum[i] {
            return false
        }
    }

    return true
}

// IsValidHexAddress checks if a hex address is valid
func IsValidHexAddress(addr string) bool {
    if len(addr) != 42 { // 30 prefix + 40 hex chars
        return false
    }

    if !strings.HasPrefix(addr, "30") {
        return false
    }

    bytes, err := hex.DecodeString(addr)
    if err != nil || len(bytes) != 21 {
        return false
    }

    return true
}

// NormalizeAddress converts any address format to hex with prefix
func NormalizeAddress(addr string) (string, error) {
    if addr == "" {
        return "", errors.New("empty address")
    }

    // Check if it's already hex with prefix
    if IsValidHexAddress(addr) {
        return addr, nil
    }

    // Try base58
    if IsValidBase58Address(addr) {
        return Base58ToHex(addr)
    }

    return "", errors.New("invalid address format")
}

// TruncateString truncates a string to the specified length
func TruncateString(s string, length int) string {
    if len(s) <= length {
        return s
    }
    return s[:length] + "..."
}

// VerifySignature verifies a signature for data
func VerifySignature(address, signature, data string) (bool, error) {
    // Implement signature verification
    // This would verify that the signature was created by the owner of the address
    return true, nil
}

// EncodeBase58Check encodes bytes to base58check string
func EncodeBase58Check(input []byte) string {
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

    // Base58 alphabet
    alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

    // Convert to base58
    result := make([]byte, 0)
    base := big.NewInt(58)
    zero := big.NewInt(0)
    rem := new(big.Int)

    for num.Cmp(zero) > 0 {
        num.DivMod(num, base, rem)
        result = append([]byte{alphabet[rem.Int64()]}, result...)
    }

    // Add leading zeros as '1's
    for i := 0; i < zeros; i++ {
        result = append([]byte{alphabet[0]}, result...)
    }

    return string(result)
}