package models

import (
	"encoding/json"
	"time"
)

// Event represents a blockchain event
type Event struct {
	ID                 uint            `gorm:"primarykey" json:"-"`
	BlockNumber        int64           `gorm:"index" json:"block_number"`
	BlockTimestamp     int64           `gorm:"index" json:"block_timestamp"`
	ContractAddress    string          `gorm:"index;type:varchar(42)" json:"contract_address"`
	EventIndex         string          `json:"event_index"`
	EventName          string          `gorm:"index;type:varchar(100)" json:"event_name"`
	EventSignature     string          `gorm:"type:varchar(256)" json:"event"`
	TransactionID      string          `gorm:"index;type:varchar(64)" json:"transaction_id"`
	Result             JSON            `gorm:"type:jsonb" json:"result"`
	ResultType         JSON            `gorm:"type:jsonb" json:"result_type"`
	Unconfirmed        bool            `gorm:"default:false" json:"_unconfirmed"`
	CreatedAt          time.Time       `json:"-"`
}

// EventResponse is the API response for event queries
type EventResponse struct {
	BlockNumber          int64                  `json:"block_number"`
	BlockTimestamp       int64                  `json:"block_timestamp"`
	CallerContractAddress string                 `json:"caller_contract_address"`
	ContractAddress      string                 `json:"contract_address"`
	EventIndex           string                 `json:"event_index"`
	EventName            string                 `json:"event_name"`
	Event                string                 `json:"event"`
	TransactionID        string                 `json:"transaction_id"`
	Result               map[string]interface{} `json:"result"`
	ResultType           map[string]string      `json:"result_type"`
	Unconfirmed          bool                   `json:"_unconfirmed,omitempty"`
}

// EventTransactionResponse represents a transaction in event responses
type EventTransactionResponse struct {
	ID              string `json:"id"`
	BlockNumber     int64  `json:"blockNumber"`
	BlockTimestamp  int64  `json:"blockTimestamp"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	Fee             int64  `json:"fee"`
	ContractAddress string `json:"contractAddress,omitempty"`
}

// EventTransactionsResponse is the API response for multiple event transactions
type EventTransactionsResponse struct {
	Data  []EventTransactionResponse `json:"data"`
	Meta  *PaginationMeta            `json:"meta"`
	Success bool                     `json:"success"`
}

// TransferResponse represents a transfer in event responses
type TransferResponse struct {
	TransactionID   string `json:"transaction_id"`
	BlockNumber     int64  `json:"block_number"`
	BlockTimestamp  int64  `json:"block_timestamp"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	TokenID         string `json:"token_id,omitempty"`
	TokenName       string `json:"token_name,omitempty"`
	TokenSymbol     string `json:"token_symbol,omitempty"`
	TokenDecimals   int32  `json:"token_decimals,omitempty"`
	ContractAddress string `json:"contract_address,omitempty"`
}

// TransferListResponse is the API response for multiple transfers
type TransferListResponse struct {
	Data  []TransferResponse `json:"data"`
	Meta  *PaginationMeta    `json:"meta"`
	Success bool             `json:"success"`
}

// EventListResponse is the API response for multiple events
type EventListResponse struct {
	Data  []EventResponse   `json:"data"`
	Meta  *PaginationMeta   `json:"meta"`
	Success bool            `json:"success"`
}

// EventBlockResponse represents a block in event responses
type EventBlockResponse struct {
	Hash              string `json:"hash"`
	Number            int64  `json:"number"`
	Timestamp         int64  `json:"timestamp"`
	ParentHash        string `json:"parentHash"`
	WitnessAddress    string `json:"witnessAddress"`
	TransactionCount  int    `json:"transactionCount"`
}

// ContractLogResponse represents a contract log
type ContractLogResponse struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      int64    `json:"blockNumber"`
	BlockTimestamp   int64    `json:"blockTimestamp"`
	TransactionID    string   `json:"transactionId"`
	TransactionIndex int      `json:"transactionIndex"`
	LogIndex         int      `json:"logIndex"`
}

// ContractLogsResponse is the API response for multiple contract logs
type ContractLogsResponse struct {
	Data  []ContractLogResponse `json:"data"`
	Meta  *PaginationMeta       `json:"meta"`
	Success bool                `json:"success"`
}

// ContractWithAbiResponse represents contract with ABI
type ContractWithAbiResponse struct {
	ContractAddress string          `json:"contractAddress"`
	ABI             json.RawMessage `json:"abi"`
	Logs            []interface{}   `json:"logs"`
}