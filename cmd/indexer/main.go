package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/indexer"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/postgres"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	configPath = flag.String("config", "./internal/config/config.yaml", "configuration file path")
)

func main() {
	flag.Parse()

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

	// Create gRPC connection
	conn, err := grpc.Dial(
		cfg.Linda.FullnodeEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to fullnode: %v", err)
	}
	defer conn.Close()

	solidityConn, err := grpc.Dial(
		cfg.Linda.SolidityEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to solidity node: %v", err)
	}
	defer solidityConn.Close()

	// Initialize blockchain client
	blockchainClient := blockchain.NewClient(conn, solidityConn, cfg.Linda)

	// Initialize repositories
	blockRepo := repository.NewBlockRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	eventRepo := repository.NewEventRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	// Initialize indexer
	idx := indexer.NewIndexer(
		&cfg.Indexer,
		blockchainClient,
		blockRepo,
		txRepo,
		tokenRepo,
		eventRepo,
		statsRepo,
	)

	// Start indexer
	if err := idx.Start(); err != nil {
		log.Fatalf("Failed to start indexer: %v", err)
	}

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Stop indexer
	if err := idx.Stop(); err != nil {
		log.Printf("Error stopping indexer: %v", err)
	}
}