package models

import (
	"time"
)

// Block represents a blockchain block
type Block struct {
	Number           int64     `gorm:"primaryKey" json:"number"`
	Hash             string    `gorm:"uniqueIndex;type:varchar(64)" json:"hash"`
	ParentHash       string    `gorm:"type:varchar(64)" json:"parent_hash"`
	Timestamp        int64     `json:"timestamp"`
	WitnessAddress   string    `gorm:"index;type:varchar(42)" json:"witness_address"`
	WitnessID        int       `json:"witness_id"`
	TxTrieRoot       string    `gorm:"type:varchar(64)" json:"tx_trie_root"`
	TransactionCount int       `json:"transaction_count"`
	Size             int       `json:"size"`
	Version          int       `json:"version"`
	CreatedAt        time.Time `json:"created_at"`
}

// BlockHeader represents the header of a block
type BlockHeader struct {
	RawData          *BlockRawData `json:"raw_data"`
	WitnessSignature string        `json:"witness_signature"`
}

// BlockRawData represents the raw data of a block header
type BlockRawData struct {
	Timestamp        int64  `json:"timestamp"`
	TxTrieRoot       string `json:"txTrieRoot"`
	ParentHash       string `json:"parentHash"`
	Number           int64  `json:"number"`
	WitnessID        int64  `json:"witness_id"`
	WitnessAddress   string `json:"witness_address"`
	Version          int32  `json:"version"`
	AccountStateRoot string `json:"accountStateRoot,omitempty"`
}

// BlockResponse is the API response for block queries
type BlockResponse struct {
	BlockID     string                 `json:"blockID"`
	BlockHeader *BlockHeader           `json:"block_header"`
	Transactions []TransactionResponse `json:"transactions,omitempty"`
}

// BlockListResponse is the API response for multiple blocks
type BlockListResponse struct {
	Block []BlockResponse `json:"block"`
}

// BlockStatsResponse represents block statistics
type BlockStatsResponse struct {
	TxStat           *TxStat           `json:"txStat,omitempty"`
	FeeStat          *FeeStat          `json:"feeStat"`
}

// TxStat represents transaction statistics for a block
type TxStat struct {
	LindAnd10TransferCount   int         `json:"lindAnd10TransferCount"`
	Lrc20And721TransferCount int        `json:"lrc20And721TransferCount"`
	Lrc1155TransferCount     int        `json:"lrc1155TransferCount"`
	TransferCount            int        `json:"transferCount"`
	FailTxCount              int        `json:"failTxCount"`
	InternalTxCount          int        `json:"internalTxCount"`
	ContainInternalTxCount   int        `json:"containInternalTxCount"`
	ContractTypeDistribute   map[int]int `json:"contractTypeDistribute"`
}

// FeeStat represents fee statistics for a block
type FeeStat struct {
	NetUsage                    int64 `json:"netUsage"`
	EnergyUsage                 int64 `json:"energyUsage"`
	OtherFee                    int64 `json:"otherFee"`
	SrCandidateRegistrationFee  int64 `json:"srCandidateRegistrationFee"`
	AccountActivationFee        int64 `json:"accountActivationFee"`
	PermissionUpdateFee         int64 `json:"permissionUpdateFee"`
	MultiSignatureFee           int64 `json:"multiSignatureFee"`
	MemoFee                     int64 `json:"memoFee"`
	Lrc10AssetIssueFee          int64 `json:"lrc10AssetIssueFee"`
	DexPairCreateFee            int64 `json:"dexPairCreateFee"`
	DexOrderSellFee             int64 `json:"dexOrderSellFee"`
	DexOrderCancelFee           int64 `json:"dexOrderCancelFee"`
	EnergyBurnFeeSunAmt         int64 `json:"energyBurnFeeSunAmt"`
	BandwidthConsumedFromBurnCnt int64 `json:"bandwidthConsumedFromBurnCnt"`
	FreeBandwidthUsageCnt       int64 `json:"freeBandwidthUsageCnt"`
	BandwidthBurnFeeSunAmt      int64 `json:"bandwidthBurnFeeSunAmt"`
	EnergyConsumedFromOwnerBurnCnt int64 `json:"energyConsumedFromOwnerBurnCnt"`
	FreeEnergyUsageCnt          int64 `json:"freeEnergyUsageCnt"`
}