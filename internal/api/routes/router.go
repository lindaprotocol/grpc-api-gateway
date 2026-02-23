package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/api/handlers"
	"github.com/lindaprotocol/grpc-api-gateway/internal/api/middleware"
	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/cache"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
)

type Router struct {
	engine           *gin.Engine
	config           *config.Config
	blockchainClient *blockchain.Client
	authService      *auth.Service
	cacheClient      *cache.RedisClient
	
	// Handlers
	accountHandler     *handlers.AccountHandler
	blockHandler       *handlers.BlockHandler
	transactionHandler *handlers.TransactionHandler
	tokenHandler       *handlers.TokenHandler
	contractHandler    *handlers.ContractHandler
	nodeHandler        *handlers.NodeHandler
	statsHandler       *handlers.StatsHandler
	searchHandler      *handlers.SearchHandler
	eventHandler       *handlers.EventHandler
}

func NewRouter(
	cfg *config.Config,
	client *blockchain.Client,
	authSvc *auth.Service,
	cacheClient *cache.RedisClient,
	accountRepo *repository.AccountRepository,
	blockRepo *repository.BlockRepository,
	txRepo *repository.TransactionRepository,
	tokenRepo *repository.TokenRepository,
	eventRepo *repository.EventRepository,
	tagRepo *repository.TagRepository,
	statsRepo *repository.StatsRepository,
) *Router {
	router := &Router{
		engine:           gin.New(),
		config:           cfg,
		blockchainClient: client,
		authService:      authSvc,
		cacheClient:      cacheClient,
	}

	// Initialize handlers
	router.accountHandler = handlers.NewAccountHandler(client, accountRepo, tagRepo)
	router.blockHandler = handlers.NewBlockHandler(client, blockRepo)
	router.transactionHandler = handlers.NewTransactionHandler(client, txRepo)
	router.tokenHandler = handlers.NewTokenHandler(client, tokenRepo)
	router.contractHandler = handlers.NewContractHandler(client)
	router.nodeHandler = handlers.NewNodeHandler(client)
	router.statsHandler = handlers.NewStatsHandler(client, statsRepo)
	router.searchHandler = handlers.NewSearchHandler(client, accountRepo, blockRepo, txRepo, tokenRepo)
	router.eventHandler = handlers.NewEventHandler(client, eventRepo)

	router.setupMiddleware()
	router.setupRoutes()

	return router
}

func (r *Router) setupMiddleware() {
	// Recovery middleware
	r.engine.Use(gin.Recovery())

	// Logger middleware
	r.engine.Use(middleware.LoggerMiddleware(r.config.Logging))

	// CORS middleware
	r.engine.Use(middleware.CorsMiddleware(r.config.CORS.AllowedOrigins))

	// Rate limiting middleware
	r.engine.Use(middleware.RateLimitMiddleware(r.cacheClient, r.config.RateLimit))

	// Auth middleware
	r.engine.Use(middleware.AuthMiddleware(r.authService, r.config.Auth))

	// Allowlist middleware
	r.engine.Use(middleware.AllowlistMiddleware(r.authService))

	// Response interceptor
	r.engine.Use(middleware.ResponseInterceptor())
}

func (r *Router) setupRoutes() {
	// Health check
	r.engine.GET("/health", r.handleHealth)

	// Metrics
	r.engine.GET("/metrics", r.handleMetrics)

	// API v1 routes (Event Query Service style)
	v1 := r.engine.Group("/v1")
	{
		// Transactions
		v1.GET("/transactions", r.transactionHandler.GetTransactionsEvent)
		v1.GET("/transactions/:hash", r.transactionHandler.GetTransactionByHashEvent)
		
		// Transfers
		v1.GET("/transfers", r.tokenHandler.GetTransfersEvent)
		v1.GET("/transfers/:hash", r.tokenHandler.GetTransferByHashEvent)
		
		// Events
		v1.GET("/events", r.eventHandler.GetEvents)
		v1.GET("/events/transaction/:transaction_id", r.eventHandler.GetEventsByTransactionId)
		v1.GET("/events/:contract_address", r.eventHandler.GetEventsByContractAddress)
		v1.GET("/events/contract/:contract_address/:event_name", r.eventHandler.GetEventsByContractAndEventName)
		v1.GET("/events/contract/:contract_address/:event_name/:block_number", r.eventHandler.GetEventsByContractEventAndBlock)
		v1.GET("/events/timestamp", r.eventHandler.GetEventsByTimestamp)
		v1.GET("/events/confirmed", r.eventHandler.GetConfirmedEvents)
		
		// Blocks
		v1.GET("/blocks/:hash", r.blockHandler.GetBlockByHashEvent)
		v1.GET("/blocks", r.blockHandler.GetBlocksEvent)
		v1.GET("/blocks/latestSolidifiedBlockNumber", r.blockHandler.GetLatestSolidifiedBlockNumber)
		v1.GET("/blocks/:blockNum/stats", r.blockHandler.GetBlockStats)
		
		// Contract logs
		v1.GET("/contractlogs", r.contractHandler.GetContractLogs)
		v1.GET("/contractlogs/transaction/:transaction_id", r.contractHandler.GetContractLogsByTransactionId)
		v1.GET("/contractlogs/contract/:contract_address", r.contractHandler.GetContractLogsByContractAddress)
		v1.POST("/contract/transaction/:transaction_id", r.contractHandler.GetContractWithAbi)
		v1.POST("/contract/contractAddress/:contract_address", r.contractHandler.GetContractByAddressWithAbi)
		
		// Accounts (v1 style)
		v1.GET("/accounts/:address", r.accountHandler.GetAccountV1)
		v1.GET("/accounts/:address/transactions", r.transactionHandler.GetAccountTransactionsV1)
		v1.GET("/accounts/:address/transactions/lrc20", r.tokenHandler.GetAccountLRC20TransactionsV1)
		v1.GET("/accounts/:address/internal-transactions", r.transactionHandler.GetAccountInternalTransactionsV1)
		v1.GET("/accounts/:address/lrc20/balance", r.tokenHandler.GetAccountLRC20BalanceV1)
		
		// Assets (v1 style)
		v1.GET("/assets", r.tokenHandler.GetAssetsV1)
		v1.GET("/assets/:name/list", r.tokenHandler.GetAssetsByNameV1)
		v1.GET("/assets/:identifier", r.tokenHandler.GetAssetsByIdentifierV1)
		
		// Contracts (v1 style)
		v1.GET("/contracts/:contractAddress/transactions", r.transactionHandler.GetContractTransactionsV1)
		v1.GET("/accounts/:contractAddress/internal-transactions", r.transactionHandler.GetContractInternalTransactionsV1)
		v1.GET("/contracts/:contractAddress/tokens", r.tokenHandler.GetContractTokensV1)
	}

	// Lindascan custom API routes
	api := r.engine.Group("/api")
	{
		// System
		api.GET("/system/homepage-bundle", r.statsHandler.GetHomepageBundle)
		api.POST("/system/proxy", r.handleProxy)
		
		// Node map
		api.GET("/nodemap", r.nodeHandler.GetNodeMap)
		
		// Top 10
		api.GET("/top10", r.statsHandler.GetTop10)
		
		// Token APIs
		api.GET("/token", r.tokenHandler.GetTokens)
		api.GET("/token/:id", r.tokenHandler.GetTokenById)
		api.GET("/token/price", r.tokenHandler.GetTokenPrice)
		api.GET("/tokens/overview", r.tokenHandler.GetTokensOverview)
		api.GET("/token_lrc20", r.tokenHandler.GetLRC20Tokens)
		api.GET("/token_lrc20/:contract", r.tokenHandler.GetLRC20TokenByContract)
		api.GET("/token_lrc20/holders", r.tokenHandler.GetTokenHolders)
		api.GET("/token_lrc20/transfers", r.tokenHandler.GetTokenTransfers)
		api.GET("/tokens/position-distribution", r.tokenHandler.GetTokenPositionDistribution)
		api.GET("/tokens/participateassetissue", r.tokenHandler.GetParticipateAssetIssue)
		
		// Special tokens
		api.GET("/wink/fund", r.tokenHandler.GetWinkFund)
		api.GET("/wink/graphic", r.tokenHandler.GetWinkGraphic)
		api.GET("/jst/fund", r.tokenHandler.GetJSTFund)
		api.GET("/jst/graphic", r.tokenHandler.GetJSTGraphic)
		api.GET("/bittorrent/graphic", r.tokenHandler.GetBitTorrentGraphic)
		
		// Asset transfers
		api.GET("/asset/transfer", r.tokenHandler.GetAssetTransfer)
		
		// Account APIs
		api.GET("/account/list", r.accountHandler.GetAccountList)
		api.GET("/account/resource", r.accountHandler.GetAccountResourceInfo)
		api.GET("/account-proposal", r.accountHandler.GetAccountProposals)
		
		// Statistics
		api.GET("/stats/overview", r.statsHandler.GetOverview)
		api.GET("/energystatistic", r.statsHandler.GetEnergyStatistic)
		api.GET("/triggerstatistic", r.statsHandler.GetTriggerStatistic)
		api.GET("/calleraddressstatistic", r.statsHandler.GetCallerAddressStatistic)
		api.GET("/energydailystatistic", r.statsHandler.GetEnergyDailyStatistic)
		api.GET("/triggeramountstatistic", r.statsHandler.GetTriggerAmountStatistic)
		api.GET("/freezeresource", r.statsHandler.GetFreezeResource)
		api.GET("/turnover", r.statsHandler.GetTurnover)
		api.GET("/onecontractenergystatistic", r.statsHandler.GetOneContractEnergyStatistic)
		api.GET("/onecontracttriggerstatistic", r.statsHandler.GetOneContractTriggerStatistic)
		api.GET("/onecontractcallerstatistic", r.statsHandler.GetOneContractCallerStatistic)
		api.GET("/onecontractcallers", r.statsHandler.GetOneContractCallers)
		
		// Block and transaction
		api.GET("/block", r.blockHandler.GetBlocksV2)
		api.GET("/transaction", r.transactionHandler.GetTransactionsV2)
		api.GET("/internal-transaction", r.transactionHandler.GetInternalTransactions)
		api.GET("/contracts/transaction", r.transactionHandler.GetContractTransactions)
		api.GET("/lrc10lrc20-transfer", r.tokenHandler.GetLRC10LRC20Transfers)
		api.GET("/contract_account_history", r.contractHandler.GetContractAccountHistory)
		api.GET("/contracts/smart-contract-triggers-batch", r.contractHandler.GetSmartContractTriggersBatch)
		
		// Chain parameters
		api.GET("/chainparameters", r.nodeHandler.GetChainParametersV2)
		
		// Vote
		api.GET("/vote", r.nodeHandler.GetVoteInfo)
		
		// Search
		api.GET("/search", r.searchHandler.Search)
		
		// Fund
		api.GET("/fund", r.statsHandler.GetFund)
		
		// Ledger
		api.GET("/ledger", r.statsHandler.GetLedger)
		
		// Testnet
		api.POST("/testnet/request-coins", r.handleTestnetRequest)
		
		// CSV Export
		api.GET("/export", r.handleCSVExport)
		
		// Monitor
		api.POST("/monitor", r.handleMonitor)
		
		// V2 node endpoints
		api.POST("/v2/node/overview_upload", r.nodeHandler.UploadNodeOverview)
		api.POST("/v2/node/info_upload", r.nodeHandler.UploadNodeInfo)
	}

	// External APIs (tag system, uploads)
	external := r.engine.Group("/external")
	{
		// Tags
		external.GET("/tag", r.accountHandler.GetTags)
		external.POST("/tag/insert", r.accountHandler.InsertTag)
		external.POST("/tag/update", r.accountHandler.UpdateTag)
		external.POST("/tag/delete", r.accountHandler.DeleteTag)
		external.GET("/tag/recommend", r.accountHandler.RecommendTag)
		
		// Upload
		external.POST("/upload/logo", r.handleLogoUpload)
	}

	// FullNode HTTP API routes (gRPC gateway style)
	fullnode := r.engine.Group("/wallet")
	{
		// Account
		fullnode.POST("/getaccount", r.accountHandler.GetAccount)
		fullnode.POST("/getaccountbalance", r.accountHandler.GetAccountBalance)
		fullnode.POST("/getaccountresource", r.accountHandler.GetAccountResource)
		fullnode.POST("/getaccountnet", r.accountHandler.GetAccountNet)
		fullnode.POST("/createaccount", r.accountHandler.CreateAccount)
		fullnode.POST("/updateaccount", r.accountHandler.UpdateAccount)
		fullnode.POST("/setaccountid", r.accountHandler.SetAccountId)
		fullnode.POST("/accountpermissionupdate", r.accountHandler.AccountPermissionUpdate)
		
		// Transaction
		fullnode.POST("/createtransaction", r.transactionHandler.CreateTransaction)
		fullnode.POST("/broadcasttransaction", r.transactionHandler.BroadcastTransaction)
		fullnode.POST("/broadcasthex", r.transactionHandler.BroadcastHex)
		fullnode.POST("/gettransactionbyid", r.transactionHandler.GetTransactionById)
		fullnode.POST("/gettransactioninfobyid", r.transactionHandler.GetTransactionInfoById)
		fullnode.POST("/gettransactionreceiptbyid", r.transactionHandler.GetTransactionReceiptById)
		fullnode.POST("/gettransactioncountbyblocknum", r.transactionHandler.GetTransactionCountByBlockNum)
		fullnode.POST("/gettransactioninfobyblocknum", r.transactionHandler.GetTransactionInfoByBlockNum)
		fullnode.POST("/gettransactionsign", r.transactionHandler.GetTransactionSign)
		fullnode.POST("/getsignweight", r.transactionHandler.GetTransactionSignWeight)
		fullnode.POST("/getapprovedlist", r.transactionHandler.GetTransactionApprovedList)
		
		// Block
		fullnode.POST("/getnowblock", r.blockHandler.GetNowBlock)
		fullnode.POST("/getblockbynum", r.blockHandler.GetBlockByNum)
		fullnode.POST("/getblockbyid", r.blockHandler.GetBlockById)
		fullnode.POST("/getblockbylimitnext", r.blockHandler.GetBlockByLimitNext)
		fullnode.POST("/getblockbylatestnum", r.blockHandler.GetBlockByLatestNum)
		fullnode.POST("/getblock", r.blockHandler.GetBlock)
		fullnode.POST("/getblockbalance", r.blockHandler.GetBlockBalance)
		
		// Node
		fullnode.POST("/listnodes", r.nodeHandler.ListNodes)
		fullnode.POST("/getnodeinfo", r.nodeHandler.GetNodeInfo)
		
		// Asset (LRC-10)
		fullnode.POST("/getassetissuebyaccount", r.tokenHandler.GetAssetIssueByAccount)
		fullnode.POST("/getassetissuebyid", r.tokenHandler.GetAssetIssueById)
		fullnode.POST("/getassetissuebyname", r.tokenHandler.GetAssetIssueByName)
		fullnode.GET("/getassetissuelist", r.tokenHandler.GetAssetIssueList)
		fullnode.POST("/getassetissuelistbyname", r.tokenHandler.GetAssetIssueListByName)
		fullnode.POST("/getpaginatedassetissuelist", r.tokenHandler.GetPaginatedAssetIssueList)
		fullnode.POST("/createassetissue", r.tokenHandler.CreateAssetIssue)
		fullnode.POST("/transferasset", r.tokenHandler.TransferAsset)
		fullnode.POST("/participateassetissue", r.tokenHandler.ParticipateAssetIssue)
		fullnode.POST("/unfreezeasset", r.tokenHandler.UnfreezeAsset)
		fullnode.POST("/updateasset", r.tokenHandler.UpdateAsset)
		
		// Resource (Stake 1.0)
		fullnode.POST("/freezebalance", r.contractHandler.FreezeBalance)
		fullnode.POST("/unfreezebalance", r.contractHandler.UnfreezeBalance)
		fullnode.POST("/withdrawbalance", r.contractHandler.WithdrawBalance)
		
		// Resource (Stake 2.0)
		fullnode.POST("/freezebalancev2", r.contractHandler.FreezeBalanceV2)
		fullnode.POST("/unfreezebalancev2", r.contractHandler.UnfreezeBalanceV2)
		fullnode.POST("/withdrawexpireunfreeze", r.contractHandler.WithdrawExpireUnfreeze)
		fullnode.POST("/delegateresource", r.contractHandler.DelegateResource)
		fullnode.POST("/undelegateresource", r.contractHandler.UnDelegateResource)
		fullnode.POST("/cancelallunfreezev2", r.contractHandler.CancelAllUnfreezeV2)
		fullnode.POST("/getavailableunfreezecount", r.contractHandler.GetAvailableUnfreezeCount)
		fullnode.POST("/getcanwithdrawunfreezeamount", r.contractHandler.GetCanWithdrawUnfreezeAmount)
		fullnode.POST("/getdelegatedresource", r.contractHandler.GetDelegatedResource)
		fullnode.POST("/getdelegatedresourcev2", r.contractHandler.GetDelegatedResourceV2)
		fullnode.POST("/getdelegatedresourceaccountindex", r.contractHandler.GetDelegatedResourceAccountIndex)
		fullnode.POST("/getdelegatedresourceaccountindexv2", r.contractHandler.GetDelegatedResourceAccountIndexV2)
		fullnode.POST("/getcandelegatedmaxsize", r.contractHandler.GetCanDelegatedMaxSize)
		
		// Witness & Voting
		fullnode.POST("/listwitnesses", r.nodeHandler.ListWitnesses)
		fullnode.POST("/createwitness", r.nodeHandler.CreateWitness)
		fullnode.POST("/updatewitness", r.nodeHandler.UpdateWitness)
		fullnode.POST("/votewitnessaccount", r.nodeHandler.VoteWitnessAccount)
		fullnode.POST("/getBrokerage", r.nodeHandler.GetBrokerage)
		fullnode.POST("/updateBrokerage", r.nodeHandler.UpdateBrokerage)
		fullnode.POST("/getReward", r.nodeHandler.GetReward)
		fullnode.GET("/getnextmaintenancetime", r.nodeHandler.GetNextMaintenanceTime)
		fullnode.POST("/getpaginatednowwitnesslist", r.nodeHandler.GetPaginatedNowWitnessList)
		
		// Proposal
		fullnode.GET("/listproposals", r.nodeHandler.ListProposals)
		fullnode.POST("/getproposalbyid", r.nodeHandler.GetProposalById)
		fullnode.POST("/proposalcreate", r.nodeHandler.ProposalCreate)
		fullnode.POST("/proposalapprove", r.nodeHandler.ProposalApprove)
		fullnode.POST("/proposaldelete", r.nodeHandler.ProposalDelete)
		fullnode.POST("/getpaginatedproposallist", r.nodeHandler.GetPaginatedProposalList)
		
		// Exchange
		fullnode.POST("/exchangecreate", r.contractHandler.ExchangeCreate)
		fullnode.POST("/exchangeinject", r.contractHandler.ExchangeInject)
		fullnode.POST("/exchangewithdraw", r.contractHandler.ExchangeWithdraw)
		fullnode.POST("/exchangetransaction", r.contractHandler.ExchangeTransaction)
		fullnode.POST("/getexchangebyid", r.contractHandler.GetExchangeById)
		fullnode.GET("/listexchanges", r.contractHandler.ListExchanges)
		fullnode.POST("/getpaginatedexchangelist", r.contractHandler.GetPaginatedExchangeList)
		
		// Market (DEX)
		fullnode.POST("/marketsellasset", r.contractHandler.MarketSellAsset)
		fullnode.POST("/marketcancelorder", r.contractHandler.MarketCancelOrder)
		fullnode.POST("/getmarketorderbyaccount", r.contractHandler.GetMarketOrderByAccount)
		fullnode.POST("/getmarketorderbyid", r.contractHandler.GetMarketOrderById)
		fullnode.POST("/getmarketpricebypair", r.contractHandler.GetMarketPriceByPair)
		fullnode.POST("/getmarketorderlistbypair", r.contractHandler.GetMarketOrderListByPair)
		fullnode.GET("/getmarketpairlist", r.contractHandler.GetMarketPairList)
		
		// Smart Contract
		fullnode.POST("/deploycontract", r.contractHandler.DeployContract)
		fullnode.POST("/triggersmartcontract", r.contractHandler.TriggerSmartContract)
		fullnode.POST("/triggerconstantcontract", r.contractHandler.TriggerConstantContract)
		fullnode.POST("/estimateenergy", r.contractHandler.EstimateEnergy)
		fullnode.POST("/getcontract", r.contractHandler.GetContract)
		fullnode.POST("/getcontractinfo", r.contractHandler.GetContractInfo)
		fullnode.POST("/updatesetting", r.contractHandler.UpdateSetting)
		fullnode.POST("/updateenergylimit", r.contractHandler.UpdateEnergyLimit)
		fullnode.POST("/clearabi", r.contractHandler.ClearAbi)
		
		// Shielded Transaction
		fullnode.POST("/getnewshieldedaddress", r.contractHandler.GetNewShieldedAddress)
		fullnode.POST("/getspendingkey", r.contractHandler.GetSpendingKey)
		fullnode.POST("/getexpandedspendingkey", r.contractHandler.GetExpandedSpendingKey)
		fullnode.POST("/getakfromask", r.contractHandler.GetAkFromAsk)
		fullnode.POST("/getnkfromnsk", r.contractHandler.GetNkFromNsk)
		fullnode.POST("/getincomingviewingkey", r.contractHandler.GetIncomingViewingKey)
		fullnode.POST("/getdiversifier", r.contractHandler.GetDiversifier)
		fullnode.POST("/getzenpaymentaddress", r.contractHandler.GetZenPaymentAddress)
		fullnode.POST("/createshieldedtransaction", r.contractHandler.CreateShieldedTransaction)
		fullnode.POST("/createshieldedtransactionwithoutspendauthsig", r.contractHandler.CreateShieldedTransactionWithoutSpendAuthSig)
		fullnode.POST("/getmerkletreevoucherinfo", r.contractHandler.GetMerkleTreeVoucherInfo)
		fullnode.POST("/scannotebyivk", r.contractHandler.ScanNoteByIvk)
		fullnode.POST("/scanandmarknotebyivk", r.contractHandler.ScanAndMarkNoteByIvk)
		fullnode.POST("/scannotebyovk", r.contractHandler.ScanNoteByOvk)
		fullnode.POST("/isspend", r.contractHandler.IsSpend)
		fullnode.POST("/getrcm", r.contractHandler.GetRcm)
		fullnode.POST("/createspendauthsig", r.contractHandler.CreateSpendAuthSig)
		fullnode.POST("/createshieldnullifier", r.contractHandler.CreateShieldNullifier)
		fullnode.POST("/getshieldtransactionhash", r.contractHandler.GetShieldTransactionHash)
		fullnode.POST("/scanshieldedlrc20notesbyivk", r.contractHandler.GetShieldedLRC20NotesByIvk)
		fullnode.POST("/scanshieldedlrc20notesbyovk", r.contractHandler.GetShieldedLRC20NotesByOvk)
		fullnode.POST("/isshieldedlrc20contractnotespent", r.contractHandler.IsShieldedLRC20ContractNoteSpent)
		fullnode.POST("/createshieldedcontractparameters", r.contractHandler.CreateShieldedContractParameters)
		fullnode.POST("/createshieldedcontractparameterswithoutask", r.contractHandler.CreateShieldedContractParametersWithoutAsk)
		fullnode.POST("/gettriggerinputforshieldedlrc20contract", r.contractHandler.GetTriggerInputForShieldedLRC20Contract)
		
		// Utility
		fullnode.POST("/validateaddress", r.accountHandler.ValidateAddress)
		fullnode.POST("/generateaddress", r.accountHandler.GenerateAddress)
		fullnode.POST("/createaddress", r.accountHandler.CreateAddress)
		fullnode.POST("/easytransfer", r.accountHandler.EasyTransfer)
		fullnode.POST("/easytransferbyprivate", r.accountHandler.EasyTransferByPrivate)
		fullnode.POST("/easytransferasset", r.accountHandler.EasyTransferAsset)
		fullnode.POST("/easytransferassetbyprivate", r.accountHandler.EasyTransferAssetByPrivate)
		
		// Pending Pool
		fullnode.POST("/gettransactionfrompending", r.transactionHandler.GetTransactionFromPending)
		fullnode.GET("/gettransactionlistfrompending", r.transactionHandler.GetTransactionListFromPending)
		fullnode.GET("/getpendingsize", r.transactionHandler.GetPendingSize)
		
		// Chain Info
		fullnode.GET("/getchainparameters", r.nodeHandler.GetChainParameters)
		fullnode.GET("/totaltransaction", r.nodeHandler.TotalTransaction)
		fullnode.GET("/getburnlind", r.nodeHandler.GetBurnLind)
		fullnode.GET("/getenergyprices", r.nodeHandler.GetEnergyPrices)
		fullnode.GET("/getbandwidthprices", r.nodeHandler.GetBandwidthPrices)
		fullnode.GET("/getmemofee", r.nodeHandler.GetMemoFeePrices)
	}

	// Solidity Node API routes
	solidity := r.engine.Group("/walletsolidity")
	{
		// Account
		solidity.POST("/getaccount", r.accountHandler.GetAccountSolidity)
		solidity.POST("/getaccountbyid", r.accountHandler.GetAccountByIdSolidity)
		solidity.POST("/getaccountresource", r.accountHandler.GetAccountResource)
		solidity.POST("/getaccountnet", r.accountHandler.GetAccountNet)
		
		// Transaction
		solidity.POST("/gettransactionbyid", r.transactionHandler.GetTransactionByIdSolidity)
		solidity.POST("/gettransactioninfobyid", r.transactionHandler.GetTransactionInfoByIdSolidity)
		solidity.POST("/gettransactioncountbyblocknum", r.transactionHandler.GetTransactionCountByBlockNum)
		solidity.POST("/gettransactioninfobyblocknum", r.transactionHandler.GetTransactionInfoByBlockNumSolidity)
		
		// Block
		solidity.POST("/getnowblock", r.blockHandler.GetNowBlockSolidity)
		solidity.POST("/getblockbynum", r.blockHandler.GetBlockByNumSolidity)
		solidity.POST("/getblockbyid", r.blockHandler.GetBlockById)
		solidity.POST("/getblockbylimitnext", r.blockHandler.GetBlockByLimitNext)
		solidity.POST("/getblockbylatestnum", r.blockHandler.GetBlockByLatestNum)
		solidity.POST("/getblock", r.blockHandler.GetBlock)
		
		// Node Info
		solidity.POST("/getnodeinfo", r.nodeHandler.GetNodeInfoSolidity)
		
		// Asset (LRC-10)
		solidity.POST("/getassetissuebyid", r.tokenHandler.GetAssetIssueByIdSolidity)
		solidity.POST("/getassetissuebyname", r.tokenHandler.GetAssetIssueByNameSolidity)
		solidity.GET("/getassetissuelist", r.tokenHandler.GetAssetIssueListSolidity)
		solidity.POST("/getassetissuelistbyname", r.tokenHandler.GetAssetIssueListByName)
		solidity.POST("/getpaginatedassetissuelist", r.tokenHandler.GetPaginatedAssetIssueList)
		
		// Exchange
		solidity.POST("/getexchangebyid", r.contractHandler.GetExchangeById)
		solidity.GET("/listexchanges", r.contractHandler.ListExchanges)
		
		// Market
		solidity.POST("/getmarketorderbyaccount", r.contractHandler.GetMarketOrderByAccount)
		solidity.POST("/getmarketorderbyid", r.contractHandler.GetMarketOrderById)
		solidity.POST("/getmarketpricebypair", r.contractHandler.GetMarketPriceByPair)
		solidity.POST("/getmarketorderlistbypair", r.contractHandler.GetMarketOrderListByPair)
		solidity.GET("/getmarketpairlist", r.contractHandler.GetMarketPairList)
		
		// Smart Contract
		solidity.POST("/triggerconstantcontract", r.contractHandler.TriggerConstantContractSolidity)
		solidity.POST("/estimateenergy", r.contractHandler.EstimateEnergySolidity)
		solidity.POST("/getcontract", r.contractHandler.GetContract)
		solidity.POST("/getcontractinfo", r.contractHandler.GetContractInfo)
		
		// Witness
		solidity.GET("/listwitnesses", r.nodeHandler.ListWitnessesSolidity)
		solidity.POST("/getBrokerage", r.nodeHandler.GetBrokerageSolidity)
		solidity.POST("/getReward", r.nodeHandler.GetRewardSolidity)
		solidity.GET("/getpaginatednowwitnesslist", r.nodeHandler.GetPaginatedNowWitnessListSolidity)
		
		// Resource Delegation
		solidity.POST("/getdelegatedresource", r.contractHandler.GetDelegatedResource)
		solidity.POST("/getdelegatedresourcev2", r.contractHandler.GetDelegatedResourceV2)
		solidity.POST("/getdelegatedresourceaccountindex", r.contractHandler.GetDelegatedResourceAccountIndex)
		solidity.POST("/getdelegatedresourceaccountindexv2", r.contractHandler.GetDelegatedResourceAccountIndexV2)
		solidity.POST("/getcandelegatedmaxsize", r.contractHandler.GetCanDelegatedMaxSize)
		solidity.POST("/getavailableunfreezecount", r.contractHandler.GetAvailableUnfreezeCount)
		solidity.POST("/getcanwithdrawunfreezeamount", r.contractHandler.GetCanWithdrawUnfreezeAmount)
		
		// Shielded Transaction
		solidity.POST("/scannotebyivk", r.contractHandler.ScanNoteByIvkSolidity)
		solidity.POST("/scanandmarknotebyivk", r.contractHandler.ScanAndMarkNoteByIvkSolidity)
		solidity.POST("/scannotebyovk", r.contractHandler.ScanNoteByOvkSolidity)
		solidity.POST("/getmerkletreevoucherinfo", r.contractHandler.GetMerkleTreeVoucherInfo)
		solidity.POST("/isspend", r.contractHandler.IsSpend)
		solidity.POST("/scanshieldedlrc20notesbyivk", r.contractHandler.GetShieldedLRC20NotesByIvk)
		solidity.POST("/scanshieldedlrc20notesbyovk", r.contractHandler.GetShieldedLRC20NotesByOvk)
		solidity.POST("/isshieldedlrc20contractnotespent", r.contractHandler.IsShieldedLRC20ContractNoteSpent)
		
		// Burn LIND
		solidity.GET("/getburnlind", r.nodeHandler.GetBurnLindSolidity)
	}

	// JSON-RPC API
	jsonrpc := r.engine.Group("/jsonrpc")
	{
		jsonrpc.POST("", r.handleJsonRPC)
	}

	// Monitor endpoints
	monitor := r.engine.Group("/monitor")
	{
		monitor.GET("/getstatsinfo", r.nodeHandler.GetStatsInfo)
		monitor.GET("/getnodeinfo", r.nodeHandler.GetNodeInfo)
	}

	// Network endpoints
	net := r.engine.Group("/net")
	{
		net.GET("/listnodes", r.nodeHandler.ListNodes)
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

// Health check handler
func (r *Router) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

// Metrics handler
func (r *Router) handleMetrics(c *gin.Context) {
	// Prometheus metrics would be exposed here
	c.String(200, "metrics endpoint")
}

// Proxy handler for external APIs
func (r *Router) handleProxy(c *gin.Context) {
	var req models.ProxyRequestMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Forward request through blockchain client
	resp, err := r.blockchainClient.ProxyRequest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(int(resp.StatusCode), resp.Body)
}

// Testnet request handler
func (r *Router) handleTestnetRequest(c *gin.Context) {
	var req models.TestnetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Process testnet request
	// This would typically interact with a faucet service
	c.JSON(200, gin.H{
		"success": true,
		"txid":    "test_txid",
	})
}

// CSV export handler
func (r *Router) handleCSVExport(c *gin.Context) {
	var req models.CSVExportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Set CSV headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=export.csv")

	// Generate CSV based on type
	switch req.Type {
	case "block":
		r.blockHandler.ExportBlocksToCSV(c, req)
	case "freezeresource":
		r.statsHandler.ExportFreezeResourceToCSV(c, req)
	case "turnover":
		r.statsHandler.ExportTurnoverToCSV(c, req)
	case "tokenholders":
		r.tokenHandler.ExportTokenHoldersToCSV(c, req)
	default:
		c.JSON(400, gin.H{"error": "invalid export type"})
	}
}

// Monitor handler
func (r *Router) handleMonitor(c *gin.Context) {
	var req models.MonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Process monitor event
	c.JSON(200, gin.H{"success": true})
}

// Logo upload handler
func (r *Router) handleLogoUpload(c *gin.Context) {
	var req models.UploadLogoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Process logo upload
	c.JSON(200, gin.H{
		"success": true,
		"url":     "https://example.com/logo.png",
	})
}

// JSON-RPC handler
func (r *Router) handleJsonRPC(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"jsonrpc": "2.0",
			"id":      nil,
			"error": gin.H{
				"code":    -32700,
				"message": "Parse error",
			},
		})
		return
	}

	// Forward to blockchain client
	resp, err := r.blockchainClient.JsonRpcForward(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{
			"jsonrpc": "2.0",
			"id":      req["id"],
			"error": gin.H{
				"code":    -32000,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(200, resp)
}