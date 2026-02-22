package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"github.com/lindaprotocol/grpc-api-gateway/internal/api/middleware"
	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/cache"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/postgres"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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

	// Initialize tracer if enabled
	if cfg.Tracing.Enabled {
		tracer.Start(
			tracer.WithServiceName("grpc-api-gateway"),
			tracer.WithEnv(cfg.Environment),
		)
		defer tracer.Stop()
	}

	// Initialize database
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		panic(err)
	}

	// Initialize Redis cache
	redisCache, err := cache.NewRedisClientFromConfig(cfg.Redis)
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create gateway mux with custom options
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions:   protojson.MarshalOptions{UseProtoNames: false, EmitUnpopulated: true},
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
		}),
		runtime.WithIncomingHeaderMatcher(middleware.CustomHeaderMatcher),
		runtime.WithMetadata(middleware.AddRequestMetadata),
		runtime.WithErrorHandler(middleware.CustomErrorHandler),
	)

	// Register all services
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(32*1024*1024)),
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
	if err := lindapb.RegisterLindascanHandlerServer(ctx, gwmux, blockchainClient); err != nil {
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

	// Rate limiting middleware
	handler = middleware.RateLimit(redisCache.Client(), middleware.RateLimitConfig{
		Enabled:     cfg.RateLimit.Enabled,
		DefaultQPS:  cfg.RateLimit.DefaultQPS,
		DefaultBurst: cfg.RateLimit.DefaultBurst,
		Strategy:    cfg.RateLimit.Strategy,
		Store:       cfg.RateLimit.Store,
	})(handler)

	// Allowlist middleware
	handler = middleware.Allowlist(authService)(handler)

	// Response interceptor for address conversion
	handler = middleware.ResponseInterceptor(handler)

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