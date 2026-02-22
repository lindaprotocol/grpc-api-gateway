-- Create additional indexes for performance

-- Account indexes
CREATE INDEX IF NOT EXISTS idx_accounts_balance ON accounts(balance DESC);
CREATE INDEX IF NOT EXISTS idx_accounts_witness ON accounts(is_witness) WHERE is_witness = true;
CREATE INDEX IF NOT EXISTS idx_frozen_expire ON frozen(expire_time);
CREATE INDEX IF NOT EXISTS idx_unfreeze_v2_expire ON unfreeze_v2s(unfreeze_expire_time);

-- Block indexes
CREATE INDEX IF NOT EXISTS idx_blocks_timestamp ON blocks(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_blocks_witness ON blocks(witness_address);

-- Transaction indexes
CREATE INDEX IF NOT EXISTS idx_transactions_from_to ON transactions(from_address, to_address);
CREATE INDEX IF NOT EXISTS idx_transactions_contract ON transactions(contract_address);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(contract_type);
CREATE INDEX IF NOT EXISTS idx_transactions_result ON transactions(result);

-- Token indexes
CREATE INDEX IF NOT EXISTS idx_lrc20_tokens_symbol ON lrc20_token_infos(symbol);
CREATE INDEX IF NOT EXISTS idx_token_holders_balance ON token_holders(balance DESC);
CREATE INDEX IF NOT EXISTS idx_token_transfers_token_time ON token_transfer_responses(token_address, block_timestamp DESC);

-- Event indexes
CREATE INDEX IF NOT EXISTS idx_events_contract_name_time ON events(contract_address, event_name, block_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_events_tx_id ON events(transaction_id);
CREATE INDEX IF NOT EXISTS idx_events_time ON events(block_timestamp DESC);

-- Stats indexes
CREATE INDEX IF NOT EXISTS idx_statistics_type_time ON statistics(type, timestamp DESC);

-- Tag indexes
CREATE INDEX IF NOT EXISTS idx_tags_address_votes ON tags(address, votes DESC);
CREATE INDEX IF NOT EXISTS idx_tags_tag ON tags(tag);

-- Auth indexes
CREATE INDEX IF NOT EXISTS idx_api_keys_user ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_last_used ON api_keys(last_used_at) WHERE last_used_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_jwt_keys_user ON jwt_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_allowlists_key_type ON allowlists(api_key_id, type);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_transactions_address_time ON transactions(from_address, block_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_events_contract_time ON events(contract_address, block_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_token_holders_contract_balance ON token_holders(contract_address, balance DESC);

-- Full text search indexes
CREATE INDEX IF NOT EXISTS idx_lrc20_tokens_name_trgm ON lrc20_token_infos USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_lrc20_tokens_symbol_trgm ON lrc20_token_infos USING gin (symbol gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_tags_tag_trgm ON tags USING gin (tag gin_trgm_ops);

-- Partition by time for large tables (example for transactions)
-- CREATE TABLE transactions_partitioned (
--     LIKE transactions INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES
-- ) PARTITION BY RANGE (block_timestamp);
-- 
-- CREATE TABLE transactions_2024 PARTITION OF transactions_partitioned
--     FOR VALUES FROM (1704067200000) TO (1735689600000);