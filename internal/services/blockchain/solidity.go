package blockchain

import (
	"context"

	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"google.golang.org/grpc"
)

// SolidityNodeClient wraps the gRPC client for Solidity Node operations
type SolidityNodeClient struct {
	client lindapb.WalletSolidityClient
	conn   *grpc.ClientConn
}

func NewSolidityNodeClient(conn *grpc.ClientConn) *SolidityNodeClient {
	return &SolidityNodeClient{
		client: lindapb.NewWalletSolidityClient(conn),
		conn:   conn,
	}
}

// Account methods (Confirmed)
func (c *SolidityNodeClient) GetAccount(ctx context.Context, req *lindapb.Account) (*lindapb.Account, error) {
	return c.client.GetAccount(ctx, req)
}

func (c *SolidityNodeClient) GetAccountById(ctx context.Context, req *lindapb.Account) (*lindapb.Account, error) {
	return c.client.GetAccountById(ctx, req)
}

func (c *SolidityNodeClient) GetAccountResource(ctx context.Context, req *lindapb.Account) (*lindapb.AccountResourceMessage, error) {
	return c.client.GetAccountResource(ctx, req)
}

func (c *SolidityNodeClient) GetAccountNet(ctx context.Context, req *lindapb.Account) (*lindapb.AccountNetMessage, error) {
	return c.client.GetAccountNet(ctx, req)
}

// Transaction methods (Confirmed)
func (c *SolidityNodeClient) GetTransactionById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Transaction, error) {
	return c.client.GetTransactionById(ctx, req)
}

func (c *SolidityNodeClient) GetTransactionInfoById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.TransactionInfo, error) {
	return c.client.GetTransactionInfoById(ctx, req)
}

func (c *SolidityNodeClient) GetTransactionCountByBlockNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.NumberMessage, error) {
	return c.client.GetTransactionCountByBlockNum(ctx, req)
}

func (c *SolidityNodeClient) GetTransactionInfoByBlockNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.TransactionInfoList, error) {
	return c.client.GetTransactionInfoByBlockNum(ctx, req)
}

// Block methods (Confirmed)
func (c *SolidityNodeClient) GetNowBlock(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.Block, error) {
	return c.client.GetNowBlock(ctx, req)
}

func (c *SolidityNodeClient) GetBlockByNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Block, error) {
	return c.client.GetBlockByNum(ctx, req)
}

func (c *SolidityNodeClient) GetBlockById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Block, error) {
	return c.client.GetBlockById(ctx, req)
}

func (c *SolidityNodeClient) GetBlockByLimitNext(ctx context.Context, req *lindapb.BlockLimit) (*lindapb.BlockList, error) {
	return c.client.GetBlockByLimitNext(ctx, req)
}

func (c *SolidityNodeClient) GetBlockByLatestNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.BlockList, error) {
	return c.client.GetBlockByLatestNum(ctx, req)
}

func (c *SolidityNodeClient) GetBlock(ctx context.Context, req *lindapb.BlockReq) (*lindapb.BlockExtention, error) {
	return c.client.GetBlock(ctx, req)
}

// Node info
func (c *SolidityNodeClient) GetNodeInfo(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeInfo, error) {
	return c.client.GetNodeInfo(ctx, req)
}

// Asset methods (Confirmed)
func (c *SolidityNodeClient) GetAssetIssueById(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.AssetIssueContract, error) {
	return c.client.GetAssetIssueById(ctx, req)
}

func (c *SolidityNodeClient) GetAssetIssueByName(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueContract, error) {
	return c.client.GetAssetIssueByName(ctx, req)
}

func (c *SolidityNodeClient) GetAssetIssueList(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AssetIssueList, error) {
	return c.client.GetAssetIssueList(ctx, req)
}

func (c *SolidityNodeClient) GetAssetIssueListByName(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueList, error) {
	return c.client.GetAssetIssueListByName(ctx, req)
}

func (c *SolidityNodeClient) GetPaginatedAssetIssueList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.AssetIssueList, error) {
	return c.client.GetPaginatedAssetIssueList(ctx, req)
}

// Exchange methods (Confirmed)
func (c *SolidityNodeClient) GetExchangeById(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Exchange, error) {
	return c.client.GetExchangeById(ctx, req)
}

func (c *SolidityNodeClient) ListExchanges(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ExchangeList, error) {
	return c.client.ListExchanges(ctx, req)
}

// Market methods (Confirmed)
func (c *SolidityNodeClient) GetMarketOrderByAccount(ctx context.Context, req *lindapb.MarketOrderReq) (*lindapb.MarketOrderList, error) {
	return c.client.GetMarketOrderByAccount(ctx, req)
}

func (c *SolidityNodeClient) GetMarketOrderById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.MarketOrder, error) {
	return c.client.GetMarketOrderById(ctx, req)
}

func (c *SolidityNodeClient) GetMarketPriceByPair(ctx context.Context, req *lindapb.MarketPriceReq) (*lindapb.MarketPriceList, error) {
	return c.client.GetMarketPriceByPair(ctx, req)
}

func (c *SolidityNodeClient) GetMarketOrderListByPair(ctx context.Context, req *lindapb.MarketOrderListReq) (*lindapb.MarketOrderList, error) {
	return c.client.GetMarketOrderListByPair(ctx, req)
}

func (c *SolidityNodeClient) GetMarketPairList(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.MarketPairList, error) {
	return c.client.GetMarketPairList(ctx, req)
}

// Smart Contract methods (Confirmed)
func (c *SolidityNodeClient) TriggerConstantContract(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	return c.client.TriggerConstantContract(ctx, req)
}

func (c *SolidityNodeClient) EstimateEnergy(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.EstimateEnergyResponse, error) {
	return c.client.EstimateEnergy(ctx, req)
}

func (c *SolidityNodeClient) GetContract(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.SmartContract, error) {
	return c.client.GetContract(ctx, req)
}

func (c *SolidityNodeClient) GetContractInfo(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.ContractInfo, error) {
	return c.client.GetContractInfo(ctx, req)
}

// Witness methods (Confirmed)
func (c *SolidityNodeClient) ListWitnesses(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.WitnessList, error) {
	return c.client.ListWitnesses(ctx, req)
}

func (c *SolidityNodeClient) GetBrokerage(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	return c.client.GetBrokerage(ctx, req)
}

func (c *SolidityNodeClient) GetReward(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	return c.client.GetReward(ctx, req)
}

func (c *SolidityNodeClient) GetPaginatedNowWitnessList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.WitnessList, error) {
	return c.client.GetPaginatedNowWitnessList(ctx, req)
}

// Resource Delegation methods (Confirmed)
func (c *SolidityNodeClient) GetDelegatedResource(ctx context.Context, req *lindapb.DelegatedResourceReq) (*lindapb.DelegatedResourceList, error) {
	return c.client.GetDelegatedResource(ctx, req)
}

func (c *SolidityNodeClient) GetDelegatedResourceV2(ctx context.Context, req *lindapb.DelegatedResourceReq) (*lindapb.DelegatedResourceList, error) {
	return c.client.GetDelegatedResourceV2(ctx, req)
}

func (c *SolidityNodeClient) GetDelegatedResourceAccountIndex(ctx context.Context, req *lindapb.Account) (*lindapb.DelegatedResourceAccountIndex, error) {
	return c.client.GetDelegatedResourceAccountIndex(ctx, req)
}

func (c *SolidityNodeClient) GetDelegatedResourceAccountIndexV2(ctx context.Context, req *lindapb.Account) (*lindapb.DelegatedResourceAccountIndex, error) {
	return c.client.GetDelegatedResourceAccountIndexV2(ctx, req)
}

func (c *SolidityNodeClient) GetCanDelegatedMaxSize(ctx context.Context, req *lindapb.CanDelegatedMaxSizeReq) (*lindapb.NumberMessage, error) {
	return c.client.GetCanDelegatedMaxSize(ctx, req)
}

func (c *SolidityNodeClient) GetAvailableUnfreezeCount(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	return c.client.GetAvailableUnfreezeCount(ctx, req)
}

func (c *SolidityNodeClient) GetCanWithdrawUnfreezeAmount(ctx context.Context, req *lindapb.CanWithdrawUnfreezeAmountReq) (*lindapb.NumberMessage, error) {
	return c.client.GetCanWithdrawUnfreezeAmount(ctx, req)
}

// Shielded Transaction methods (Confirmed)
func (c *SolidityNodeClient) ScanNoteByIvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	return c.client.ScanNoteByIvk(ctx, req)
}

func (c *SolidityNodeClient) ScanAndMarkNoteByIvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	return c.client.ScanAndMarkNoteByIvk(ctx, req)
}

func (c *SolidityNodeClient) ScanNoteByOvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	return c.client.ScanNoteByOvk(ctx, req)
}

func (c *SolidityNodeClient) GetMerkleTreeVoucherInfo(ctx context.Context, req *lindapb.MerkleTreeVoucherReq) (*lindapb.MerkleTreeVoucherInfo, error) {
	return c.client.GetMerkleTreeVoucherInfo(ctx, req)
}

func (c *SolidityNodeClient) IsSpend(ctx context.Context, req *lindapb.IsSpendReq) (*lindapb.IsSpendResponse, error) {
	return c.client.IsSpend(ctx, req)
}

func (c *SolidityNodeClient) GetShieldedLRC20NotesByIvk(ctx context.Context, req *lindapb.ShieldedLRC20NoteReq) (*lindapb.ShieldedLRC20NoteResponse, error) {
	return c.client.GetShieldedLRC20NotesByIvk(ctx, req)
}

func (c *SolidityNodeClient) GetShieldedLRC20NotesByOvk(ctx context.Context, req *lindapb.ShieldedLRC20NoteReq) (*lindapb.ShieldedLRC20NoteResponse, error) {
	return c.client.GetShieldedLRC20NotesByOvk(ctx, req)
}

func (c *SolidityNodeClient) IsShieldedLRC20ContractNoteSpent(ctx context.Context, req *lindapb.IsShieldedLRC20NoteSpentReq) (*lindapb.IsShieldedLRC20NoteSpentResponse, error) {
	return c.client.IsShieldedLRC20ContractNoteSpent(ctx, req)
}

// Burn LIND
func (c *SolidityNodeClient) GetBurnLind(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	return c.client.GetBurnLind(ctx, req)
}