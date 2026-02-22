package models

import (
	"encoding/json"
	"time"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	ID              uint      `gorm:"primarykey" json:"-"`
	Hash            string    `gorm:"uniqueIndex;type:varchar(64)" json:"hash"`
	BlockNumber     int64     `gorm:"index" json:"block_number"`
	BlockTimestamp  int64     `gorm:"index" json:"block_timestamp"`
	FromAddress     string    `gorm:"index;type:varchar(42)" json:"from_address"`
	ToAddress       string    `gorm:"index;type:varchar(42)" json:"to_address"`
	ContractAddress string    `gorm:"index;type:varchar(42)" json:"contract_address,omitempty"`
	Amount          int64     `json:"amount,omitempty"`
	Fee             int64     `json:"fee,omitempty"`
	EnergyUsed      int64     `json:"energy_used,omitempty"`
	EnergyFee       int64     `json:"energy_fee,omitempty"`
	NetUsage        int64     `json:"net_usage,omitempty"`
	NetFee          int64     `json:"net_fee,omitempty"`
	Result          int       `json:"result,omitempty"` // 0: success, 1: failed
	ContractType    int       `json:"contract_type,omitempty"`
	Data            string    `gorm:"type:text" json:"data,omitempty"`
	RawData         string    `gorm:"type:text" json:"raw_data,omitempty"`
	Signature       JSON      `gorm:"type:jsonb" json:"signature,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// TransactionResponse is the API response for transaction queries
type TransactionResponse struct {
	TxID              string                 `json:"txID"`
	BlockNumber       int64                   `json:"blockNumber,omitempty"`
	BlockTimestamp    int64                   `json:"block_timestamp,omitempty"`
	Ret               []TransactionResult     `json:"ret,omitempty"`
	Signature         []string                `json:"signature,omitempty"`
	RawDataHex        string                  `json:"raw_data_hex,omitempty"`
	RawData           *TransactionRawData     `json:"raw_data,omitempty"`
	EnergyFee         int64                    `json:"energy_fee,omitempty"`
	EnergyUsage       int64                    `json:"energy_usage,omitempty"`
	EnergyUsageTotal  int64                    `json:"energy_usage_total,omitempty"`
	NetFee            int64                    `json:"net_fee,omitempty"`
	NetUsage          int64                    `json:"net_usage,omitempty"`
	InternalTransactions []*InternalTransaction `json:"internal_transactions,omitempty"`
	FeeLimit          int64                    `json:"fee_limit,omitempty"`
	RefBlockBytes     string                   `json:"ref_block_bytes,omitempty"`
	RefBlockHash      string                   `json:"ref_block_hash,omitempty"`
	Expiration        int64                    `json:"expiration,omitempty"`
	Timestamp         int64                    `json:"timestamp,omitempty"`
}

// TransactionResult represents the result of a transaction
type TransactionResult struct {
	ContractRet string `json:"contractRet"`
	Fee         int64  `json:"fee,omitempty"`
}

// TransactionRawData represents the raw data of a transaction
type TransactionRawData struct {
	Contract      []TransactionContract `json:"contract"`
	RefBlockBytes string                 `json:"ref_block_bytes"`
	RefBlockNum   int64                  `json:"ref_block_num"`
	RefBlockHash  string                 `json:"ref_block_hash"`
	Expiration    int64                  `json:"expiration"`
	Auths         []Authority            `json:"auths,omitempty"`
	Data          string                 `json:"data,omitempty"`
	Scripts       string                 `json:"scripts,omitempty"`
	Timestamp     int64                  `json:"timestamp"`
}

// TransactionContract represents a contract in a transaction
type TransactionContract struct {
	Type      string          `json:"type"`
	Parameter json.RawMessage `json:"parameter"`
}

// Authority represents an authority in a transaction
type Authority struct {
	Account *AccountID `json:"account"`
}

// AccountID represents an account ID
type AccountID struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

// TransactionInfoResponse represents transaction info response
type TransactionInfoResponse struct {
	ID                        string                 `json:"id"`
	Fee                       int64                   `json:"fee"`
	BlockNumber               int64                   `json:"blockNumber"`
	BlockTimeStamp            int64                   `json:"blockTimeStamp"`
	ContractResult            []string                `json:"contractResult"`
	ContractAddress           string                  `json:"contract_address"`
	Receipt                   *ResourceReceipt        `json:"receipt"`
	Log                       []*EventLog              `json:"log"`
	Result                    int                     `json:"result,omitempty"`
	ResMessage                string                  `json:"resMessage,omitempty"`
	AssetIssueID              string                  `json:"assetIssueID,omitempty"`
	WithdrawAmount            int64                   `json:"withdraw_amount,omitempty"`
	UnfreezeAmount            int64                   `json:"unfreeze_amount,omitempty"`
	InternalTransactions      []*InternalTransaction   `json:"internal_transactions,omitempty"`
	WithdrawExpireAmount      int64                   `json:"withdraw_expire_amount,omitempty"`
	CancelUnfreezeV2Amount    map[string]int64        `json:"cancel_unfreezeV2_amount,omitempty"`
	ExchangeReceivedAmount    int64                   `json:"exchange_received_amount,omitempty"`
	ExchangeInjectAnotherAmount int64                 `json:"exchange_inject_another_amount,omitempty"`
	ExchangeWithdrawAnotherAmount int64               `json:"exchange_withdraw_another_amount,omitempty"`
	ExchangeID                int64                   `json:"exchange_id,omitempty"`
	ShieldedTransactionFee    int64                   `json:"shielded_transaction_fee,omitempty"`
}

// ResourceReceipt represents resource receipt in transaction info
type ResourceReceipt struct {
	EnergyUsage        int64  `json:"energy_usage"`
	EnergyFee          int64  `json:"energy_fee"`
	OriginEnergyUsage  int64  `json:"origin_energy_usage"`
	EnergyUsageTotal   int64  `json:"energy_usage_total"`
	NetUsage           int64  `json:"net_usage"`
	NetFee             int64  `json:"net_fee"`
	Result             string `json:"result"`
	EnergyPenaltyTotal int64  `json:"energy_penalty_total,omitempty"`
}

// EventLog represents an event log in transaction info
type EventLog struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

// InternalTransaction represents an internal transaction
type InternalTransaction struct {
	InternalTxID    string                 `json:"internal_tx_id"`
	Data            *InternalTransactionData `json:"data"`
	BlockTimestamp  int64                   `json:"block_timestamp"`
	ToAddress       string                   `json:"to_address"`
	TxID            string                   `json:"tx_id"`
	FromAddress     string                   `json:"from_address"`
}

// InternalTransactionData represents internal transaction data
type InternalTransactionData struct {
	Note     string `json:"note"`
	Rejected bool   `json:"rejected"`
}