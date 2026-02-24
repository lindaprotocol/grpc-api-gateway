// internal/models/api_responses.go
package models

import "encoding/json"

// JSON is a type alias for json.RawMessage
// type JSON json.RawMessage

// ==================== Account API Responses ====================
type AccountResponse struct {
	Address            string                 `json:"address"`
	Balance            int64                  `json:"balance"`
	AccountName        string                 `json:"account_name,omitempty"`
	CreateTime         int64                  `json:"create_time,omitempty"`
	IsWitness          bool                   `json:"is_witness"`
	Allowance          int64                  `json:"allowance,omitempty"`
	LatestWithdrawTime int64                  `json:"latest_withdraw_time,omitempty"`
	LatestOperationTime int64                 `json:"latest_opration_time,omitempty"`
	LatestConsumeTime  int64                  `json:"latest_consume_time,omitempty"`
	LatestConsumeFreeTime int64               `json:"latest_consume_free_time,omitempty"`
	NetWindowSize      int64                  `json:"net_window_size,omitempty"`
	NetWindowOptimized bool                   `json:"net_window_optimized,omitempty"`
	Frozen             []Frozen               `json:"frozen,omitempty"`
	DelegatedFrozenBalanceForBandwidth int64  `json:"delegated_frozen_balance_for_bandwidth,omitempty"`
	AcquiredDelegatedFrozenBalanceForBandwidth int64 `json:"acquired_delegated_frozen_balance_for_bandwidth,omitempty"`
	FrozenV2           []FreezeV2             `json:"frozenV2,omitempty"`
	UnfrozenV2         []UnFreezeV2           `json:"unfrozenV2,omitempty"`
	DelegatedFrozenV2BalanceForBandwidth int64 `json:"delegated_frozenV2_balance_for_bandwidth,omitempty"`
	AcquiredDelegatedFrozenV2BalanceForBandwidth int64 `json:"acquired_delegated_frozenV2_balance_for_bandwidth,omitempty"`
	AccountResource    *AccountResource       `json:"account_resource,omitempty"`
	OwnerPermission    *Permission            `json:"owner_permission,omitempty"`
	WitnessPermission  *Permission            `json:"witness_permission,omitempty"`
	ActivePermissions  []*Permission          `json:"active_permission,omitempty"`
	Asset              map[string]int64       `json:"asset,omitempty"`
	AssetV2            map[string]int64       `json:"assetV2,omitempty"`
	LRC20              []map[string]string    `json:"lrc20,omitempty"`
	Votes              []Vote                 `json:"votes,omitempty"`
	NetUsage           int64                   `json:"net_usage,omitempty"`
	FreeNetUsage       int64                   `json:"free_net_usage,omitempty"`
	FreeAssetNetUsageV2 map[string]int64       `json:"free_asset_net_usageV2,omitempty"`
	Tags               []TagResponse           `json:"tags,omitempty"`
}

// ==================== Transaction API Responses ====================
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

type TransactionResult struct {
	ContractRet string `json:"contractRet"`
	Fee         int64  `json:"fee,omitempty"`
}

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

type TransactionContract struct {
	Type      string          `json:"type"`
	Parameter json.RawMessage `json:"parameter"`
}

type Authority struct {
	Account *AccountID `json:"account"`
}

type AccountID struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

// InternalTransactionData represents the data of an internal transaction in API responses
type InternalTransactionData struct {
	Note     string `json:"note"`
	Rejected bool   `json:"rejected"`
}

// InternalTransactionResponse is the API response for internal transactions
type InternalTransactionResponse struct {
	InternalTxID    string                 `json:"internal_tx_id"`
	Data            *InternalTransactionData `json:"data"`
	BlockTimestamp  int64                   `json:"block_timestamp"`
	ToAddress       string                   `json:"to_address"`
	TxID            string                   `json:"tx_id"`
	FromAddress     string                   `json:"from_address"`
}

// ==================== TransactionInfo API Responses ====================
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

type EventLog struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

// ==================== Block API Responses ====================
type BlockResponse struct {
	BlockID     string                 `json:"blockID"`
	BlockHeader *BlockHeader           `json:"block_header"`
	Transactions []TransactionResponse `json:"transactions,omitempty"`
}

type BlockHeader struct {
	RawData          *BlockRawData `json:"raw_data"`
	WitnessSignature string        `json:"witness_signature"`
}

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

type BlockListResponse struct {
	Block []BlockResponse `json:"block"`
}

// ==================== Block Statistics API Responses ====================
type BlockStatsResponse struct {
	TxStat           *TxStat           `json:"txStat,omitempty"`
	FeeStat          *FeeStat          `json:"feeStat"`
}

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

// ==================== Node API Responses ====================
type NodeInfoResponse struct {
	BeginSyncNum       int64            `json:"beginSyncNum"`
	Block              string           `json:"block"`
	SolidityBlock      string           `json:"solidityBlock"`
	CurrentConnectCount int32           `json:"currentConnectCount"`
	ActiveConnectCount  int32           `json:"activeConnectCount"`
	PassiveConnectCount int32           `json:"passiveConnectCount"`
	TotalFlow          int64            `json:"totalFlow"`
	PeerInfoList       []PeerInfo       `json:"peerInfoList"`
	ConfigNodeInfo     *ConfigNodeInfo  `json:"configNodeInfo"`
	MachineInfo        *MachineInfo     `json:"machineInfo"`
	CheatWitnessInfoMap map[string]string `json:"cheatWitnessInfoMap"`
}

type PeerInfo struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	LastBlockTime  int64  `json:"lastBlockTime"`
	Score          int    `json:"score"`
	SyncToPeer     bool   `json:"syncToPeer"`
	SyncFromPeer   bool   `json:"syncFromPeer"`
}

type ConfigNodeInfo struct {
	CodeVersion     string `json:"codeVersion"`
	VersionNum      string `json:"versionNum"`
	P2pVersion      string `json:"p2pVersion"`
	ListenPort      int    `json:"listenPort"`
	DiscoverEnable  bool   `json:"discoverEnable"`
	ActiveNodeCount int    `json:"activeNodeCount"`
	PassiveNodeCount int   `json:"passiveNodeCount"`
	SendNodeCount   int    `json:"sendNodeCount"`
	MaxConnectCount int    `json:"maxConnectCount"`
	SameIpMaxConnectCount int `json:"sameIpMaxConnectCount"`
	BackupListenPort int   `json:"backupListenPort"`
	BackupMemberSize int    `json:"backupMemberSize"`
	BackupPriority   int    `json:"backupPriority"`
	DbVersion       string `json:"dbVersion"`
	MinParticipationRate int `json:"minParticipationRate"`
	SupportConstant  bool   `json:"supportConstant"`
	MinTimeRatio     float64 `json:"minTimeRatio"`
	MaxTimeRatio     float64 `json:"maxTimeRatio"`
	AllowCreationOfContracts bool `json:"allowCreationOfContracts"`
	AllowAdaptiveEnergy      bool `json:"allowAdaptiveEnergy"`
}

type MachineInfo struct {
	ThreadCount     int32   `json:"threadCount"`
	DeadLockThreadCount int32 `json:"deadLockThreadCount"`
	CpuCount        int32   `json:"cpuCount"`
	TotalMemory     int64   `json:"totalMemory"`
	FreeMemory      int64   `json:"freeMemory"`
	CpuRate         float64 `json:"cpuRate"`
	JavaCpuRate     float64 `json:"javaCpuRate"`
	ProcessCpuRate  float64 `json:"processCpuRate"`
	MemoryDescInfoList []MemoryDescInfo `json:"memoryDescInfoList"`
}

type MemoryDescInfo struct {
	Name string `json:"name"`
	Init int64  `json:"init"`
	Used int64  `json:"used"`
	Max  int64  `json:"max"`
	Committed int64 `json:"committed"`
}

type NodeListResponse struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Address *NodeAddress `json:"address"`
}

type NodeAddress struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

// ==================== Asset/Token API Responses ====================
type AssetIssueResponse struct {
	ID                       string         `json:"id"`
	OwnerAddress             string         `json:"owner_address"`
	Name                     string         `json:"name"`
	Abbr                     string         `json:"abbr"`
	TotalSupply              int64          `json:"total_supply"`
	FrozenSupply             []FrozenSupply `json:"frozen_supply,omitempty"`
	LindNum                  int32          `json:"lind_num"`
	Num                      int32          `json:"num"`
	Precision                int32          `json:"precision,omitempty"`
	StartTime                int64          `json:"start_time"`
	EndTime                  int64          `json:"end_time"`
	VoteScore                int32          `json:"vote_score"`
	Description              string         `json:"description"`
	URL                      string         `json:"url"`
	FreeAssetNetLimit        int64          `json:"free_asset_net_limit,omitempty"`
	PublicFreeAssetNetLimit  int64          `json:"public_free_asset_net_limit,omitempty"`
	PublicFreeAssetNetUsage  int64          `json:"public_free_asset_net_usage,omitempty"`
	PublicLatestFreeNetTime  int64          `json:"public_latest_free_net_time,omitempty"`
}

type FrozenSupply struct {
	FrozenAmount int64 `json:"frozen_amount"`
	FrozenDays   int64 `json:"frozen_days"`
}

type AssetIssueListResponse struct {
	AssetIssue []AssetIssueResponse `json:"assetIssue"`
}

// ==================== LRC20 Token API Responses ====================
// type LRC20TokenInfo struct {
// 	Contract    string `json:"contract"`
// 	Name        string `json:"name"`
// 	Symbol      string `json:"symbol"`
// 	Decimals    int32  `json:"decimals"`
// 	TotalSupply string `json:"total_supply"`
// 	Owner       string `json:"owner"`
// 	IssueTime   int64  `json:"issue_time"`
// 	Holders     int64  `json:"holders,omitempty"`
// 	Transfers   int64  `json:"transfers,omitempty"`
// }

type LRC20TokenListResponse struct {
	Tokens []LRC20TokenInfo `json:"tokens"`
	Total  int64            `json:"total"`
}

type TokenHolderResponse struct {
	Address    string  `json:"address"`
	Balance    string  `json:"balance"`
	Percentage float64 `json:"percentage"`
	Rank       int64   `json:"rank,omitempty"`
}

type TokenHoldersResponse struct {
	Holders []TokenHolderResponse `json:"holders"`
	Total   int64                  `json:"total"`
}

type TokenTransferResponse struct {
	TransactionID string `json:"transaction_id"`
	BlockNumber   int64  `json:"block_number"`
	BlockTimestamp int64 `json:"block_timestamp"`
	From          string `json:"from"`
	To            string `json:"to"`
	Value         string `json:"value"`
	TokenAddress  string `json:"token_address"`
	TokenSymbol   string `json:"token_symbol"`
	TokenDecimals int32  `json:"token_decimals"`
}

type TokenTransfersResponse struct {
	Transfers []TokenTransferResponse `json:"transfers"`
	Total     int64                   `json:"total"`
}

// ==================== Account Resource API Responses ====================
type AccountResourceResponse struct {
	FreeNetUsed     int64            `json:"freeNetUsed"`
	FreeNetLimit    int64            `json:"freeNetLimit"`
	NetUsed         int64            `json:"NetUsed"`
	NetLimit        int64            `json:"NetLimit"`
	TotalNetLimit   int64            `json:"TotalNetLimit"`
	TotalNetWeight  int64            `json:"TotalNetWeight"`
	TotalLindaPowerWeight int64       `json:"totalLindaPowerWeight"`
	LindaPowerLimit  int64            `json:"lindaPowerLimit"`
	LindaPowerUsed   int64            `json:"lindaPowerUsed"`
	EnergyUsed      int64            `json:"EnergyUsed"`
	EnergyLimit     int64            `json:"EnergyLimit"`
	TotalEnergyLimit int64           `json:"TotalEnergyLimit"`
	TotalEnergyWeight int64          `json:"TotalEnergyWeight"`
	AssetNetUsed    map[string]int64 `json:"assetNetUsed,omitempty"`
	AssetNetLimit   map[string]int64 `json:"assetNetLimit,omitempty"`
}

// ==================== Delegated Resource API Responses ====================
type DelegatedResourceResponse struct {
	From                     string `json:"from"`
	To                       string `json:"to"`
	FrozenBalanceForBandwidth int64  `json:"frozen_balance_for_bandwidth,omitempty"`
	FrozenBalanceForEnergy    int64  `json:"frozen_balance_for_energy,omitempty"`
	ExpireTimeForBandwidth    int64  `json:"expire_time_for_bandwidth,omitempty"`
	ExpireTimeForEnergy       int64  `json:"expire_time_for_energy,omitempty"`
}

type DelegatedResourceListResponse struct {
	DelegatedResource []DelegatedResourceResponse `json:"delegatedResource"`
}

type DelegatedResourceAccountIndexResponse struct {
	Account      string   `json:"account"`
	FromAccounts []string `json:"fromAccounts,omitempty"`
	ToAccounts   []string `json:"toAccounts,omitempty"`
}

// ==================== Witness API Responses ====================
type WitnessResponse struct {
	Address        string `json:"address"`
	VoteCount      int64  `json:"voteCount"`
	URL            string `json:"url"`
	TotalProduced  int64  `json:"totalProduced"`
	TotalMissed    int64  `json:"totalMissed"`
	LatestBlockNum int64  `json:"latestBlockNum"`
	LatestSlotNum  int64  `json:"latestSlotNum"`
	IsJobs         bool   `json:"isJobs"`
}

type WitnessListResponse struct {
	Witnesses []WitnessResponse `json:"witnesses"`
}

// ==================== Proposal API Responses ====================
type ProposalResponse struct {
	ProposalID      int64             `json:"proposal_id"`
	ProposerAddress string            `json:"proposer_address"`
	Parameters      map[int64]int64   `json:"parameters"`
	ExpirationTime  int64             `json:"expiration_time"`
	CreateTime      int64             `json:"create_time"`
	Approvals       []string          `json:"approvals"`
	State           string            `json:"state"`
}

type ProposalListResponse struct {
	Proposals []ProposalResponse `json:"proposals"`
}

// ==================== Exchange API Responses ====================
type ExchangeResponse struct {
	ExchangeID        int64  `json:"exchange_id"`
	CreatorAddress    string `json:"creator_address"`
	CreateTime        int64  `json:"create_time"`
	FirstTokenID      string `json:"first_token_id"`
	FirstTokenBalance int64  `json:"first_token_balance"`
	SecondTokenID     string `json:"second_token_id"`
	SecondTokenBalance int64 `json:"second_token_balance"`
}

type ExchangeListResponse struct {
	Exchanges []ExchangeResponse `json:"exchanges"`
}

// ==================== Market API Responses ====================
type MarketOrderResponse struct {
	OrderID        string `json:"order_id"`
	OwnerAddress   string `json:"owner_address"`
	CreateTime     int64  `json:"create_time"`
	SellTokenID    string `json:"sell_token_id"`
	SellTokenValue int64  `json:"sell_token_value"`
	BuyTokenID     string `json:"buy_token_id"`
	BuyTokenValue  int64  `json:"buy_token_value"`
	OrderStatus    string `json:"order_status"`
}

type MarketOrderListResponse struct {
	Orders []MarketOrderResponse `json:"orders"`
	Total  int64                 `json:"total"`
}

type MarketPriceResponse struct {
	SellTokenID string `json:"sell_token_id"`
	BuyTokenID  string `json:"buy_token_id"`
	Price       string `json:"price"`
}

type MarketPriceListResponse struct {
	Prices []MarketPriceResponse `json:"prices"`
}

type MarketPairResponse struct {
	SellTokenID string `json:"sell_token_id"`
	BuyTokenID  string `json:"buy_token_id"`
}

type MarketPairListResponse struct {
	Pairs []MarketPairResponse `json:"pairs"`
}

// ==================== Smart Contract API Responses ====================
type SmartContractResponse struct {
	OriginAddress               string          `json:"origin_address"`
	ContractAddress             string          `json:"contract_address"`
	ABI                         json.RawMessage `json:"abi"`
	Bytecode                    string          `json:"bytecode"`
	CallValue                   int64           `json:"call_value,omitempty"`
	ConsumeUserResourcePercent  int32           `json:"consume_user_resource_percent"`
	Name                        string          `json:"name"`
	OriginEnergyLimit           int64           `json:"origin_energy_limit"`
	CodeHash                    string          `json:"code_hash"`
}

type ContractInfoResponse struct {
	RuntimeCode   string                `json:"runtimecode"`
	SmartContract *SmartContractResponse `json:"smart_contract"`
	ContractState *ContractState        `json:"contract_state,omitempty"`
}

type ContractState struct {
	EnergyUsage  int64 `json:"energy_usage"`
	EnergyFactor int64 `json:"energy_factor"`
	UpdateCycle  int64 `json:"update_cycle"`
}

type EstimateEnergyResponse struct {
	Result         *ResultResponse `json:"result"`
	EnergyRequired int64           `json:"energy_required"`
}

type ResultResponse struct {
	Result  bool   `json:"result"`
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type TransactionExtentionResponse struct {
	Transaction *TransactionResponse `json:"transaction"`
	TxID        string               `json:"txid"`
	Result      *ResultResponse      `json:"result"`
}

// ==================== Chain Parameters API Responses ====================
type ChainParameter struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type ChainParametersResponse struct {
	ChainParameter []ChainParameter `json:"chainParameter"`
}

// ==================== Statistics API Responses ====================

// DailyStat represents daily statistics data
type DailyStat struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// EnergyStatisticResponse represents energy statistics response
type EnergyStatisticResponse struct {
	Address string     `json:"address"`
	Total   int64      `json:"total"`
	Daily   []DailyStat `json:"daily"`
}

// TriggerStatisticResponse represents trigger statistics response
type TriggerStatisticResponse struct {
	Contract string     `json:"contract"`
	Count    int64      `json:"count"`
	Daily    []DailyStat `json:"daily"`
}

// CallerAddressStatisticResponse represents caller statistics response
type CallerAddressStatisticResponse struct {
	Contract string       `json:"contract"`
	Callers  []CallerStat `json:"callers"`
}

// CallerStat represents a caller statistic
type CallerStat struct {
	Address string `json:"address"`
	Count   int64  `json:"count"`
}

// EnergyDailyStatisticResponse represents daily energy statistics
type EnergyDailyStatisticResponse struct {
	Daily []EnergyDailyStat `json:"daily"`
}

// EnergyDailyStat represents a daily energy statistic
type EnergyDailyStat struct {
	Date        string `json:"date"`
	EnergyUsage int64  `json:"energy_usage"`
	EnergyFee   int64  `json:"energy_fee"`
}

// FreezeResourceResponse represents freeze resource response
type FreezeResourceResponse struct {
	Frozen     []FrozenRecord    `json:"frozen"`
	Delegated  []DelegatedRecord `json:"delegated"`
}

// FrozenRecord represents a frozen record
type FrozenRecord struct {
	Amount   int64  `json:"amount"`
	ExpireAt int64  `json:"expire_at"`
	Type     string `json:"type"`
}

// DelegatedRecord represents a delegated record
type DelegatedRecord struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Amount   int64  `json:"amount"`
	ExpireAt int64  `json:"expire_at"`
	Type     string `json:"type"`
}

// TurnoverResponse represents turnover response
type TurnoverResponse struct {
	Total int64           `json:"total"`
	Daily []DailyTurnover `json:"daily"`
}

// DailyTurnover represents daily turnover
type DailyTurnover struct {
	Date     string `json:"date"`
	Turnover int64  `json:"turnover"`
}

// LindHolderResponse represents LIND holders response
type LindHolderResponse struct {
	Holders []LindHolder `json:"holders"`
	Total   int64        `json:"total"`
}

// LindHolder represents a LIND holder
type LindHolder struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
	Rank    int    `json:"rank"`
}

// TokenPriceResponse represents token price response
type TokenPriceResponse struct {
	Price     float64 `json:"price"`
	Change24h float64 `json:"change_24h"`
	Volume24h float64 `json:"volume_24h"`
	MarketCap float64 `json:"market_cap"`
}

// ================================================================

// TokenPositionResponse represents token position distribution
type TokenPositionResponse struct {
	Addresses []TokenPosition `json:"addresses"`
}

type TokenPosition struct {
	Address    string  `json:"address"`
	Balance    string  `json:"balance"`
	Percentage float64 `json:"percentage"`
	Rank       int     `json:"rank"`
}

// WinkFundResponse represents WINK fund information
type WinkFundResponse struct {
	Total       int64 `json:"total"`
	Burned      int64 `json:"burned"`
	Circulating int64 `json:"circulating"`
}

// JSTFundResponse represents JST fund information
type JSTFundResponse struct {
	Total       int64 `json:"total"`
	Burned      int64 `json:"burned"`
	Circulating int64 `json:"circulating"`
}

// GraphicResponse represents graphic data
type GraphicResponse struct {
	Labels []string  `json:"labels"`
	Data   []float64 `json:"data"`
}

// ==================== JSON-RPC API Responses ====================
type JsonRpcResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JsonRpcError   `json:"error,omitempty"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JsonRpcBlockResponse struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []interface{} `json:"transactions"`
	Uncles           []string      `json:"uncles"`
}

type JsonRpcTransactionResponse struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type JsonRpcTransactionReceiptResponse struct {
	TransactionHash   string        `json:"transactionHash"`
	TransactionIndex  string        `json:"transactionIndex"`
	BlockHash         string        `json:"blockHash"`
	BlockNumber       string        `json:"blockNumber"`
	From              string        `json:"from"`
	To                string        `json:"to"`
	CumulativeGasUsed string        `json:"cumulativeGasUsed"`
	GasUsed           string        `json:"gasUsed"`
	ContractAddress   string        `json:"contractAddress"`
	Logs              []JsonRpcLog  `json:"logs"`
	LogsBloom         string        `json:"logsBloom"`
	Status            string        `json:"status"`
	EffectiveGasPrice string        `json:"effectiveGasPrice"`
}

type JsonRpcLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

// ==================== Event Query API Responses ====================
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

type EventTransactionsResponse struct {
	Data  []EventTransactionResponse `json:"data"`
	Meta  *PaginationMeta            `json:"meta"`
	Success bool                     `json:"success"`
}

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

type TransferListResponse struct {
	Data  []TransferResponse `json:"data"`
	Meta  *PaginationMeta    `json:"meta"`
	Success bool             `json:"success"`
}

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

type EventListResponse struct {
	Data  []EventResponse   `json:"data"`
	Meta  *PaginationMeta   `json:"meta"`
	Success bool            `json:"success"`
}

type EventBlockResponse struct {
	Hash              string `json:"hash"`
	Number            int64  `json:"number"`
	Timestamp         int64  `json:"timestamp"`
	ParentHash        string `json:"parentHash"`
	WitnessAddress    string `json:"witnessAddress"`
	TransactionCount  int    `json:"transactionCount"`
}

type BlockListV2Response struct {
	Data  []EventBlockResponse `json:"data"`
	Meta  *PaginationMeta      `json:"meta"`
	Success bool               `json:"success"`
}

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

type ContractLogsResponse struct {
	Data  []ContractLogResponse `json:"data"`
	Meta  *PaginationMeta       `json:"meta"`
	Success bool                `json:"success"`
}

type ContractWithAbiResponse struct {
	ContractAddress string          `json:"contractAddress"`
	ABI             json.RawMessage `json:"abi"`
	Logs            []interface{}   `json:"logs"`
}

// ==================== Lindascan Custom API Responses ====================
type HomepageBundleResponse struct {
	TotalBlocks        int64                 `json:"totalBlocks"`
	TotalTransactions  int64                 `json:"totalTransactions"`
	TotalAccounts      int64                 `json:"totalAccounts"`
	TotalContracts     int64                 `json:"totalContracts"`
	TotalTokens        int64                 `json:"totalTokens"`
	PriceUSD           int64                 `json:"priceUSD"`
	MarketCap          int64                 `json:"marketCap"`
	Volume24h          int64                 `json:"volume24h"`
	RecentBlocks       []BlockResponse       `json:"recentBlocks"`
	RecentTransactions []TransactionResponse `json:"recentTransactions"`
}

type NodeMapResponse struct {
	Nodes []NodeInfo `json:"nodes"`
}

type NodeInfo struct {
	IP        string  `json:"ip"`
	Host      string  `json:"host"`
	Port      int32   `json:"port"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	NodeType  string  `json:"nodeType"`
}

type Top10Response struct {
	Witnesses []WitnessResponse `json:"witnesses"`
	Accounts  []AccountResponse `json:"accounts"`
	Tokens    []LRC20TokenInfo  `json:"tokens"`
}

type TagResponse struct {
	ID          int32  `json:"id,omitempty"`
	Address     string `json:"address"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	CreatedAt   int64  `json:"createdAt"`
	Votes       int32  `json:"votes"`
}

type TagListResponse struct {
	Tags  []TagResponse `json:"tags"`
	Total int64         `json:"total"`
}

type TagInsertResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      int32  `json:"id,omitempty"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// ==================== Pagination Meta ====================
type PaginationMeta struct {
	At          int64  `json:"at"`
	PageSize    int    `json:"page_size"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Links       *Links `json:"links,omitempty"`
}

type Links struct {
	Next string `json:"next,omitempty"`
}