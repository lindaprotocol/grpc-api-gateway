package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/postgres"
)

var (
	configPath = flag.String("config", "./internal/config/config.yaml", "configuration file path")
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	apiKeyService := auth.NewAPIKeyService(db)

	switch os.Args[1] {
	case "generate":
		if len(os.Args) < 5 {
			fmt.Println("Usage: apikey generate <user_id> <name> [daily_limit] [rate_limit_qps]")
			return
		}
		userID := os.Args[2]
		name := os.Args[3]
		
		dailyLimit := int64(100000)
		if len(os.Args) > 4 {
			fmt.Sscanf(os.Args[4], "%d", &dailyLimit)
		}
		
		rateLimitQPS := 15
		if len(os.Args) > 5 {
			fmt.Sscanf(os.Args[5], "%d", &rateLimitQPS)
		}

		key, plainKey, err := apiKeyService.GenerateAPIKey(userID, name, dailyLimit, rateLimitQPS)
		if err != nil {
			log.Fatalf("Failed to generate API key: %v", err)
		}

		fmt.Printf("API Key generated successfully:\n")
		fmt.Printf("ID: %s\n", key.ID)
		fmt.Printf("Key: %s\n", plainKey)
		fmt.Printf("Name: %s\n", key.Name)
		fmt.Printf("Daily Limit: %d\n", key.DailyLimit)
		fmt.Printf("Rate Limit: %d QPS\n", key.RateLimitQPS)

	case "revoke":
		if len(os.Args) < 3 {
			fmt.Println("Usage: apikey revoke <key_id>")
			return
		}
		keyID := os.Args[2]
		
		if err := apiKeyService.RevokeAPIKey(keyID); err != nil {
			log.Fatalf("Failed to revoke API key: %v", err)
		}
		fmt.Printf("API key %s revoked successfully\n", keyID)

	case "list":
		if len(os.Args) < 3 {
			fmt.Println("Usage: apikey list <user_id>")
			return
		}
		userID := os.Args[2]
		
		keys, err := apiKeyService.GetAPIKeys(userID)
		if err != nil {
			log.Fatalf("Failed to list API keys: %v", err)
		}

		fmt.Printf("API Keys for user %s:\n", userID)
		for _, key := range keys {
			status := "active"
			if !key.IsActive {
				status = "inactive"
			}
			if key.BlockedUntil != nil && key.BlockedUntil.After(key.CreatedAt) {
				status = "blocked"
			}
			fmt.Printf("  %s: %s (%s) - %s\n", key.ID[:8], key.Name, key.Key[:16], status)
		}

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`API Key Management Tool

Usage:
  apikey generate <user_id> <name> [daily_limit] [rate_limit_qps]  Generate a new API key
  apikey revoke <key_id>                                            Revoke an API key
  apikey list <user_id>                                             List API keys for a user
`)
}