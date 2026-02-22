module github.com/lindaprotocol/grpc-api-gateway

go 1.26

// Resolve ambiguous import: chaincfg/chainhash exists in both btcd and as standalone module
exclude github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1

require (
    gopkg.in/DataDog/dd-trace-go.v1 v1.61.0
    github.com/btcsuite/btcutil v1.0.2
    github.com/ethereum/go-ethereum v1.13.4
    github.com/gin-gonic/gin v1.9.1
    github.com/go-redis/redis/v8 v8.11.5
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/golang/protobuf v1.5.3
    github.com/google/uuid v1.3.1
    github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.0
    github.com/rs/cors v1.10.1
    github.com/sirupsen/logrus v1.9.3
    golang.org/x/crypto v0.14.0
    google.golang.org/genproto/googleapis/api v0.0.0-20231016165738-49dd2c1f3d0b
    google.golang.org/grpc v1.59.0
    google.golang.org/protobuf v1.31.0
    gopkg.in/yaml.v3 v3.0.1
    gorm.io/driver/postgres v1.5.3
    gorm.io/gorm v1.25.5
)

require (
    github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
    github.com/bytedance/sonic v1.9.1 // indirect
    github.com/cespare/xxhash/v2 v2.2.0 // indirect
    github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
    github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
    github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
    github.com/gabriel-vasile/mimetype v1.4.2 // indirect
    github.com/gin-contrib/sse v0.1.0 // indirect
    github.com/go-playground/locales v0.14.1 // indirect
    github.com/go-playground/universal-translator v0.18.1 // indirect
    github.com/go-playground/validator/v10 v10.14.0 // indirect
    github.com/goccy/go-json v0.10.2 // indirect
    github.com/holiman/uint256 v1.2.3 // indirect
    github.com/jackc/pgpassfile v1.0.0 // indirect
    github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
    github.com/jackc/pgx/v5 v5.4.3 // indirect
    github.com/jinzhu/inflection v1.0.0 // indirect
    github.com/jinzhu/now v1.1.5 // indirect
    github.com/json-iterator/go v1.1.12 // indirect
    github.com/klauspost/cpuid/v2 v2.2.4 // indirect
    github.com/leodido/go-urn v1.2.4 // indirect
    github.com/mattn/go-isatty v0.0.19 // indirect
    github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
    github.com/modern-go/reflect2 v1.0.2 // indirect
    github.com/pelletier/go-toml/v2 v2.0.8 // indirect
    github.com/rogpeppe/go-internal v1.11.0 // indirect
    github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
    github.com/ugorji/go/codec v1.2.11 // indirect
    golang.org/x/arch v0.3.0 // indirect
    golang.org/x/net v0.17.0 // indirect
    golang.org/x/sys v0.13.0 // indirect
    golang.org/x/text v0.13.0 // indirect
    google.golang.org/genproto v0.0.0-20231016165738-49dd2c1f3d0b // indirect
    google.golang.org/genproto/googleapis/rpc v0.0.0-20231016165738-49dd2c1f3d0b // indirect
)