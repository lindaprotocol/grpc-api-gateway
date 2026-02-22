// scripts/migrate.go
package main

import (
    "log"
    
    "github.com/lindaprotocol/grpc-api-gateway/internal/config"
    "github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/postgres"
)

func main() {
    // Load configuration
    cfg, err := config.Load("./internal/config/config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Connect to database (this will auto-migrate)
    _, err = postgres.NewConnection(cfg.Database)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    log.Println("Database migrations completed successfully")
}