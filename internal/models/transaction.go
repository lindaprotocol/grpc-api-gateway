// internal/models/transaction.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// JSON is a type alias for json.RawMessage
// type JSON json.RawMessage

// Transaction represents a blockchain transaction database model
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
	RawDataHex      string    `gorm:"type:text" json:"raw_data_hex,omitempty"`
	Signature       JSON      `gorm:"type:jsonb" json:"signature,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// CallValueInfo represents token transfer information in internal transactions (for database storage)
type CallValueInfo struct {
	CallValue int64  `json:"callValue,omitempty"`
	TokenID   string `json:"tokenId,omitempty"`
}

// InternalTransaction represents an internal transaction database model
// This matches the structure from the Linda blockchain API
type InternalTransaction struct {
	Hash              string          `gorm:"primaryKey;type:varchar(64)" json:"hash"`
	CallerAddress     string          `gorm:"index;type:varchar(42)" json:"caller_address"`
	TransferToAddress string          `gorm:"index;type:varchar(42)" json:"transferTo_address"`
	CallValueInfo     []CallValueInfo `gorm:"type:jsonb" json:"callValueInfo"`
	Note              string          `gorm:"type:text" json:"note"` // Hex string that decodes to instruction type
	Rejected          bool            `json:"rejected"`
	Extra             JSON            `gorm:"type:jsonb" json:"extra,omitempty"` // For voting details and other extra info
	BlockTimestamp    int64           `gorm:"index" json:"block_timestamp"`
	TransactionID     string          `gorm:"index;type:varchar(64)" json:"transaction_id"` // Parent transaction ID
	CreatedAt         time.Time       `json:"created_at"`
}

// CallValueInfoWrapper implements sql.Scanner and driver.Valuer for CallValueInfo slice
type CallValueInfoWrapper []CallValueInfo

// Scan implements the sql.Scanner interface
func (c *CallValueInfoWrapper) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal CallValueInfo: expected []byte, got %T", value)
	}
	
	var info []CallValueInfo
	if err := json.Unmarshal(bytes, &info); err != nil {
		return fmt.Errorf("failed to unmarshal CallValueInfo: %w", err)
	}
	
	*c = info
	return nil
}

// Value implements the driver.Valuer interface
func (c CallValueInfoWrapper) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

// Note: The note field is stored as hex string in the database
// Common note values (hex decoded):
// - "63616c6c" -> "call"
// - "667265657a6542616c616e63655632466f72456e65726779" -> "freezeBalanceV2ForEnergy"
// - "756e667265657a6542616c616e63655632466f7242616e647769647468" -> "unfreezeBalanceV2ForBandwidth"
// - "64656c65676174655265736f757263654f66456e65726779" -> "delegateResourceOfEnergy"
// - "756e44656c65676174655265736f757263654f66456e65726779" -> "unDelegateResourceOfEnergy"
// - "766f74655769746e657373" -> "voteWitness"
// - "7769746864726177526577617264" -> "withdrawReward"

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
	InternalTransactions []*InternalTransactionResponse `json:"internal_transactions,omitempty"`
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
	InternalTransactions      []*InternalTransactionResponse   `json:"internal_transactions,omitempty"`
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