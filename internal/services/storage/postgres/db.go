package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection creates a new database connection
func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.IdleConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// Auto migrate schemas
	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

// autoMigrate runs database migrations
func autoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Account related tables
	if err := db.AutoMigrate(
		&models.Account{},
		&models.AccountResource{},
		&models.Frozen{},
		&models.FreezeV2{},
		&models.UnFreezeV2{},
		&models.Permission{},
		&models.Vote{},
	); err != nil {
		return err
	}

	// Block and transaction tables
	if err := db.AutoMigrate(
		&models.Block{},
		&models.Transaction{},
		&models.InternalTransaction{},
	); err != nil {
		return err
	}

	// Token related tables
	if err := db.AutoMigrate(
		&models.TokenInfo{},           // LRC-10 token database model
		&models.LRC20TokenInfo{},       // LRC20 token database model
		&models.TokenHolder{},          // Token holder database model
		&models.TokenTransferDB{},       // Token transfer database model AssetIssueDB
		&models.AssetIssueDB{},         // Asset issue database model
		// &models.TokenTransferResponse{}, // NULL THIS - it's an API response type
		// &models.AssetIssueResponse{},   // NULL THIS - it's an API response type
	); err != nil {
		return err
	}

	// Event tables
	if err := db.AutoMigrate(
		&models.Event{},
	); err != nil {
		return err
	}

	// Stats tables
	if err := db.AutoMigrate(
		&models.Statistic{},
	); err != nil {
		return err
	}

	// Tag tables
	if err := db.AutoMigrate(
		&models.TagResponse{},
	); err != nil {
		return err
	}

	// Auth tables
	if err := db.AutoMigrate(
		&auth.APIKey{},
		&auth.JWTKey{},
		&auth.Allowlist{},
	); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}