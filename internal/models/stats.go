package models

// Statistic represents a stored statistic
type Statistic struct {
	ID        uint            `gorm:"primarykey" json:"-"`
	Type      string          `gorm:"index;type:varchar(50)" json:"type"`
	Value     json.RawMessage `gorm:"type:jsonb" json:"value"`
	Timestamp int64           `gorm:"index" json:"timestamp"`
}

// HomepageBundleResponse represents the homepage bundle data
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

// Top10Request represents the request for top 10 data
type Top10Request struct {
	Type string `json:"type" form:"type"`
}

// Top10Response represents the response for top 10 data
type Top10Response struct {
	Witnesses []WitnessResponse `json:"witnesses"`
	Accounts  []AccountResponse `json:"accounts"`
	Tokens    []LRC20TokenInfo  `json:"tokens"`
}

// EnergyStatisticRequest represents request for energy statistics
type EnergyStatisticRequest struct {
	Address string `json:"address" form:"address"`
	From    int64  `json:"from" form:"from"`
	To      int64  `json:"to" form:"to"`
}

// EnergyStatisticResponse represents energy statistics response
type EnergyStatisticResponse struct {
	Address string     `json:"address"`
	Total   int64      `json:"total"`
	Daily   []DailyStat `json:"daily"`
}

// DailyStat represents daily statistics
type DailyStat struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// TriggerStatisticRequest represents request for trigger statistics
type TriggerStatisticRequest struct {
	Contract string `json:"contract" form:"contract"`
	From     int64  `json:"from" form:"from"`
	To       int64  `json:"to" form:"to"`
}

// TriggerStatisticResponse represents trigger statistics response
type TriggerStatisticResponse struct {
	Contract string     `json:"contract"`
	Count    int64      `json:"count"`
	Daily    []DailyStat `json:"daily"`
}

// CallerAddressStatisticRequest represents request for caller statistics
type CallerAddressStatisticRequest struct {
	Contract string `json:"contract" form:"contract"`
	From     int64  `json:"from" form:"from"`
	To       int64  `json:"to" form:"to"`
}

// CallerAddressStatisticResponse represents caller statistics response
type CallerAddressStatisticResponse struct {
	Contract string       `json:"contract"`
	Callers  []CallerStat `json:"callers"`
}

// CallerStat represents caller statistics
type CallerStat struct {
	Address string `json:"address"`
	Count   int64  `json:"count"`
}

// EnergyDailyStatisticRequest represents request for daily energy statistics
type EnergyDailyStatisticRequest struct {
	From int64 `json:"from" form:"from"`
	To   int64 `json:"to" form:"to"`
}

// EnergyDailyStatisticResponse represents daily energy statistics response
type EnergyDailyStatisticResponse struct {
	Daily []EnergyDailyStat `json:"daily"`
}

// EnergyDailyStat represents daily energy statistics
type EnergyDailyStat struct {
	Date        string `json:"date"`
	EnergyUsage int64  `json:"energy_usage"`
	EnergyFee   int64  `json:"energy_fee"`
}

// TriggerAmountStatisticRequest represents request for trigger amount statistics
type TriggerAmountStatisticRequest struct {
	Contract string `json:"contract" form:"contract"`
	From     int64  `json:"from" form:"from"`
	To       int64  `json:"to" form:"to"`
}

// TriggerAmountStatisticResponse represents trigger amount statistics response
type TriggerAmountStatisticResponse struct {
	Contract string       `json:"contract"`
	Triggers []TriggerStat `json:"triggers"`
}

// TriggerStat represents trigger statistics
type TriggerStat struct {
	TxID    string `json:"tx_id"`
	Amount  int64  `json:"amount"`
	Address string `json:"address"`
}

// FreezeResourceRequest represents request for freeze resource
type FreezeResourceRequest struct {
	Address string `json:"address" form:"address"`
	Type    string `json:"type" form:"type"`
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

// TurnoverRequest represents request for turnover statistics
type TurnoverRequest struct {
	From int64 `json:"from" form:"from"`
	To   int64 `json:"to" form:"to"`
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

// LindHolderRequest represents request for LIND holders
type LindHolderRequest struct {
	PaginationRequest
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

// OneContractEnergyStatisticRequest represents request for one contract energy statistics
type OneContractEnergyStatisticRequest struct {
	Contract string `json:"contract" form:"contract"`
	From     int64  `json:"from" form:"from"`
	To       int64  `json:"to" form:"to"`
}

// OneContractEnergyStatisticResponse represents one contract energy statistics response
type OneContractEnergyStatisticResponse struct {
	Contract string     `json:"contract"`
	Total    int64      `json:"total"`
	Daily    []DailyStat `json:"daily"`
}

// OneContractTriggerStatisticRequest represents request for one contract trigger statistics
type OneContractTriggerStatisticRequest struct {
	Contract string `json:"contract" form:"contract"`
	From     int64  `json:"from" form:"from"`
	To       int64  `json:"to" form:"to"`
}

// OneContractTriggerStatisticResponse represents one contract trigger statistics response
type OneContractTriggerStatisticResponse struct {
	Contract string     `json:"contract"`
	Count    int64      `json:"count"`
	Daily    []DailyStat `json:"daily"`
}

// OneContractCallerStatisticRequest represents request for one contract caller statistics
type OneContractCallerStatisticRequest struct {
	Contract string `json:"contract" form:"contract"`
	From     int64  `json:"from" form:"from"`
	To       int64  `json:"to" form:"to"`
}

// OneContractCallerStatisticResponse represents one contract caller statistics response
type OneContractCallerStatisticResponse struct {
	Contract string       `json:"contract"`
	Callers  []CallerStat `json:"callers"`
}

// OneContractCallersRequest represents request for one contract callers
type OneContractCallersRequest struct {
	Contract string `json:"contract" form:"contract"`
	PaginationRequest
}

// OneContractCallersResponse represents one contract callers response
type OneContractCallersResponse struct {
	Callers []CallerDetail `json:"callers"`
	Total   int64          `json:"total"`
}

// CallerDetail represents caller details
type CallerDetail struct {
	Address string `json:"address"`
	Count   int64  `json:"count"`
	FirstTx string `json:"first_tx"`
	LastTx  string `json:"last_tx"`
}

// NodeOverviewRequest represents request for node overview
type NodeOverviewRequest struct {
	Address string `json:"address" binding:"required"`
}

// NodeOverviewResponse represents node overview response
type NodeOverviewResponse struct {
	Uptime      int64   `json:"uptime"`
	BlockHeight int64   `json:"block_height"`
	Peers       int     `json:"peers"`
	Version     string  `json:"version"`
}

// NodeInfoUploadRequest represents request for node info upload
type NodeInfoUploadRequest struct {
	Address     string `json:"address" binding:"required"`
	Version     string `json:"version"`
	Location    string `json:"location"`
	BlockHeight int64  `json:"block_height"`
	Peers       int    `json:"peers"`
}

// NodeInfoUploadResponse represents node info upload response
type NodeInfoUploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}