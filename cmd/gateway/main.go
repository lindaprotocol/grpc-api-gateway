// cmd/gateway/main.go
package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lindaprotocol/grpc-api-gateway/internal/api/middleware"
	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/cache"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/lindascan"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/postgres"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	configPath = flag.String("config", "./internal/config/config.yaml", "configuration file path")
)

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	// Initialize database
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		panic(err)
	}

	// Initialize Redis cache - use config.Redis directly
	redisCache, err := cache.NewRedisClient(cache.RedisConfig{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
	})
	if err != nil {
		panic(err)
	}

	// Initialize auth service
	authService := auth.NewService(cfg.Auth, db, redisCache)

	// Create gRPC connection to blockchain nodes
	conn, err := grpc.Dial(
		cfg.Linda.FullnodeEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(32*1024*1024)),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	solidityConn, err := grpc.Dial(
		cfg.Linda.SolidityEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(32*1024*1024)),
	)
	if err != nil {
		panic(err)
	}
	defer solidityConn.Close()

	// Initialize blockchain clients
	blockchainClient := blockchain.NewClient(conn, solidityConn, cfg.Linda)

	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db)
	blockRepo := repository.NewBlockRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	eventRepo := repository.NewEventRepository(db)
	tagRepo := repository.NewTagRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	// Initialize Lindascan service
	lindascanService := lindascan.NewService(
		blockchainClient,
		accountRepo,
		blockRepo,
		txRepo,
		tokenRepo,
		eventRepo,
		tagRepo,
		statsRepo,
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create gateway mux with custom options
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   false,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithIncomingHeaderMatcher(middleware.CustomHeaderMatcher),
		runtime.WithMetadata(middleware.AddRequestMetadata),
		runtime.WithErrorHandler(middleware.CustomErrorHandler),
	)

	// Register all services
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)),
	}

	// Register Wallet service
	if err := lindapb.RegisterWalletHandlerFromEndpoint(ctx, gwmux, cfg.Linda.FullnodeEndpoint, opts); err != nil {
		panic(err)
	}

	// Register WalletSolidity service
	if err := lindapb.RegisterWalletSolidityHandlerFromEndpoint(ctx, gwmux, cfg.Linda.SolidityEndpoint, opts); err != nil {
		panic(err)
	}

	// Register JsonRpc service
	if err := lindapb.RegisterJsonRpcHandlerFromEndpoint(ctx, gwmux, cfg.Linda.FullnodeEndpoint, opts); err != nil {
		panic(err)
	}

	// Register EventService
	if err := lindapb.RegisterEventServiceHandlerFromEndpoint(ctx, gwmux, cfg.Linda.EventEndpoint, opts); err != nil {
		panic(err)
	}

	// Register Lindascan custom service
	if err := lindapb.RegisterLindascanHandlerServer(ctx, gwmux, lindascanService); err != nil {
		panic(err)
	}

	// Build middleware chain
	handler := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	}).Handler(gwmux)

	// Auth middleware (API Key & JWT)
	handler = middleware.Auth(authService, cfg.Auth)(handler)

	// Rate limiting middleware - convert config types
	rateLimitConfig := middleware.RateLimitConfig{
		Enabled:      cfg.RateLimit.Enabled,
		DefaultQPS:   cfg.RateLimit.DefaultQPS,
		DefaultBurst: cfg.RateLimit.DefaultBurst,
		Strategy:     cfg.RateLimit.Strategy,
		Store:        cfg.RateLimit.Store,
	}
	handler = middleware.RateLimit(redisCache.Client(), rateLimitConfig)(handler)

	// Allowlist middleware
	handler = middleware.Allowlist(authService)(handler)

	// Response interceptor - Fix the gin handler adapter
	responseInterceptor := middleware.ResponseInterceptor()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a Gin context
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = r
		// Call the interceptor and then continue with the next handler
		responseInterceptor(ginCtx)
		// The next handler in the chain is already called by the middleware
		// So we don't need to call handler.ServeHTTP here
	})

	// Logging middleware
	handler = middleware.Logger(cfg.Logging)(handler)

	// Recovery middleware
	handler = middleware.Recovery()(handler)

	// Start server
	server := &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}