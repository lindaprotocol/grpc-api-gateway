package models

import (
	"time"
)

// Account represents a blockchain account
type Account struct {
	ID                    uint      `gorm:"primarykey" json:"-"`
	Address               string    `gorm:"uniqueIndex;type:varchar(42)" json:"address"`
	Balance               int64     `json:"balance"`
	AccountName           string    `gorm:"type:varchar(100)" json:"account_name,omitempty"`
	AccountType           string    `gorm:"type:varchar(20)" json:"account_type,omitempty"`
	CreateTime            int64     `json:"create_time,omitempty"`
	IsWitness             bool      `json:"is_witness"`
	Allowance             int64     `json:"allowance,omitempty"`
	LatestWithdrawTime    int64     `json:"latest_withdraw_time,omitempty"`
	LatestOperationTime   int64     `json:"latest_operation_time,omitempty"`
	Transactions          int64     `gorm:"default:0" json:"transactions,omitempty"`
	Bandwidth             int64     `gorm:"default:0" json:"bandwidth,omitempty"`
	Energy                int64     `gorm:"default:0" json:"energy,omitempty"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// AccountResource represents account resource information
type AccountResource struct {
	Address                     string `json:"address"`
	FreeNetUsed                 int64  `json:"free_net_used"`
	FreeNetLimit                int64  `json:"free_net_limit"`
	NetUsed                     int64  `json:"net_used"`
	NetLimit                    int64  `json:"net_limit"`
	TotalNetLimit               int64  `json:"total_net_limit"`
	TotalNetWeight              int64  `json:"total_net_weight"`
	TotalLindaPowerWeight        int64  `json:"total_linda_power_weight"`
	LindaPowerLimit              int64  `json:"linda_power_limit"`
	LindaPowerUsed               int64  `json:"linda_power_used"`
	EnergyUsed                  int64  `json:"energy_used"`
	EnergyLimit                 int64  `json:"energy_limit"`
	TotalEnergyLimit            int64  `json:"total_energy_limit"`
	TotalEnergyWeight           int64  `json:"total_energy_weight"`
	AssetNetUsed                JSON   `gorm:"type:jsonb" json:"asset_net_used,omitempty"`
	AssetNetLimit               JSON   `gorm:"type:jsonb" json:"asset_net_limit,omitempty"`
}

// Frozen represents a frozen balance (Stake 1.0)
type Frozen struct {
	ID            uint      `gorm:"primarykey" json:"-"`
	Address       string    `gorm:"index;type:varchar(42)" json:"address"`
	FrozenBalance int64     `json:"frozen_balance"`
	ExpireTime    int64     `json:"expire_time"`
	ResourceType  string    `gorm:"type:varchar(20)" json:"resource_type"` // BANDWIDTH, ENERGY
	CreatedAt     time.Time `json:"created_at"`
}

// FreezeV2 represents a Stake 2.0 freeze
type FreezeV2 struct {
	ID        uint      `gorm:"primarykey" json:"-"`
	Address   string    `gorm:"index;type:varchar(42)" json:"address"`
	Amount    int64     `json:"amount"`
	Type      string    `gorm:"type:varchar(20)" json:"type"` // BANDWIDTH, ENERGY
	CreatedAt time.Time `json:"created_at"`
}

// UnFreezeV2 represents a Stake 2.0 unfreeze operation
type UnFreezeV2 struct {
	ID                 uint      `gorm:"primarykey" json:"-"`
	Address            string    `gorm:"index;type:varchar(42)" json:"address"`
	Type               string    `gorm:"type:varchar(20)" json:"type"`
	UnfreezeAmount     int64     `json:"unfreeze_amount"`
	UnfreezeExpireTime int64     `json:"unfreeze_expire_time"`
	CreatedAt          time.Time `json:"created_at"`
}

// Permission represents account permission
type Permission struct {
	ID             uint   `gorm:"primarykey" json:"-"`
	Address        string `gorm:"index;type:varchar(42)" json:"-"`
	Type           int    `json:"type"`
	PermissionID   int    `json:"id"`
	PermissionName string `gorm:"type:varchar(50)" json:"permission_name"`
	Threshold      int64  `json:"threshold"`
	ParentID       int    `json:"parent_id,omitempty"`
	Operations     string `gorm:"type:text" json:"operations,omitempty"`
	Keys           JSON   `gorm:"type:jsonb" json:"keys"`
}

// Vote represents a witness vote
type Vote struct {
	ID          uint      `gorm:"primarykey" json:"-"`
	Voter       string    `gorm:"index;type:varchar(42)" json:"voter"`
	VoteAddress string    `gorm:"index;type:varchar(42)" json:"vote_address"`
	VoteCount   int64     `json:"vote_count"`
	Timestamp   int64     `json:"timestamp"`
	CreatedAt   time.Time `json:"-"`
}

// AccountResponse is the API response for account queries
type AccountResponse struct {
	Address            string                 `json:"address"`
	Balance            int64                  `json:"balance"`
	AccountName        string                 `json:"account_name,omitempty"`
	CreateTime         int64                  `json:"create_time,omitempty"`
	IsWitness          bool                   `json:"is_witness"`
	Allowance          int64                   `json:"allowance,omitempty"`
	LatestWithdrawTime int64                   `json:"latest_withdraw_time,omitempty"`
	LatestOperationTime int64                  `json:"latest_opration_time,omitempty"`
	LatestConsumeTime  int64                   `json:"latest_consume_time,omitempty"`
	LatestConsumeFreeTime int64                `json:"latest_consume_free_time,omitempty"`
	NetWindowSize      int64                   `json:"net_window_size,omitempty"`
	NetWindowOptimized bool                    `json:"net_window_optimized,omitempty"`
	Frozen             []Frozen                `json:"frozen,omitempty"`
	DelegatedFrozenBalanceForBandwidth int64   `json:"delegated_frozen_balance_for_bandwidth,omitempty"`
	AcquiredDelegatedFrozenBalanceForBandwidth int64 `json:"acquired_delegated_frozen_balance_for_bandwidth,omitempty"`
	FrozenV2           []FreezeV2              `json:"frozenV2,omitempty"`
	UnfrozenV2         []UnFreezeV2            `json:"unfrozenV2,omitempty"`
	DelegatedFrozenV2BalanceForBandwidth int64 `json:"delegated_frozenV2_balance_for_bandwidth,omitempty"`
	AcquiredDelegatedFrozenV2BalanceForBandwidth int64 `json:"acquired_delegated_frozenV2_balance_for_bandwidth,omitempty"`
	AccountResource    *AccountResource        `json:"account_resource,omitempty"`
	OwnerPermission    *Permission             `json:"owner_permission,omitempty"`
	WitnessPermission  *Permission             `json:"witness_permission,omitempty"`
	ActivePermissions  []*Permission           `json:"active_permission,omitempty"`
	Asset              map[string]int64        `json:"asset,omitempty"`
	AssetV2            map[string]int64        `json:"assetV2,omitempty"`
	LRC20              []map[string]string     `json:"lrc20,omitempty"`
	Votes              []Vote                   `json:"votes,omitempty"`
	NetUsage           int64                    `json:"net_usage,omitempty"`
	FreeNetUsage       int64                    `json:"free_net_usage,omitempty"`
	FreeAssetNetUsageV2 map[string]int64        `json:"free_asset_net_usageV2,omitempty"`
	Tags               []TagResponse            `json:"tags,omitempty"`
}