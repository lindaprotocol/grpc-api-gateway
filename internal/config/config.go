package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Environment string            `yaml:"environment"`
	Linda       LindaConfig       `yaml:"linda"`
	Database    DatabaseConfig    `yaml:"database"`
	Redis       RedisConfig       `yaml:"redis"`
	Auth        AuthConfig        `yaml:"auth"`
	RateLimit   RateLimitConfig   `yaml:"rate_limit"`
	Allowlist   AllowlistConfig   `yaml:"allowlist"`
	CORS        CORSConfig        `yaml:"cors"`
	Cache       CacheConfig       `yaml:"cache"`
	Indexer     IndexerConfig     `yaml:"indexer"`
	Logging     LoggingConfig     `yaml:"logging"`
	Tracing     TracingConfig     `yaml:"tracing"`
	Metrics     MetricsConfig     `yaml:"metrics"`
	External    ExternalAPIConfig `yaml:"external_apis"`
}

type ServerConfig struct {
	HTTPPort   string `yaml:"http_port"`
	GRPCPort   string `yaml:"grpc_port"`
	EnableTLS  bool   `yaml:"enable_tls"`
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
}

type LindaConfig struct {
	FullnodeEndpoint string        `yaml:"fullnode_endpoint"`
	SolidityEndpoint string        `yaml:"solidity_endpoint"`
	EventEndpoint    string        `yaml:"event_endpoint"`
	GRPCTimeout      time.Duration `yaml:"grpc_timeout"`
	MaxMsgSize       int           `yaml:"max_msg_size"`
}

type DatabaseConfig struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	SSLMode         string `yaml:"sslmode"`
	MaxConnections  int    `yaml:"max_connections"`
	IdleConnections int    `yaml:"idle_connections"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	MaxRetries   int    `yaml:"max_retries"`
}

type AuthConfig struct {
	APIKeyEnabled               bool   `yaml:"api_key_enabled"`
	JWTEnabled                  bool   `yaml:"jwt_enabled"`
	JWTSecret                   string `yaml:"jwt_secret"`
	JWTExpiry                   int64  `yaml:"jwt_expiry"`
	MaxKeysPerAccount           int    `yaml:"max_keys_per_account"`
	DefaultRateLimitQPS         int    `yaml:"default_rate_limit_qps"`
	DefaultDailyLimit           int64  `yaml:"default_daily_limit"`
	UnauthenticatedRateLimitQPS int    `yaml:"unauthenticated_rate_limit_qps"`
	UnauthenticatedDailyLimit   int64  `yaml:"unauthenticated_daily_limit"`
	PenaltyDuration             int    `yaml:"penalty_duration"`
	AllowAnonymous              bool   `yaml:"allow_anonymous"`
}

type RateLimitConfig struct {
	Enabled      bool   `yaml:"enabled"`
	DefaultQPS   int    `yaml:"default_qps"`
	DefaultBurst int    `yaml:"default_burst"`
	Strategy     string `yaml:"strategy"`
	Store        string `yaml:"store"`
}

type AllowlistConfig struct {
	Enabled               bool `yaml:"enabled"`
	UserAgentEnabled      bool `yaml:"user_agent_enabled"`
	OriginEnabled         bool `yaml:"origin_enabled"`
	ContractAddressEnabled bool `yaml:"contract_address_enabled"`
	APIMethodEnabled      bool `yaml:"api_method_enabled"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type CacheConfig struct {
	DefaultTTL   int `yaml:"default_ttl"`
	AccountTTL   int `yaml:"account_ttl"`
	BlockTTL     int `yaml:"block_ttl"`
	TransactionTTL int `yaml:"transaction_ttl"`
	TokenTTL     int `yaml:"token_ttl"`
	StatsTTL     int `yaml:"stats_ttl"`
}

type IndexerConfig struct {
	Enabled            bool          `yaml:"enabled"`
	BlockBatchSize     int           `yaml:"block_batch_size"`
	TransactionBatchSize int         `yaml:"transaction_batch_size"`
	ContractBatchSize  int           `yaml:"contract_batch_size"`
	SyncInterval       time.Duration `yaml:"sync_interval"`
	StartBlock         int64         `yaml:"start_block"`
	MaxWorkers         int           `yaml:"max_workers"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
	Traces bool   `yaml:"traces"`
}

type TracingConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AgentHost string `yaml:"agent_host"`
	AgentPort int    `yaml:"agent_port"`
}

type MetricsConfig struct {
	Enabled           bool   `yaml:"enabled"`
	PrometheusEndpoint string `yaml:"prometheus_endpoint"`
	Port              int    `yaml:"port"`
}

type ExternalAPIConfig struct {
	CoinMarketCap CoinMarketCapConfig `yaml:"coinmarketcap"`
	IPGeo         IPGeoConfig         `yaml:"ipgeo"`
}

type CoinMarketCapConfig struct {
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url"`
}

type IPGeoConfig struct {
	Provider      string `yaml:"provider"`
	APIKey        string `yaml:"api_key"`
	MaxMindDBPath string `yaml:"maxmind_db_path"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}