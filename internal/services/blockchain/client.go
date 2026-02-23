// internal/services/blockchain/client.go
package blockchain

import (
	"context"

	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"google.golang.org/grpc"
)

// Client struct - now using config.LindaConfig directly
type Client struct {
	fullnodeClient   lindapb.WalletClient
	solidityClient   lindapb.WalletSolidityClient
	jsonRpcClient    lindapb.JsonRpcClient
	eventClient      lindapb.EventServiceClient
	config           config.LindaConfig // Use config.LindaConfig directly
}

// NewClient creates a new blockchain client
func NewClient(conn, solidityConn *grpc.ClientConn, cfg config.LindaConfig) *Client {
	return &Client{
		fullnodeClient:  lindapb.NewWalletClient(conn),
		solidityClient:  lindapb.NewWalletSolidityClient(solidityConn),
		jsonRpcClient:   lindapb.NewJsonRpcClient(conn),
		config:          cfg,
	}
}

// ==================== Wallet Service Methods ====================

// GetAccount retrieves account information
func (c *Client) GetAccount(ctx context.Context, req *lindapb.Account) (*lindapb.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAccount(ctx, req)
}

// GetAccountBalance retrieves account balance at a specific block
func (c *Client) GetAccountBalance(ctx context.Context, req *lindapb.AccountBalanceRequest) (*lindapb.AccountBalanceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAccountBalance(ctx, req)
}

// GetAccountResource retrieves account resource information
func (c *Client) GetAccountResource(ctx context.Context, req *lindapb.Account) (*lindapb.AccountResourceMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAccountResource(ctx, req)
}

// GetAccountNet retrieves account bandwidth information
func (c *Client) GetAccountNet(ctx context.Context, req *lindapb.Account) (*lindapb.AccountNetMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAccountNet(ctx, req)
}

// CreateAccount creates a new account
func (c *Client) CreateAccount(ctx context.Context, req *lindapb.AccountCreateContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateAccount(ctx, req)
}

// UpdateAccount updates account name
func (c *Client) UpdateAccount(ctx context.Context, req *lindapb.AccountUpdateContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UpdateAccount(ctx, req)
}

// SetAccountId sets account ID
func (c *Client) SetAccountId(ctx context.Context, req *lindapb.SetAccountIdContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.SetAccountId(ctx, req)
}

// AccountPermissionUpdate updates account permissions
func (c *Client) AccountPermissionUpdate(ctx context.Context, req *lindapb.AccountPermissionUpdateContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.AccountPermissionUpdate(ctx, req)
}

// CreateTransaction creates a TRX transfer transaction
func (c *Client) CreateTransaction(ctx context.Context, req *lindapb.TransferContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateTransaction(ctx, req)
}

// BroadcastTransaction broadcasts a signed transaction
func (c *Client) BroadcastTransaction(ctx context.Context, req *lindapb.Transaction) (*lindapb.Return, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.BroadcastTransaction(ctx, req)
}

// BroadcastHex broadcasts a transaction from hex string
func (c *Client) BroadcastHex(ctx context.Context, req *lindapb.BroadcastHexMessage) (*lindapb.Return, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.BroadcastHex(ctx, req)
}

// GetTransactionById retrieves a transaction by ID
func (c *Client) GetTransactionById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionById(ctx, req)
}

// GetTransactionInfoById retrieves transaction info by ID
func (c *Client) GetTransactionInfoById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.TransactionInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionInfoById(ctx, req)
}

// GetTransactionReceiptById retrieves transaction receipt by ID
func (c *Client) GetTransactionReceiptById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.TransactionReceipt, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionReceiptById(ctx, req)
}

// GetTransactionCountByBlockNum gets transaction count in a block
func (c *Client) GetTransactionCountByBlockNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionCountByBlockNum(ctx, req)
}

// GetTransactionInfoByBlockNum gets all transaction infos in a block
func (c *Client) GetTransactionInfoByBlockNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.TransactionInfoList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionInfoByBlockNum(ctx, req)
}

// GetTransactionSign signs a transaction
func (c *Client) GetTransactionSign(ctx context.Context, req *lindapb.TransactionSign) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionSign(ctx, req)
}

// GetTransactionSignWeight gets sign weight of a transaction
func (c *Client) GetTransactionSignWeight(ctx context.Context, req *lindapb.Transaction) (*lindapb.TransactionSignWeight, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionSignWeight(ctx, req)
}

// GetTransactionApprovedList gets list of addresses that signed the transaction
func (c *Client) GetTransactionApprovedList(ctx context.Context, req *lindapb.Transaction) (*lindapb.TransactionApprovedList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionApprovedList(ctx, req)
}

// GetNowBlock gets the most recent block
func (c *Client) GetNowBlock(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetNowBlock(ctx, req)
}

// GetBlockByNum gets block by number
func (c *Client) GetBlockByNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBlockByNum(ctx, req)
}

// GetBlockById gets block by hash
func (c *Client) GetBlockById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBlockById(ctx, req)
}

// GetBlockByLimitNext gets blocks in range
func (c *Client) GetBlockByLimitNext(ctx context.Context, req *lindapb.BlockLimit) (*lindapb.BlockList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBlockByLimitNext(ctx, req)
}

// GetBlockByLatestNum gets latest N blocks
func (c *Client) GetBlockByLatestNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.BlockList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBlockByLatestNum(ctx, req)
}

// GetBlock gets block by id or number
func (c *Client) GetBlock(ctx context.Context, req *lindapb.BlockReq) (*lindapb.BlockExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBlock(ctx, req)
}

// GetBlockBalance gets balance changes in block
func (c *Client) GetBlockBalance(ctx context.Context, req *lindapb.BlockBalanceReq) (*lindapb.BlockBalance, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBlockBalance(ctx, req)
}

// ListNodes lists connected nodes
func (c *Client) ListNodes(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ListNodes(ctx, req)
}

// GetNodeInfo gets node information
func (c *Client) GetNodeInfo(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetNodeInfo(ctx, req)
}

// GetAssetIssueByAccount gets TRC-10 tokens issued by an account
func (c *Client) GetAssetIssueByAccount(ctx context.Context, req *lindapb.Account) (*lindapb.AssetIssueList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAssetIssueByAccount(ctx, req)
}

// GetAssetIssueById gets TRC-10 token by ID
func (c *Client) GetAssetIssueById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueContract, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAssetIssueById(ctx, req)
}

// GetAssetIssueByName gets TRC-10 token by name
func (c *Client) GetAssetIssueByName(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueContract, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAssetIssueByName(ctx, req)
}

// GetAssetIssueList gets all TRC-10 tokens
func (c *Client) GetAssetIssueList(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AssetIssueList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAssetIssueList(ctx, req)
}

// GetAssetIssueListByName gets all TRC-10 tokens with given name
func (c *Client) GetAssetIssueListByName(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAssetIssueListByName(ctx, req)
}

// GetPaginatedAssetIssueList gets paginated TRC-10 tokens
func (c *Client) GetPaginatedAssetIssueList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.AssetIssueList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetPaginatedAssetIssueList(ctx, req)
}

// CreateAssetIssue creates a new TRC-10 token
func (c *Client) CreateAssetIssue(ctx context.Context, req *lindapb.AssetIssueContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateAssetIssue(ctx, req)
}

// TransferAsset transfers TRC-10 token
func (c *Client) TransferAsset(ctx context.Context, req *lindapb.TransferAssetContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.TransferAsset(ctx, req)
}

// ParticipateAssetIssue participates in token issuance
func (c *Client) ParticipateAssetIssue(ctx context.Context, req *lindapb.ParticipateAssetIssueContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ParticipateAssetIssue(ctx, req)
}

// UnfreezeAsset unfreezes frozen TRC-10 token
func (c *Client) UnfreezeAsset(ctx context.Context, req *lindapb.UnfreezeAssetContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UnfreezeAsset(ctx, req)
}

// UpdateAsset updates TRC-10 token information
func (c *Client) UpdateAsset(ctx context.Context, req *lindapb.UpdateAssetContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UpdateAsset(ctx, req)
}

// FreezeBalance stakes TRX for resources (Stake 1.0)
func (c *Client) FreezeBalance(ctx context.Context, req *lindapb.FreezeBalanceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.FreezeBalance(ctx, req)
}

// UnfreezeBalance unstakes TRX (Stake 1.0)
func (c *Client) UnfreezeBalance(ctx context.Context, req *lindapb.UnfreezeBalanceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UnfreezeBalance(ctx, req)
}

// WithdrawBalance withdraws rewards
func (c *Client) WithdrawBalance(ctx context.Context, req *lindapb.WithdrawBalanceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.WithdrawBalance(ctx, req)
}

// FreezeBalanceV2 stakes TRX for resources (Stake 2.0)
func (c *Client) FreezeBalanceV2(ctx context.Context, req *lindapb.FreezeBalanceV2Contract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.FreezeBalanceV2(ctx, req)
}

// UnfreezeBalanceV2 unstakes TRX (Stake 2.0)
func (c *Client) UnfreezeBalanceV2(ctx context.Context, req *lindapb.UnfreezeBalanceV2Contract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UnfreezeBalanceV2(ctx, req)
}

// WithdrawExpireUnfreeze withdraws unfrozen balance (Stake 2.0)
func (c *Client) WithdrawExpireUnfreeze(ctx context.Context, req *lindapb.WithdrawExpireUnfreezeContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.WithdrawExpireUnfreeze(ctx, req)
}

// DelegateResource delegates resources (Stake 2.0)
func (c *Client) DelegateResource(ctx context.Context, req *lindapb.DelegateResourceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.DelegateResource(ctx, req)
}

// UnDelegateResource cancels resource delegation (Stake 2.0)
func (c *Client) UnDelegateResource(ctx context.Context, req *lindapb.UnDelegateResourceContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UnDelegateResource(ctx, req)
}

// CancelAllUnfreezeV2 cancels all unstaking (Stake 2.0)
func (c *Client) CancelAllUnfreezeV2(ctx context.Context, req *lindapb.CancelAllUnfreezeV2Contract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CancelAllUnfreezeV2(ctx, req)
}

// GetAvailableUnfreezeCount gets remaining unstake operations (Stake 2.0)
func (c *Client) GetAvailableUnfreezeCount(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAvailableUnfreezeCount(ctx, req)
}

// GetCanWithdrawUnfreezeAmount gets withdrawable balance (Stake 2.0)
func (c *Client) GetCanWithdrawUnfreezeAmount(ctx context.Context, req *lindapb.CanWithdrawUnfreezeAmountReq) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetCanWithdrawUnfreezeAmount(ctx, req)
}

// GetDelegatedResource gets delegated resources (Stake 1.0)
func (c *Client) GetDelegatedResource(ctx context.Context, req *lindapb.DelegatedResourceReq) (*lindapb.DelegatedResourceList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetDelegatedResource(ctx, req)
}

// GetDelegatedResourceV2 gets delegated resources (Stake 2.0)
func (c *Client) GetDelegatedResourceV2(ctx context.Context, req *lindapb.DelegatedResourceReq) (*lindapb.DelegatedResourceList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetDelegatedResourceV2(ctx, req)
}

// GetDelegatedResourceAccountIndex gets resource delegation index (Stake 1.0)
func (c *Client) GetDelegatedResourceAccountIndex(ctx context.Context, req *lindapb.Account) (*lindapb.DelegatedResourceAccountIndex, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetDelegatedResourceAccountIndex(ctx, req)
}

// GetDelegatedResourceAccountIndexV2 gets resource delegation index (Stake 2.0)
func (c *Client) GetDelegatedResourceAccountIndexV2(ctx context.Context, req *lindapb.Account) (*lindapb.DelegatedResourceAccountIndex, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetDelegatedResourceAccountIndexV2(ctx, req)
}

// GetCanDelegatedMaxSize gets delegatable resources (Stake 2.0)
func (c *Client) GetCanDelegatedMaxSize(ctx context.Context, req *lindapb.CanDelegatedMaxSizeReq) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetCanDelegatedMaxSize(ctx, req)
}

// ListWitnesses lists all witnesses
func (c *Client) ListWitnesses(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.WitnessList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ListWitnesses(ctx, req)
}

// CreateWitness creates a witness
func (c *Client) CreateWitness(ctx context.Context, req *lindapb.WitnessCreateContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateWitness(ctx, req)
}

// UpdateWitness updates witness URL
func (c *Client) UpdateWitness(ctx context.Context, req *lindapb.WitnessUpdateContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UpdateWitness(ctx, req)
}

// VoteWitnessAccount votes for witnesses
func (c *Client) VoteWitnessAccount(ctx context.Context, req *lindapb.VoteWitnessContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.VoteWitnessAccount(ctx, req)
}

// GetBrokerage gets SR brokerage ratio
func (c *Client) GetBrokerage(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBrokerage(ctx, req)
}

// UpdateBrokerage updates SR brokerage setting
func (c *Client) UpdateBrokerage(ctx context.Context, req *lindapb.UpdateBrokerageContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UpdateBrokerage(ctx, req)
}

// GetReward gets voting rewards
func (c *Client) GetReward(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetReward(ctx, req)
}

// GetNextMaintenanceTime gets next maintenance time
func (c *Client) GetNextMaintenanceTime(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetNextMaintenanceTime(ctx, req)
}

// GetPaginatedNowWitnessList gets paginated witness list
func (c *Client) GetPaginatedNowWitnessList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.WitnessList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetPaginatedNowWitnessList(ctx, req)
}

// ListProposals lists all proposals
func (c *Client) ListProposals(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ProposalList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ListProposals(ctx, req)
}

// GetProposalById gets proposal by ID
func (c *Client) GetProposalById(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Proposal, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetProposalById(ctx, req)
}

// ProposalCreate creates a proposal
func (c *Client) ProposalCreate(ctx context.Context, req *lindapb.ProposalCreateContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ProposalCreate(ctx, req)
}

// ProposalApprove approves a proposal
func (c *Client) ProposalApprove(ctx context.Context, req *lindapb.ProposalApproveContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ProposalApprove(ctx, req)
}

// ProposalDelete deletes a proposal
func (c *Client) ProposalDelete(ctx context.Context, req *lindapb.ProposalDeleteContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ProposalDelete(ctx, req)
}

// GetPaginatedProposalList gets paginated proposals
func (c *Client) GetPaginatedProposalList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.ProposalList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetPaginatedProposalList(ctx, req)
}

// ExchangeCreate creates an exchange
func (c *Client) ExchangeCreate(ctx context.Context, req *lindapb.ExchangeCreateContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ExchangeCreate(ctx, req)
}

// ExchangeInject injects into an exchange
func (c *Client) ExchangeInject(ctx context.Context, req *lindapb.ExchangeInjectContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ExchangeInject(ctx, req)
}

// ExchangeWithdraw withdraws from an exchange
func (c *Client) ExchangeWithdraw(ctx context.Context, req *lindapb.ExchangeWithdrawContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ExchangeWithdraw(ctx, req)
}

// ExchangeTransaction performs an exchange transaction
func (c *Client) ExchangeTransaction(ctx context.Context, req *lindapb.ExchangeTransactionContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ExchangeTransaction(ctx, req)
}

// GetExchangeById gets exchange by ID
func (c *Client) GetExchangeById(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Exchange, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetExchangeById(ctx, req)
}

// ListExchanges lists all exchanges
func (c *Client) ListExchanges(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ExchangeList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ListExchanges(ctx, req)
}

// GetPaginatedExchangeList gets paginated exchanges
func (c *Client) GetPaginatedExchangeList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.ExchangeList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetPaginatedExchangeList(ctx, req)
}

// MarketSellAsset sells asset on market
func (c *Client) MarketSellAsset(ctx context.Context, req *lindapb.MarketSellAssetContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.MarketSellAsset(ctx, req)
}

// MarketCancelOrder cancels market order
func (c *Client) MarketCancelOrder(ctx context.Context, req *lindapb.MarketCancelOrderContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.MarketCancelOrder(ctx, req)
}

// GetMarketOrderByAccount gets market orders by account
func (c *Client) GetMarketOrderByAccount(ctx context.Context, req *lindapb.MarketOrderReq) (*lindapb.MarketOrderList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMarketOrderByAccount(ctx, req)
}

// GetMarketOrderById gets market order by ID
func (c *Client) GetMarketOrderById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.MarketOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMarketOrderById(ctx, req)
}

// GetMarketPriceByPair gets market price by pair
func (c *Client) GetMarketPriceByPair(ctx context.Context, req *lindapb.MarketPriceReq) (*lindapb.MarketPriceList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMarketPriceByPair(ctx, req)
}

// GetMarketOrderListByPair gets market orders by pair
func (c *Client) GetMarketOrderListByPair(ctx context.Context, req *lindapb.MarketOrderListReq) (*lindapb.MarketOrderList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMarketOrderListByPair(ctx, req)
}

// GetMarketPairList gets market pair list
func (c *Client) GetMarketPairList(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.MarketPairList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMarketPairList(ctx, req)
}

// DeployContract deploys a smart contract
func (c *Client) DeployContract(ctx context.Context, req *lindapb.CreateSmartContract) (*lindapb.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.DeployContract(ctx, req)
}

// TriggerSmartContract triggers a smart contract
func (c *Client) TriggerSmartContract(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.TriggerSmartContract(ctx, req)
}

// TriggerConstantContract triggers a constant contract
func (c *Client) TriggerConstantContract(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.TriggerConstantContract(ctx, req)
}

// EstimateEnergy estimates energy for contract execution
func (c *Client) EstimateEnergy(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.EstimateEnergyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.EstimateEnergy(ctx, req)
}

// GetContract gets contract information
func (c *Client) GetContract(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.SmartContract, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetContract(ctx, req)
}

// GetContractInfo gets detailed contract information
func (c *Client) GetContractInfo(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.ContractInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetContractInfo(ctx, req)
}

// UpdateSetting updates contract settings
func (c *Client) UpdateSetting(ctx context.Context, req *lindapb.UpdateSettingContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UpdateSetting(ctx, req)
}

// UpdateEnergyLimit updates contract energy limit
func (c *Client) UpdateEnergyLimit(ctx context.Context, req *lindapb.UpdateEnergyLimitContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.UpdateEnergyLimit(ctx, req)
}

// ClearAbi clears contract ABI
func (c *Client) ClearAbi(ctx context.Context, req *lindapb.ClearAbiContract) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ClearAbi(ctx, req)
}

// GetNewShieldedAddress gets new shielded address
func (c *Client) GetNewShieldedAddress(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ShieldedAddressInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetNewShieldedAddress(ctx, req)
}

// GetSpendingKey gets spending key
func (c *Client) GetSpendingKey(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetSpendingKey(ctx, req)
}

// GetExpandedSpendingKey gets expanded spending key
func (c *Client) GetExpandedSpendingKey(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.ExpandedSpendingKey, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetExpandedSpendingKey(ctx, req)
}

// GetAkFromAsk gets ak from ask
func (c *Client) GetAkFromAsk(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetAkFromAsk(ctx, req)
}

// GetNkFromNsk gets nk from nsk
func (c *Client) GetNkFromNsk(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetNkFromNsk(ctx, req)
}

// GetIncomingViewingKey gets incoming viewing key
func (c *Client) GetIncomingViewingKey(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.IncomingViewingKey, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetIncomingViewingKey(ctx, req)
}

// GetDiversifier gets diversifier
func (c *Client) GetDiversifier(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetDiversifier(ctx, req)
}

// GetZenPaymentAddress gets Zen payment address
func (c *Client) GetZenPaymentAddress(ctx context.Context, req *lindapb.ZenPaymentAddressReq) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetZenPaymentAddress(ctx, req)
}

// CreateShieldedTransaction creates shielded transaction
func (c *Client) CreateShieldedTransaction(ctx context.Context, req *lindapb.ShieldedTransactionReq) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateShieldedTransaction(ctx, req)
}

// CreateShieldedTransactionWithoutSpendAuthSig creates shielded transaction without spend auth sig
func (c *Client) CreateShieldedTransactionWithoutSpendAuthSig(ctx context.Context, req *lindapb.ShieldedTransactionReq) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateShieldedTransactionWithoutSpendAuthSig(ctx, req)
}

// GetMerkleTreeVoucherInfo gets Merkle tree voucher info
func (c *Client) GetMerkleTreeVoucherInfo(ctx context.Context, req *lindapb.MerkleTreeVoucherReq) (*lindapb.MerkleTreeVoucherInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMerkleTreeVoucherInfo(ctx, req)
}

// ScanNoteByIvk scans note by IVK
func (c *Client) ScanNoteByIvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ScanNoteByIvk(ctx, req)
}

// ScanAndMarkNoteByIvk scans and marks note by IVK
func (c *Client) ScanAndMarkNoteByIvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ScanAndMarkNoteByIvk(ctx, req)
}

// ScanNoteByOvk scans note by OVK
func (c *Client) ScanNoteByOvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ScanNoteByOvk(ctx, req)
}

// IsSpend checks if a note is spent
func (c *Client) IsSpend(ctx context.Context, req *lindapb.IsSpendReq) (*lindapb.IsSpendResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.IsSpend(ctx, req)
}

// GetRcm gets RCM
func (c *Client) GetRcm(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetRcm(ctx, req)
}

// CreateSpendAuthSig creates spend auth signature
func (c *Client) CreateSpendAuthSig(ctx context.Context, req *lindapb.SpendAuthSigReq) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateSpendAuthSig(ctx, req)
}

// CreateShieldNullifier creates shield nullifier
func (c *Client) CreateShieldNullifier(ctx context.Context, req *lindapb.ShieldNullifierReq) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateShieldNullifier(ctx, req)
}

// GetShieldTransactionHash gets shield transaction hash
func (c *Client) GetShieldTransactionHash(ctx context.Context, req *lindapb.Transaction) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetShieldTransactionHash(ctx, req)
}

// GetShieldedLRC20NotesByIvk gets shielded LRC20 notes by IVK
func (c *Client) GetShieldedLRC20NotesByIvk(ctx context.Context, req *lindapb.ShieldedLRC20NoteReq) (*lindapb.ShieldedLRC20NoteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetShieldedLRC20NotesByIvk(ctx, req)
}

// GetShieldedLRC20NotesByOvk gets shielded LRC20 notes by OVK
func (c *Client) GetShieldedLRC20NotesByOvk(ctx context.Context, req *lindapb.ShieldedLRC20NoteReq) (*lindapb.ShieldedLRC20NoteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetShieldedLRC20NotesByOvk(ctx, req)
}

// IsShieldedLRC20ContractNoteSpent checks if shielded LRC20 note is spent
func (c *Client) IsShieldedLRC20ContractNoteSpent(ctx context.Context, req *lindapb.IsShieldedLRC20NoteSpentReq) (*lindapb.IsShieldedLRC20NoteSpentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.IsShieldedLRC20ContractNoteSpent(ctx, req)
}

// CreateShieldedContractParameters creates shielded contract parameters
func (c *Client) CreateShieldedContractParameters(ctx context.Context, req *lindapb.ShieldedContractParametersReq) (*lindapb.ShieldedContractParameters, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateShieldedContractParameters(ctx, req)
}

// CreateShieldedContractParametersWithoutAsk creates shielded contract parameters without ask
func (c *Client) CreateShieldedContractParametersWithoutAsk(ctx context.Context, req *lindapb.ShieldedContractParametersReq) (*lindapb.ShieldedContractParameters, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateShieldedContractParametersWithoutAsk(ctx, req)
}

// GetTriggerInputForShieldedLRC20Contract gets trigger input for shielded LRC20 contract
func (c *Client) GetTriggerInputForShieldedLRC20Contract(ctx context.Context, req *lindapb.ShieldedLRC20TriggerInputReq) (*lindapb.ShieldedLRC20TriggerInput, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTriggerInputForShieldedLRC20Contract(ctx, req)
}

// ValidateAddress validates an address
func (c *Client) ValidateAddress(ctx context.Context, req *lindapb.AddressMessage) (*lindapb.AddressValidateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.ValidateAddress(ctx, req)
}

// GenerateAddress generates a new address
func (c *Client) GenerateAddress(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AddressPrKeyPairMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GenerateAddress(ctx, req)
}

// CreateAddress creates an address
func (c *Client) CreateAddress(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.BytesMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.CreateAddress(ctx, req)
}

// EasyTransfer performs easy transfer
func (c *Client) EasyTransfer(ctx context.Context, req *lindapb.EasyTransferMessage) (*lindapb.EasyTransferResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.EasyTransfer(ctx, req)
}

// EasyTransferByPrivate performs easy transfer by private key
func (c *Client) EasyTransferByPrivate(ctx context.Context, req *lindapb.EasyTransferByPrivateMessage) (*lindapb.EasyTransferResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.EasyTransferByPrivate(ctx, req)
}

// EasyTransferAsset performs easy asset transfer
func (c *Client) EasyTransferAsset(ctx context.Context, req *lindapb.EasyTransferAssetMessage) (*lindapb.EasyTransferResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.EasyTransferAsset(ctx, req)
}

// EasyTransferAssetByPrivate performs easy asset transfer by private key
func (c *Client) EasyTransferAssetByPrivate(ctx context.Context, req *lindapb.EasyTransferAssetByPrivateMessage) (*lindapb.EasyTransferResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.EasyTransferAssetByPrivate(ctx, req)
}

// GetTransactionFromPending gets transaction from pending pool
func (c *Client) GetTransactionFromPending(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionFromPending(ctx, req)
}

// GetTransactionListFromPending gets transaction list from pending pool
func (c *Client) GetTransactionListFromPending(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.TransactionIdList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetTransactionListFromPending(ctx, req)
}

// GetPendingSize gets pending pool size
func (c *Client) GetPendingSize(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetPendingSize(ctx, req)
}

// GetChainParameters gets chain parameters
func (c *Client) GetChainParameters(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ChainParameters, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetChainParameters(ctx, req)
}

// TotalTransaction gets total transaction count
func (c *Client) TotalTransaction(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.TotalTransaction(ctx, req)
}

// GetBurnLind gets burned LIND amount
func (c *Client) GetBurnLind(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBurnLind(ctx, req)
}

// GetEnergyPrices gets energy price history
func (c *Client) GetEnergyPrices(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.PriceHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetEnergyPrices(ctx, req)
}

// GetBandwidthPrices gets bandwidth price history
func (c *Client) GetBandwidthPrices(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.PriceHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetBandwidthPrices(ctx, req)
}

// GetMemoFeePrices gets memo fee price history
func (c *Client) GetMemoFeePrices(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.PriceHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetMemoFeePrices(ctx, req)
}

// GetStatsInfo gets stats information
func (c *Client) GetStatsInfo(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.StatsInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.fullnodeClient.GetStatsInfo(ctx, req)
}

// ==================== Solidity Node Methods ====================

// GetAccountSolidity gets confirmed account information
func (c *Client) GetAccountSolidity(ctx context.Context, req *lindapb.Account) (*lindapb.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetAccount(ctx, req)
}

// GetAccountByIdSolidity gets account by ID from solidity node
func (c *Client) GetAccountByIdSolidity(ctx context.Context, req *lindapb.Account) (*lindapb.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetAccountById(ctx, req)
}

// GetTransactionByIdSolidity gets confirmed transaction by ID
func (c *Client) GetTransactionByIdSolidity(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetTransactionById(ctx, req)
}

// GetTransactionInfoByIdSolidity gets confirmed transaction info by ID
func (c *Client) GetTransactionInfoByIdSolidity(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.TransactionInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetTransactionInfoById(ctx, req)
}

// GetTransactionInfoByBlockNumSolidity gets confirmed transaction infos by block number
func (c *Client) GetTransactionInfoByBlockNumSolidity(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.TransactionInfoList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetTransactionInfoByBlockNum(ctx, req)
}

// GetNowBlockSolidity gets latest confirmed block
func (c *Client) GetNowBlockSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetNowBlock(ctx, req)
}

// GetBlockByNumSolidity gets confirmed block by number
func (c *Client) GetBlockByNumSolidity(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetBlockByNum(ctx, req)
}

// GetAssetIssueByIdSolidity gets confirmed TRC-10 token by ID
func (c *Client) GetAssetIssueByIdSolidity(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.AssetIssueContract, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetAssetIssueById(ctx, req)
}

// GetAssetIssueByNameSolidity gets confirmed TRC-10 token by name
func (c *Client) GetAssetIssueByNameSolidity(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueContract, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetAssetIssueByName(ctx, req)
}

// GetAssetIssueListSolidity gets all confirmed TRC-10 tokens
func (c *Client) GetAssetIssueListSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AssetIssueList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetAssetIssueList(ctx, req)
}

// ListWitnessesSolidity lists all confirmed witnesses
func (c *Client) ListWitnessesSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.WitnessList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.ListWitnesses(ctx, req)
}

// GetNodeInfoSolidity gets confirmed node information
func (c *Client) GetNodeInfoSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetNodeInfo(ctx, req)
}

// GetBrokerageSolidity gets confirmed SR brokerage ratio
func (c *Client) GetBrokerageSolidity(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetBrokerage(ctx, req)
}

// GetRewardSolidity gets confirmed voting rewards
func (c *Client) GetRewardSolidity(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetReward(ctx, req)
}

// GetPaginatedNowWitnessListSolidity gets paginated confirmed witness list
func (c *Client) GetPaginatedNowWitnessListSolidity(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.WitnessList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetPaginatedNowWitnessList(ctx, req)
}

// TriggerConstantContractSolidity triggers constant contract on solidity node
func (c *Client) TriggerConstantContractSolidity(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.TriggerConstantContract(ctx, req)
}

// EstimateEnergySolidity estimates energy on solidity node
func (c *Client) EstimateEnergySolidity(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.EstimateEnergyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.EstimateEnergy(ctx, req)
}

// GetBurnLindSolidity gets burned LIND amount from solidity node
func (c *Client) GetBurnLindSolidity(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.GRPCTimeout)
	defer cancel()
	return c.solidityClient.GetBurnLind(ctx, req)
}

// ==================== JSON-RPC Methods ====================

// JsonRpcForward forwards a JSON-RPC request
func (c *Client) JsonRpcForward(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	// This would need to be implemented with the actual JSON-RPC client
	return nil, nil
}

// ==================== Lindascan Custom Methods ====================

// // GetHomepageBundle gets homepage bundle data
// func (c *Client) GetHomepageBundle(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.HomepageBundle, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetNodeMap gets node map data
// func (c *Client) GetNodeMap(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeMapResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTop10 gets top 10 data
// func (c *Client) GetTop10(ctx context.Context, req *lindapb.Top10Request) (*lindapb.Top10Response, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // ProxyRequest proxies a request to an external API
// func (c *Client) ProxyRequest(ctx context.Context, req *lindapb.ProxyRequestMessage) (*lindapb.ProxyResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTokens gets token list
// func (c *Client) GetTokens(ctx context.Context, req *lindapb.TokenListRequest) (*lindapb.TokenListResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTokenById gets token by ID
// func (c *Client) GetTokenById(ctx context.Context, req *lindapb.TokenByIdRequest) (*lindapb.TokenInfo, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetLRC20Tokens gets LRC20 tokens
// func (c *Client) GetLRC20Tokens(ctx context.Context, req *lindapb.LRC20TokenRequest) (*lindapb.LRC20TokenListResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetLRC20TokenByContract gets LRC20 token by contract address
// func (c *Client) GetLRC20TokenByContract(ctx context.Context, req *lindapb.LRC20TokenByContractRequest) (*lindapb.LRC20TokenInfo, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTokenHolders gets token holders
// func (c *Client) GetTokenHolders(ctx context.Context, req *lindapb.TokenHoldersRequest) (*lindapb.TokenHoldersResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTokenTransfers gets token transfers
// func (c *Client) GetTokenTransfers(ctx context.Context, req *lindapb.TokenTransfersRequest) (*lindapb.TokenTransfersResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTokensOverview gets tokens overview
// func (c *Client) GetTokensOverview(ctx context.Context, req *lindapb.TokensOverviewRequest) (*lindapb.TokensOverviewResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTokenPrice gets token price
// func (c *Client) GetTokenPrice(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.TokenPriceResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetParticipateAssetIssue gets participate asset issue data
// func (c *Client) GetParticipateAssetIssue(ctx context.Context, req *lindapb.ParticipateAssetIssueRequest) (*lindapb.ParticipateAssetIssueResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetWinkFund gets WINK fund data
// func (c *Client) GetWinkFund(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.WinkFundResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetWinkGraphic gets WINK graphic data
// func (c *Client) GetWinkGraphic(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.GraphicResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetJSTFund gets JST fund data
// func (c *Client) GetJSTFund(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.JSTFundResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetJSTGraphic gets JST graphic data
// func (c *Client) GetJSTGraphic(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.GraphicResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetBitTorrentGraphic gets BitTorrent graphic data
// func (c *Client) GetBitTorrentGraphic(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.GraphicResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTokenPositionDistribution gets token position distribution
// func (c *Client) GetTokenPositionDistribution(ctx context.Context, req *lindapb.TokenPositionRequest) (*lindapb.TokenPositionResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetAssetTransfer gets asset transfers
// func (c *Client) GetAssetTransfer(ctx context.Context, req *lindapb.AssetTransferRequest) (*lindapb.AssetTransferResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetAccountList gets account list
// func (c *Client) GetAccountList(ctx context.Context, req *lindapb.AccountListRequest) (*lindapb.AccountListResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetAccountOverview gets account overview
// func (c *Client) GetAccountOverview(ctx context.Context, req *lindapb.AccountOverviewRequest) (*lindapb.AccountOverviewResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetAccountProposals gets account proposals
// func (c *Client) GetAccountProposals(ctx context.Context, req *lindapb.AccountProposalRequest) (*lindapb.AccountProposalResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTags gets tags
// func (c *Client) GetTags(ctx context.Context, req *lindapb.TagRequest) (*lindapb.TagListResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // InsertTag inserts a tag
// func (c *Client) InsertTag(ctx context.Context, req *lindapb.TagInsertRequest) (*lindapb.TagResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // UpdateTag updates a tag
// func (c *Client) UpdateTag(ctx context.Context, req *lindapb.TagUpdateRequest) (*lindapb.TagResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // DeleteTag deletes a tag
// func (c *Client) DeleteTag(ctx context.Context, req *lindapb.TagDeleteRequest) (*lindapb.TagResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // RecommendTag recommends tags
// func (c *Client) RecommendTag(ctx context.Context, req *lindapb.TagRecommendRequest) (*lindapb.TagListResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // UploadLogo uploads a logo
// func (c *Client) UploadLogo(ctx context.Context, req *lindapb.UploadLogoRequest) (*lindapb.UploadLogoResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetVoteInfo gets vote information
// func (c *Client) GetVoteInfo(ctx context.Context, req *lindapb.VoteInfoRequest) (*lindapb.VoteInfoResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetChainParametersV2 gets chain parameters v2
// func (c *Client) GetChainParametersV2(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ChainParametersV2, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetBlocksV2 gets blocks v2
// func (c *Client) GetBlocksV2(ctx context.Context, req *lindapb.BlockListV2Request) (*lindapb.BlockListV2Response, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTransactionsV2 gets transactions v2
// func (c *Client) GetTransactionsV2(ctx context.Context, req *lindapb.TransactionListV2Request) (*lindapb.TransactionListV2Response, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetInternalTransactions gets internal transactions
// func (c *Client) GetInternalTransactions(ctx context.Context, req *lindapb.InternalTransactionRequest) (*lindapb.InternalTransactionResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetContractTransactions gets contract transactions
// func (c *Client) GetContractTransactions(ctx context.Context, req *lindapb.ContractTransactionRequest) (*lindapb.ContractTransactionResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetLRC10LRC20Transfers gets LRC10/LRC20 transfers
// func (c *Client) GetLRC10LRC20Transfers(ctx context.Context, req *lindapb.LRCTransferRequest) (*lindapb.LRCTransferResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetContractAccountHistory gets contract account history
// func (c *Client) GetContractAccountHistory(ctx context.Context, req *lindapb.ContractAccountHistoryRequest) (*lindapb.ContractAccountHistoryResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetSmartContractTriggersBatch gets smart contract triggers batch
// func (c *Client) GetSmartContractTriggersBatch(ctx context.Context, req *lindapb.SmartContractTriggersRequest) (*lindapb.SmartContractTriggersResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetEnergyStatistic gets energy statistics
// func (c *Client) GetEnergyStatistic(ctx context.Context, req *lindapb.EnergyStatisticRequest) (*lindapb.EnergyStatisticResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTriggerStatistic gets trigger statistics
// func (c *Client) GetTriggerStatistic(ctx context.Context, req *lindapb.TriggerStatisticRequest) (*lindapb.TriggerStatisticResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetCallerAddressStatistic gets caller address statistics
// func (c *Client) GetCallerAddressStatistic(ctx context.Context, req *lindapb.CallerAddressStatisticRequest) (*lindapb.CallerAddressStatisticResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetEnergyDailyStatistic gets daily energy statistics
// func (c *Client) GetEnergyDailyStatistic(ctx context.Context, req *lindapb.EnergyDailyStatisticRequest) (*lindapb.EnergyDailyStatisticResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTriggerAmountStatistic gets trigger amount statistics
// func (c *Client) GetTriggerAmountStatistic(ctx context.Context, req *lindapb.TriggerAmountStatisticRequest) (*lindapb.TriggerAmountStatisticResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetFreezeResource gets freeze resource data
// func (c *Client) GetFreezeResource(ctx context.Context, req *lindapb.FreezeResourceRequest) (*lindapb.FreezeResourceResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetTurnover gets turnover data
// func (c *Client) GetTurnover(ctx context.Context, req *lindapb.TurnoverRequest) (*lindapb.TurnoverResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetLindHolderStats gets LIND holder statistics
// func (c *Client) GetLindHolderStats(ctx context.Context, req *lindapb.LindHolderRequest) (*lindapb.LindHolderResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetOneContractEnergyStatistic gets one contract energy statistics
// func (c *Client) GetOneContractEnergyStatistic(ctx context.Context, req *lindapb.OneContractEnergyStatisticRequest) (*lindapb.OneContractEnergyStatisticResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetOneContractTriggerStatistic gets one contract trigger statistics
// func (c *Client) GetOneContractTriggerStatistic(ctx context.Context, req *lindapb.OneContractTriggerStatisticRequest) (*lindapb.OneContractTriggerStatisticResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetOneContractCallerStatistic gets one contract caller statistics
// func (c *Client) GetOneContractCallerStatistic(ctx context.Context, req *lindapb.OneContractCallerStatisticRequest) (*lindapb.OneContractCallerStatisticResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetOneContractCallers gets one contract callers
// func (c *Client) GetOneContractCallers(ctx context.Context, req *lindapb.OneContractCallersRequest) (*lindapb.OneContractCallersResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetNodeOverviewUpload gets node overview upload
// func (c *Client) GetNodeOverviewUpload(ctx context.Context, req *lindapb.NodeOverviewRequest) (*lindapb.NodeOverviewResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetNodeInfoUpload gets node info upload
// func (c *Client) GetNodeInfoUpload(ctx context.Context, req *lindapb.NodeInfoUploadRequest) (*lindapb.NodeInfoUploadResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetLedger gets ledger data
// func (c *Client) GetLedger(ctx context.Context, req *lindapb.LedgerRequest) (*lindapb.LedgerResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // Search performs a search
// func (c *Client) Search(ctx context.Context, req *lindapb.SearchRequest) (*lindapb.SearchResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // GetFund gets fund data
// func (c *Client) GetFund(ctx context.Context, req *lindapb.FundRequest) (*lindapb.FundResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // RequestTestnetCoins requests testnet coins
// func (c *Client) RequestTestnetCoins(ctx context.Context, req *lindapb.TestnetRequest) (*lindapb.TestnetResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // ExportToCSV exports data to CSV
// func (c *Client) ExportToCSV(ctx context.Context, req *lindapb.CSVExportRequest) (*lindapb.CSVExportResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }

// // Monitor handles monitor requests
// func (c *Client) Monitor(ctx context.Context, req *lindapb.MonitorRequest) (*lindapb.MonitorResponse, error) {
// 	// Implementation would call the appropriate handler or service
// 	return nil, nil
// }