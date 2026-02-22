-- API Keys table
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(100) NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    daily_limit BIGINT DEFAULT 100000,
    rate_limit_qps INT DEFAULT 15,
    blocked_until TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}',
    INDEX idx_api_key (api_key),
    INDEX idx_user_id (user_id)
);

-- JWT Keys table
CREATE TABLE jwt_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(100) NOT NULL,
    key_id VARCHAR(100) UNIQUE NOT NULL,
    public_key TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    fingerprint VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    INDEX idx_key_id (key_id),
    INDEX idx_user_id (user_id)
);

-- Allowlist table
CREATE TABLE allowlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key_id UUID REFERENCES api_keys(id) ON DELETE CASCADE,
    user_id VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'user_agent', 'origin', 'contract_address', 'api_method'
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_api_key_id (api_key_id),
    INDEX idx_user_id (user_id),
    UNIQUE(api_key_id, type, value)
);

-- Rate limit tracking
CREATE TABLE rate_limit_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key_id UUID REFERENCES api_keys(id) ON DELETE CASCADE,
    user_id VARCHAR(100),
    client_ip INET,
    endpoint VARCHAR(255),
    method VARCHAR(10),
    status_code INT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_api_key_id_timestamp (api_key_id, timestamp),
    INDEX idx_user_id_timestamp (user_id, timestamp),
    INDEX idx_ip_timestamp (client_ip, timestamp)
);

-- Violations tracking
CREATE TABLE violations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key_id UUID REFERENCES api_keys(id) ON DELETE CASCADE,
    user_id VARCHAR(100),
    client_ip INET,
    violation_type VARCHAR(50), -- 'rate_limit', 'allowlist', 'auth'
    details JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_api_key_id (api_key_id),
    INDEX idx_user_id (user_id)
);

-- Blocks indexed data
CREATE TABLE blocks (
    number BIGINT PRIMARY KEY,
    hash VARCHAR(64) UNIQUE NOT NULL,
    parent_hash VARCHAR(64) NOT NULL,
    timestamp BIGINT NOT NULL,
    witness_address VARCHAR(34),
    witness_id INT,
    tx_trie_root VARCHAR(64),
    transaction_count INT,
    size INT,
    version INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_timestamp (timestamp),
    INDEX idx_witness_address (witness_address)
);

-- Transactions
CREATE TABLE transactions (
    hash VARCHAR(64) PRIMARY KEY,
    block_number BIGINT REFERENCES blocks(number) ON DELETE CASCADE,
    block_timestamp BIGINT,
    from_address VARCHAR(34),
    to_address VARCHAR(34),
    contract_address VARCHAR(34),
    amount BIGINT,
    fee BIGINT,
    energy_used BIGINT,
    energy_fee BIGINT,
    net_usage BIGINT,
    net_fee BIGINT,
    result INT,
    contract_type INT,
    data TEXT,
    raw_data TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_block_number (block_number),
    INDEX idx_from_address (from_address),
    INDEX idx_to_address (to_address),
    INDEX idx_contract_address (contract_address),
    INDEX idx_block_timestamp (block_timestamp)
);

-- LRC-20/LRC20 Tokens
CREATE TABLE lrc20_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_address VARCHAR(34) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    decimals INT DEFAULT 18,
    total_supply VARCHAR(100),
    owner_address VARCHAR(34),
    issue_time BIGINT,
    holders_count BIGINT DEFAULT 0,
    transfer_count BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    INDEX idx_symbol (symbol),
    INDEX idx_owner (owner_address)
);

-- Token holders
CREATE TABLE token_holders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_address VARCHAR(34) NOT NULL REFERENCES lrc20_tokens(contract_address) ON DELETE CASCADE,
    address VARCHAR(34) NOT NULL,
    balance VARCHAR(100) NOT NULL,
    percentage DECIMAL(10,4),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(contract_address, address),
    INDEX idx_contract_address (contract_address),
    INDEX idx_address (address)
);

-- Tags
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    address VARCHAR(34) NOT NULL,
    tag VARCHAR(100) NOT NULL,
    description TEXT,
    owner VARCHAR(34) NOT NULL,
    signature TEXT,
    votes INT DEFAULT 0,
    created_at BIGINT NOT NULL,
    INDEX idx_address (address),
    INDEX idx_tag (tag),
    INDEX idx_owner (owner)
);

-- Statistics
CREATE TABLE statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL,
    value JSONB NOT NULL,
    timestamp BIGINT NOT NULL,
    INDEX idx_type_timestamp (type, timestamp)
);

-- Market data
CREATE TABLE market_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pair VARCHAR(20) NOT NULL,
    price VARCHAR(100),
    volume_24h VARCHAR(100),
    high_24h VARCHAR(100),
    low_24h VARCHAR(100),
    change_24h VARCHAR(20),
    timestamp BIGINT NOT NULL,
    INDEX idx_pair_timestamp (pair, timestamp)
);