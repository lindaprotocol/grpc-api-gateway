package blockchain

import (
	"context"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"google.golang.org/grpc"
)

type Client struct {
	fullnodeClient   lindapb.WalletClient
	solidityClient   lindapb.WalletSolidityClient
	jsonRpcClient    lindapb.JsonRpcClient
	eventClient      lindapb.EventServiceClient
	lindascanClient  lindapb.LindascanServer
	config           LindaConfig
}

type LindaConfig struct {
	FullnodeEndpoint string
	SolidityEndpoint string
	EventEndpoint    string
	GRPCTimeout      time.Duration
	MaxMsgSize       int
}

func NewClient(conn, solidityConn *grpc.ClientConn, cfg LindaConfig) *Client {
	return &Client{
		fullnodeClient:  lindapb.NewWalletClient(conn),
		solidityClient:  lindapb.NewWalletSolidityClient(solidityConn),
		jsonRpcClient:   lindapb.NewJsonRpcClient(conn),
		// Event client would need its own connection
		config:          cfg,
	}
}

// ==================== Wallet Service Methods ====================

func (c *Client) GetAccount(ctx context.Context, req *lindapb.Account) (*lindapb.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAccount(ctx, req)
}

func (c *Client) GetAccountBalance(ctx context.Context, req *lindapb.AccountBalanceRequest) (*lindapb.AccountBalanceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAccountBalance(ctx, req)
}

func (c *Client) GetAccountResource(ctx context.Context, req *lindapb.Account) (*lindapb.AccountResourceMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAccountResource(ctx, req)
}

func (c *Client) CreateTransaction(ctx context.Context, req *lindapb.TransferContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateTransaction(ctx, req)
}

func (c *Client) BroadcastTransaction(ctx context.Context, req *lindapb.Transaction) (*lindapb.Return, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.BroadcastTransaction(ctx, req)
}

func (c *Client) GetTransactionById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionById(ctx, req)
}

func (c *Client) GetTransactionInfoById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.TransactionInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionInfoById(ctx, req)
}

func (c *Client) GetNowBlock(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetNowBlock(ctx, req)
}

func (c *Client) GetBlockByNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBlockByNum(ctx, req)
}

func (c *Client) GetBlockByLimitNext(ctx context.Context, req *lindapb.BlockLimit) (*lindapb.BlockList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBlockByLimitNext(ctx, req)
}

func (c *Client) ListNodes(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ListNodes(ctx, req)
}

func (c *Client) GetNodeInfo(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetNodeInfo(ctx, req)
}

func (c *Client) GetAssetIssueList(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AssetIssueList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAssetIssueList(ctx, req)
}

func (c *Client) GetAssetIssueById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueContract, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAssetIssueById(ctx, req)
}

func (c *Client) FreezeBalance(ctx context.Context, req *lindapb.FreezeBalanceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.FreezeBalance(ctx, req)
}

func (c *Client) UnfreezeBalance(ctx context.Context, req *lindapb.UnfreezeBalanceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UnfreezeBalance(ctx, req)
}

func (c *Client) FreezeBalanceV2(ctx context.Context, req *lindapb.FreezeBalanceV2Contract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.FreezeBalanceV2(ctx, req)
}

func (c *Client) UnfreezeBalanceV2(ctx context.Context, req *lindapb.UnfreezeBalanceV2Contract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UnfreezeBalanceV2(ctx, req)
}

func (c *Client) DelegateResource(ctx context.Context, req *lindapb.DelegateResourceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.DelegateResource(ctx, req)
}

func (c *Client) UnDelegateResource(ctx context.Context, req *lindapb.UnDelegateResourceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UnDelegateResource(ctx, req)
}

func (c *Client) ListWitnesses(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.WitnessList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ListWitnesses(ctx, req)
}

func (c *Client) VoteWitnessAccount(ctx context.Context, req *lindapb.VoteWitnessContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.VoteWitnessAccount(ctx, req)
}

func (c *Client) GetBrokerage(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBrokerage(ctx, req)
}

func (c *Client) GetReward(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetReward(ctx, req)
}

func (c *Client) ListProposals(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ProposalList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ListProposals(ctx, req)
}

func (c *Client) GetProposalById(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Proposal, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetProposalById(ctx, req)
}

func (c *Client) ProposalCreate(ctx context.Context, req *lindapb.ProposalCreateContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ProposalCreate(ctx, req)
}

func (c *Client) GetExchangeById(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Exchange, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetExchangeById(ctx, req)
}

func (c *Client) ListExchanges(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ExchangeList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ListExchanges(ctx, req)
}

func (c *Client) GetMarketOrderByAccount(ctx context.Context, req *lindapb.MarketOrderReq) (*lindapb.MarketOrderList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMarketOrderByAccount(ctx, req)
}

func (c *Client) GetMarketPriceByPair(ctx context.Context, req *lindapb.MarketPriceReq) (*lindapb.MarketPriceList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMarketPriceByPair(ctx, req)
}

func (c *Client) GetMarketPairList(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.MarketPairList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMarketPairList(ctx, req)
}

func (c *Client) DeployContract(ctx context.Context, req *lindapb.CreateSmartContract) (*lindapb.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.DeployContract(ctx, req)
}

func (c *Client) TriggerSmartContract(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.TriggerSmartContract(ctx, req)
}

func (c *Client) TriggerConstantContract(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.TriggerConstantContract(ctx, req)
}

func (c *Client) EstimateEnergy(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.EstimateEnergyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.EstimateEnergy(ctx, req)
}

func (c *Client) GetContract(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.SmartContract, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetContract(ctx, req)
}

func (c *Client) ValidateAddress(ctx context.Context, req *lindapb.AddressMessage) (*lindapb.AddressValidateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ValidateAddress(ctx, req)
}

func (c *Client) GenerateAddress(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AddressPrKeyPairMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GenerateAddress(ctx, req)
}

func (c *Client) GetChainParameters(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ChainParameters, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetChainParameters(ctx, req)
}

func (c *Client) TotalTransaction(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.TotalTransaction(ctx, req)
}

func (c *Client) GetBurnLind(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBurnLind(ctx, req)
}

// ==================== Solidity Node Methods ====================

func (c *Client) GetAccountSolidity(ctx context.Context, req *lindapb.Account) (*lindapb.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetAccount(ctx, req)
}

func (c *Client) GetTransactionByIdSolidity(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetTransactionById(ctx, req)
}

func (c *Client) GetTransactionInfoByIdSolidity(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.TransactionInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetTransactionInfoById(ctx, req)
}

func (c *Client) GetNowBlockSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetNowBlock(ctx, req)
}

func (c *Client) GetBlockByNumSolidity(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetBlockByNum(ctx, req)
}

func (c *Client) GetAssetIssueListSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AssetIssueList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetAssetIssueList(ctx, req)
}

func (c *Client) ListWitnessesSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.WitnessList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.ListWitnesses(ctx, req)
}

func (c *Client) GetNodeInfoSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetNodeInfo(ctx, req)
}

func (c *Client) TriggerConstantContractSolidity(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.TriggerConstantContract(ctx, req)
}

// ==================== JSON-RPC Methods ====================

func (c *Client) JsonRpcForward(ctx context.Context, req *lindapb.JsonRpcRequest) (*lindapb.JsonRpcResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.jsonRpcClient.Forward(ctx, req)
}

// ==================== Lindascan Custom Methods ====================

func (c *Client) GetHomepageBundle(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.HomepageBundle, error) {
	// Implementation with caching and aggregation
	return c.lindascanClient.GetHomepageBundle(ctx, req)
}

func (c *Client) GetNodeMap(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeMapResponse, error) {
	return c.lindascanClient.GetNodeMap(ctx, req)
}

func (c *Client) GetTop10(ctx context.Context, req *lindapb.Top10Request) (*lindapb.Top10Response, error) {
	return c.lindascanClient.GetTop10(ctx, req)
}

func (c *Client) GetTokens(ctx context.Context, req *lindapb.TokenListRequest) (*lindapb.TokenListResponse, error) {
	return c.lindascanClient.GetTokens(ctx, req)
}

func (c *Client) GetLRC20Tokens(ctx context.Context, req *lindapb.LRC20TokenRequest) (*lindapb.LRC20TokenListResponse, error) {
	return c.lindascanClient.GetLRC20Tokens(ctx, req)
}

func (c *Client) GetTokenHolders(ctx context.Context, req *lindapb.TokenHoldersRequest) (*lindapb.TokenHoldersResponse, error) {
	return c.lindascanClient.GetTokenHolders(ctx, req)
}

func (c *Client) GetAccountList(ctx context.Context, req *lindapb.AccountListRequest) (*lindapb.AccountListResponse, error) {
	return c.lindascanClient.GetAccountList(ctx, req)
}

func (c *Client) GetTags(ctx context.Context, req *lindapb.TagRequest) (*lindapb.TagListResponse, error) {
	return c.lindascanClient.GetTags(ctx, req)
}

func (c *Client) InsertTag(ctx context.Context, req *lindapb.TagInsertRequest) (*lindapb.TagResponse, error) {
	return c.lindascanClient.InsertTag(ctx, req)
}

func (c *Client) Search(ctx context.Context, req *lindapb.SearchRequest) (*lindapb.SearchResponse, error) {
	return c.lindascanClient.Search(ctx, req)
}