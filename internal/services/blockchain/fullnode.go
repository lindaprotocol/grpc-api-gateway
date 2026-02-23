package blockchain

import (
	"context"

	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"google.golang.org/grpc"
)

// FullNodeClient wraps the gRPC client for FullNode operations
type FullNodeClient struct {
	client lindapb.WalletClient
	conn   *grpc.ClientConn
}

func NewFullNodeClient(conn *grpc.ClientConn) *FullNodeClient {
	return &FullNodeClient{
		client: lindapb.NewWalletClient(conn),
		conn:   conn,
	}
}

// Account methods
func (c *FullNodeClient) GetAccount(ctx context.Context, req *lindapb.Account) (*lindapb.Account, error) {
	return c.client.GetAccount(ctx, req)
}

func (c *FullNodeClient) GetAccountBalance(ctx context.Context, req *lindapb.AccountBalanceRequest) (*lindapb.AccountBalanceResponse, error) {
	return c.client.GetAccountBalance(ctx, req)
}

func (c *FullNodeClient) GetAccountResource(ctx context.Context, req *lindapb.Account) (*lindapb.AccountResourceMessage, error) {
	return c.client.GetAccountResource(ctx, req)
}

func (c *FullNodeClient) GetAccountNet(ctx context.Context, req *lindapb.Account) (*lindapb.AccountNetMessage, error) {
	return c.client.GetAccountNet(ctx, req)
}

func (c *FullNodeClient) CreateAccount(ctx context.Context, req *lindapb.AccountCreateContract) (*lindapb.Transaction, error) {
	return c.client.CreateAccount(ctx, req)
}

func (c *FullNodeClient) UpdateAccount(ctx context.Context, req *lindapb.AccountUpdateContract) (*lindapb.Transaction, error) {
	return c.client.UpdateAccount(ctx, req)
}

func (c *FullNodeClient) SetAccountId(ctx context.Context, req *lindapb.SetAccountIdContract) (*lindapb.Transaction, error) {
	return c.client.SetAccountId(ctx, req)
}

func (c *FullNodeClient) AccountPermissionUpdate(ctx context.Context, req *lindapb.AccountPermissionUpdateContract) (*lindapb.Transaction, error) {
	return c.client.AccountPermissionUpdate(ctx, req)
}

// Transaction methods
func (c *FullNodeClient) CreateTransaction(ctx context.Context, req *lindapb.TransferContract) (*lindapb.Transaction, error) {
	return c.client.CreateTransaction(ctx, req)
}

func (c *FullNodeClient) BroadcastTransaction(ctx context.Context, req *lindapb.Transaction) (*lindapb.Return, error) {
	return c.client.BroadcastTransaction(ctx, req)
}

func (c *FullNodeClient) BroadcastHex(ctx context.Context, req *lindapb.BroadcastHexMessage) (*lindapb.Return, error) {
	return c.client.BroadcastHex(ctx, req)
}

func (c *FullNodeClient) GetTransactionById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Transaction, error) {
	return c.client.GetTransactionById(ctx, req)
}

func (c *FullNodeClient) GetTransactionInfoById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.TransactionInfo, error) {
	return c.client.GetTransactionInfoById(ctx, req)
}

func (c *FullNodeClient) GetTransactionReceiptById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.TransactionReceipt, error) {
	return c.client.GetTransactionReceiptById(ctx, req)
}

func (c *FullNodeClient) GetTransactionCountByBlockNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.NumberMessage, error) {
	return c.client.GetTransactionCountByBlockNum(ctx, req)
}

func (c *FullNodeClient) GetTransactionInfoByBlockNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.TransactionInfoList, error) {
	return c.client.GetTransactionInfoByBlockNum(ctx, req)
}

func (c *FullNodeClient) GetTransactionSign(ctx context.Context, req *lindapb.TransactionSign) (*lindapb.Transaction, error) {
	return c.client.GetTransactionSign(ctx, req)
}

func (c *FullNodeClient) GetTransactionSignWeight(ctx context.Context, req *lindapb.Transaction) (*lindapb.TransactionSignWeight, error) {
	return c.client.GetTransactionSignWeight(ctx, req)
}

func (c *FullNodeClient) GetTransactionApprovedList(ctx context.Context, req *lindapb.Transaction) (*lindapb.TransactionApprovedList, error) {
	return c.client.GetTransactionApprovedList(ctx, req)
}

// Block methods
func (c *FullNodeClient) GetNowBlock(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.Block, error) {
	return c.client.GetNowBlock(ctx, req)
}

func (c *FullNodeClient) GetBlockByNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Block, error) {
	return c.client.GetBlockByNum(ctx, req)
}

func (c *FullNodeClient) GetBlockById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Block, error) {
	return c.client.GetBlockById(ctx, req)
}

func (c *FullNodeClient) GetBlockByLimitNext(ctx context.Context, req *lindapb.BlockLimit) (*lindapb.BlockList, error) {
	return c.client.GetBlockByLimitNext(ctx, req)
}

func (c *FullNodeClient) GetBlockByLatestNum(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.BlockList, error) {
	return c.client.GetBlockByLatestNum(ctx, req)
}

func (c *FullNodeClient) GetBlock(ctx context.Context, req *lindapb.BlockReq) (*lindapb.BlockExtention, error) {
	return c.client.GetBlock(ctx, req)
}

func (c *FullNodeClient) GetBlockBalance(ctx context.Context, req *lindapb.BlockBalanceReq) (*lindapb.BlockBalance, error) {
	return c.client.GetBlockBalance(ctx, req)
}

// Node methods
func (c *FullNodeClient) ListNodes(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeList, error) {
	return c.client.ListNodes(ctx, req)
}

func (c *FullNodeClient) GetNodeInfo(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NodeInfo, error) {
	return c.client.GetNodeInfo(ctx, req)
}

// Asset methods (LRC-10)
func (c *FullNodeClient) GetAssetIssueByAccount(ctx context.Context, req *lindapb.Account) (*lindapb.AssetIssueList, error) {
	return c.client.GetAssetIssueByAccount(ctx, req)
}

func (c *FullNodeClient) GetAssetIssueById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueContract, error) {
	return c.client.GetAssetIssueById(ctx, req)
}

func (c *FullNodeClient) GetAssetIssueByName(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueContract, error) {
	return c.client.GetAssetIssueByName(ctx, req)
}

func (c *FullNodeClient) GetAssetIssueList(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AssetIssueList, error) {
	return c.client.GetAssetIssueList(ctx, req)
}

func (c *FullNodeClient) GetAssetIssueListByName(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.AssetIssueList, error) {
	return c.client.GetAssetIssueListByName(ctx, req)
}

func (c *FullNodeClient) GetPaginatedAssetIssueList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.AssetIssueList, error) {
	return c.client.GetPaginatedAssetIssueList(ctx, req)
}

func (c *FullNodeClient) CreateAssetIssue(ctx context.Context, req *lindapb.AssetIssueContract) (*lindapb.Transaction, error) {
	return c.client.CreateAssetIssue(ctx, req)
}

func (c *FullNodeClient) TransferAsset(ctx context.Context, req *lindapb.TransferAssetContract) (*lindapb.Transaction, error) {
	return c.client.TransferAsset(ctx, req)
}

func (c *FullNodeClient) ParticipateAssetIssue(ctx context.Context, req *lindapb.ParticipateAssetIssueContract) (*lindapb.Transaction, error) {
	return c.client.ParticipateAssetIssue(ctx, req)
}

func (c *FullNodeClient) UnfreezeAsset(ctx context.Context, req *lindapb.UnfreezeAssetContract) (*lindapb.Transaction, error) {
	return c.client.UnfreezeAsset(ctx, req)
}

func (c *FullNodeClient) UpdateAsset(ctx context.Context, req *lindapb.UpdateAssetContract) (*lindapb.Transaction, error) {
	return c.client.UpdateAsset(ctx, req)
}

// Resource methods (Stake 1.0)
func (c *FullNodeClient) FreezeBalance(ctx context.Context, req *lindapb.FreezeBalanceContract) (*lindapb.Transaction, error) {
	return c.client.FreezeBalance(ctx, req)
}

func (c *FullNodeClient) UnfreezeBalance(ctx context.Context, req *lindapb.UnfreezeBalanceContract) (*lindapb.Transaction, error) {
	return c.client.UnfreezeBalance(ctx, req)
}

func (c *FullNodeClient) WithdrawBalance(ctx context.Context, req *lindapb.WithdrawBalanceContract) (*lindapb.Transaction, error) {
	return c.client.WithdrawBalance(ctx, req)
}

// Resource methods (Stake 2.0)
func (c *FullNodeClient) FreezeBalanceV2(ctx context.Context, req *lindapb.FreezeBalanceV2Contract) (*lindapb.Transaction, error) {
	return c.client.FreezeBalanceV2(ctx, req)
}

func (c *FullNodeClient) UnfreezeBalanceV2(ctx context.Context, req *lindapb.UnfreezeBalanceV2Contract) (*lindapb.Transaction, error) {
	return c.client.UnfreezeBalanceV2(ctx, req)
}

func (c *FullNodeClient) WithdrawExpireUnfreeze(ctx context.Context, req *lindapb.WithdrawExpireUnfreezeContract) (*lindapb.Transaction, error) {
	return c.client.WithdrawExpireUnfreeze(ctx, req)
}

func (c *FullNodeClient) DelegateResource(ctx context.Context, req *lindapb.DelegateResourceContract) (*lindapb.Transaction, error) {
	return c.client.DelegateResource(ctx, req)
}

func (c *FullNodeClient) UnDelegateResource(ctx context.Context, req *lindapb.UnDelegateResourceContract) (*lindapb.Transaction, error) {
	return c.client.UnDelegateResource(ctx, req)
}

func (c *FullNodeClient) CancelAllUnfreezeV2(ctx context.Context, req *lindapb.CancelAllUnfreezeV2Contract) (*lindapb.Transaction, error) {
	return c.client.CancelAllUnfreezeV2(ctx, req)
}

func (c *FullNodeClient) GetAvailableUnfreezeCount(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	return c.client.GetAvailableUnfreezeCount(ctx, req)
}

func (c *FullNodeClient) GetCanWithdrawUnfreezeAmount(ctx context.Context, req *lindapb.CanWithdrawUnfreezeAmountReq) (*lindapb.NumberMessage, error) {
	return c.client.GetCanWithdrawUnfreezeAmount(ctx, req)
}

func (c *FullNodeClient) GetDelegatedResource(ctx context.Context, req *lindapb.DelegatedResourceReq) (*lindapb.DelegatedResourceList, error) {
	return c.client.GetDelegatedResource(ctx, req)
}

func (c *FullNodeClient) GetDelegatedResourceV2(ctx context.Context, req *lindapb.DelegatedResourceReq) (*lindapb.DelegatedResourceList, error) {
	return c.client.GetDelegatedResourceV2(ctx, req)
}

func (c *FullNodeClient) GetDelegatedResourceAccountIndex(ctx context.Context, req *lindapb.Account) (*lindapb.DelegatedResourceAccountIndex, error) {
	return c.client.GetDelegatedResourceAccountIndex(ctx, req)
}

func (c *FullNodeClient) GetDelegatedResourceAccountIndexV2(ctx context.Context, req *lindapb.Account) (*lindapb.DelegatedResourceAccountIndex, error) {
	return c.client.GetDelegatedResourceAccountIndexV2(ctx, req)
}

func (c *FullNodeClient) GetCanDelegatedMaxSize(ctx context.Context, req *lindapb.CanDelegatedMaxSizeReq) (*lindapb.NumberMessage, error) {
	return c.client.GetCanDelegatedMaxSize(ctx, req)
}

// Witness & Voting methods
func (c *FullNodeClient) ListWitnesses(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.WitnessList, error) {
	return c.client.ListWitnesses(ctx, req)
}

func (c *FullNodeClient) CreateWitness(ctx context.Context, req *lindapb.WitnessCreateContract) (*lindapb.Transaction, error) {
	return c.client.CreateWitness(ctx, req)
}

func (c *FullNodeClient) UpdateWitness(ctx context.Context, req *lindapb.WitnessUpdateContract) (*lindapb.Transaction, error) {
	return c.client.UpdateWitness(ctx, req)
}

func (c *FullNodeClient) VoteWitnessAccount(ctx context.Context, req *lindapb.VoteWitnessContract) (*lindapb.Transaction, error) {
	return c.client.VoteWitnessAccount(ctx, req)
}

func (c *FullNodeClient) GetBrokerage(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	return c.client.GetBrokerage(ctx, req)
}

func (c *FullNodeClient) UpdateBrokerage(ctx context.Context, req *lindapb.UpdateBrokerageContract) (*lindapb.Transaction, error) {
	return c.client.UpdateBrokerage(ctx, req)
}

func (c *FullNodeClient) GetReward(ctx context.Context, req *lindapb.Account) (*lindapb.NumberMessage, error) {
	return c.client.GetReward(ctx, req)
}

func (c *FullNodeClient) GetNextMaintenanceTime(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	return c.client.GetNextMaintenanceTime(ctx, req)
}

func (c *FullNodeClient) GetPaginatedNowWitnessList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.WitnessList, error) {
	return c.client.GetPaginatedNowWitnessList(ctx, req)
}

// Proposal methods
func (c *FullNodeClient) ListProposals(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ProposalList, error) {
	return c.client.ListProposals(ctx, req)
}

func (c *FullNodeClient) GetProposalById(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Proposal, error) {
	return c.client.GetProposalById(ctx, req)
}

func (c *FullNodeClient) ProposalCreate(ctx context.Context, req *lindapb.ProposalCreateContract) (*lindapb.Transaction, error) {
	return c.client.ProposalCreate(ctx, req)
}

func (c *FullNodeClient) ProposalApprove(ctx context.Context, req *lindapb.ProposalApproveContract) (*lindapb.Transaction, error) {
	return c.client.ProposalApprove(ctx, req)
}

func (c *FullNodeClient) ProposalDelete(ctx context.Context, req *lindapb.ProposalDeleteContract) (*lindapb.Transaction, error) {
	return c.client.ProposalDelete(ctx, req)
}

func (c *FullNodeClient) GetPaginatedProposalList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.ProposalList, error) {
	return c.client.GetPaginatedProposalList(ctx, req)
}

// Exchange methods
func (c *FullNodeClient) ExchangeCreate(ctx context.Context, req *lindapb.ExchangeCreateContract) (*lindapb.Transaction, error) {
	return c.client.ExchangeCreate(ctx, req)
}

func (c *FullNodeClient) ExchangeInject(ctx context.Context, req *lindapb.ExchangeInjectContract) (*lindapb.Transaction, error) {
	return c.client.ExchangeInject(ctx, req)
}

func (c *FullNodeClient) ExchangeWithdraw(ctx context.Context, req *lindapb.ExchangeWithdrawContract) (*lindapb.Transaction, error) {
	return c.client.ExchangeWithdraw(ctx, req)
}

func (c *FullNodeClient) ExchangeTransaction(ctx context.Context, req *lindapb.ExchangeTransactionContract) (*lindapb.Transaction, error) {
	return c.client.ExchangeTransaction(ctx, req)
}

func (c *FullNodeClient) GetExchangeById(ctx context.Context, req *lindapb.NumberMessage) (*lindapb.Exchange, error) {
	return c.client.GetExchangeById(ctx, req)
}

func (c *FullNodeClient) ListExchanges(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ExchangeList, error) {
	return c.client.ListExchanges(ctx, req)
}

func (c *FullNodeClient) GetPaginatedExchangeList(ctx context.Context, req *lindapb.PaginatedMessage) (*lindapb.ExchangeList, error) {
	return c.client.GetPaginatedExchangeList(ctx, req)
}

// Market methods (DEX)
func (c *FullNodeClient) MarketSellAsset(ctx context.Context, req *lindapb.MarketSellAssetContract) (*lindapb.Transaction, error) {
	return c.client.MarketSellAsset(ctx, req)
}

func (c *FullNodeClient) MarketCancelOrder(ctx context.Context, req *lindapb.MarketCancelOrderContract) (*lindapb.Transaction, error) {
	return c.client.MarketCancelOrder(ctx, req)
}

func (c *FullNodeClient) GetMarketOrderByAccount(ctx context.Context, req *lindapb.MarketOrderReq) (*lindapb.MarketOrderList, error) {
	return c.client.GetMarketOrderByAccount(ctx, req)
}

func (c *FullNodeClient) GetMarketOrderById(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.MarketOrder, error) {
	return c.client.GetMarketOrderById(ctx, req)
}

func (c *FullNodeClient) GetMarketPriceByPair(ctx context.Context, req *lindapb.MarketPriceReq) (*lindapb.MarketPriceList, error) {
	return c.client.GetMarketPriceByPair(ctx, req)
}

func (c *FullNodeClient) GetMarketOrderListByPair(ctx context.Context, req *lindapb.MarketOrderListReq) (*lindapb.MarketOrderList, error) {
	return c.client.GetMarketOrderListByPair(ctx, req)
}

func (c *FullNodeClient) GetMarketPairList(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.MarketPairList, error) {
	return c.client.GetMarketPairList(ctx, req)
}

// Smart Contract methods
func (c *FullNodeClient) DeployContract(ctx context.Context, req *lindapb.CreateSmartContract) (*lindapb.TransactionExtention, error) {
	return c.client.DeployContract(ctx, req)
}

func (c *FullNodeClient) TriggerSmartContract(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	return c.client.TriggerSmartContract(ctx, req)
}

func (c *FullNodeClient) TriggerConstantContract(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.TransactionExtention, error) {
	return c.client.TriggerConstantContract(ctx, req)
}

func (c *FullNodeClient) EstimateEnergy(ctx context.Context, req *lindapb.TriggerSmartContractReq) (*lindapb.EstimateEnergyResponse, error) {
	return c.client.EstimateEnergy(ctx, req)
}

func (c *FullNodeClient) GetContract(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.SmartContract, error) {
	return c.client.GetContract(ctx, req)
}

func (c *FullNodeClient) GetContractInfo(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.ContractInfo, error) {
	return c.client.GetContractInfo(ctx, req)
}

func (c *FullNodeClient) UpdateSetting(ctx context.Context, req *lindapb.UpdateSettingContract) (*lindapb.Transaction, error) {
	return c.client.UpdateSetting(ctx, req)
}

func (c *FullNodeClient) UpdateEnergyLimit(ctx context.Context, req *lindapb.UpdateEnergyLimitContract) (*lindapb.Transaction, error) {
	return c.client.UpdateEnergyLimit(ctx, req)
}

func (c *FullNodeClient) ClearAbi(ctx context.Context, req *lindapb.ClearAbiContract) (*lindapb.Transaction, error) {
	return c.client.ClearAbi(ctx, req)
}

// Shielded Transaction methods
func (c *FullNodeClient) GetNewShieldedAddress(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ShieldedAddressInfo, error) {
	return c.client.GetNewShieldedAddress(ctx, req)
}

func (c *FullNodeClient) GetSpendingKey(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.BytesMessage, error) {
	return c.client.GetSpendingKey(ctx, req)
}

func (c *FullNodeClient) GetExpandedSpendingKey(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.ExpandedSpendingKey, error) {
	return c.client.GetExpandedSpendingKey(ctx, req)
}

func (c *FullNodeClient) GetAkFromAsk(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.BytesMessage, error) {
	return c.client.GetAkFromAsk(ctx, req)
}

func (c *FullNodeClient) GetNkFromNsk(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.BytesMessage, error) {
	return c.client.GetNkFromNsk(ctx, req)
}

func (c *FullNodeClient) GetIncomingViewingKey(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.IncomingViewingKey, error) {
	return c.client.GetIncomingViewingKey(ctx, req)
}

func (c *FullNodeClient) GetDiversifier(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.BytesMessage, error) {
	return c.client.GetDiversifier(ctx, req)
}

func (c *FullNodeClient) GetZenPaymentAddress(ctx context.Context, req *lindapb.ZenPaymentAddressReq) (*lindapb.BytesMessage, error) {
	return c.client.GetZenPaymentAddress(ctx, req)
}

func (c *FullNodeClient) CreateShieldedTransaction(ctx context.Context, req *lindapb.ShieldedTransactionReq) (*lindapb.Transaction, error) {
	return c.client.CreateShieldedTransaction(ctx, req)
}

func (c *FullNodeClient) CreateShieldedTransactionWithoutSpendAuthSig(ctx context.Context, req *lindapb.ShieldedTransactionReq) (*lindapb.Transaction, error) {
	return c.client.CreateShieldedTransactionWithoutSpendAuthSig(ctx, req)
}

func (c *FullNodeClient) GetMerkleTreeVoucherInfo(ctx context.Context, req *lindapb.MerkleTreeVoucherReq) (*lindapb.MerkleTreeVoucherInfo, error) {
	return c.client.GetMerkleTreeVoucherInfo(ctx, req)
}

func (c *FullNodeClient) ScanNoteByIvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	return c.client.ScanNoteByIvk(ctx, req)
}

func (c *FullNodeClient) ScanAndMarkNoteByIvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	return c.client.ScanAndMarkNoteByIvk(ctx, req)
}

func (c *FullNodeClient) ScanNoteByOvk(ctx context.Context, req *lindapb.ScanNoteReq) (*lindapb.ScanNoteResponse, error) {
	return c.client.ScanNoteByOvk(ctx, req)
}

func (c *FullNodeClient) IsSpend(ctx context.Context, req *lindapb.IsSpendReq) (*lindapb.IsSpendResponse, error) {
	return c.client.IsSpend(ctx, req)
}

func (c *FullNodeClient) GetRcm(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.BytesMessage, error) {
	return c.client.GetRcm(ctx, req)
}

func (c *FullNodeClient) CreateSpendAuthSig(ctx context.Context, req *lindapb.SpendAuthSigReq) (*lindapb.BytesMessage, error) {
	return c.client.CreateSpendAuthSig(ctx, req)
}

func (c *FullNodeClient) CreateShieldNullifier(ctx context.Context, req *lindapb.ShieldNullifierReq) (*lindapb.BytesMessage, error) {
	return c.client.CreateShieldNullifier(ctx, req)
}

func (c *FullNodeClient) GetShieldTransactionHash(ctx context.Context, req *lindapb.Transaction) (*lindapb.BytesMessage, error) {
	return c.client.GetShieldTransactionHash(ctx, req)
}

func (c *FullNodeClient) GetShieldedLRC20NotesByIvk(ctx context.Context, req *lindapb.ShieldedLRC20NoteReq) (*lindapb.ShieldedLRC20NoteResponse, error) {
	return c.client.GetShieldedLRC20NotesByIvk(ctx, req)
}

func (c *FullNodeClient) GetShieldedLRC20NotesByOvk(ctx context.Context, req *lindapb.ShieldedLRC20NoteReq) (*lindapb.ShieldedLRC20NoteResponse, error) {
	return c.client.GetShieldedLRC20NotesByOvk(ctx, req)
}

func (c *FullNodeClient) IsShieldedLRC20ContractNoteSpent(ctx context.Context, req *lindapb.IsShieldedLRC20NoteSpentReq) (*lindapb.IsShieldedLRC20NoteSpentResponse, error) {
	return c.client.IsShieldedLRC20ContractNoteSpent(ctx, req)
}

func (c *FullNodeClient) CreateShieldedContractParameters(ctx context.Context, req *lindapb.ShieldedContractParametersReq) (*lindapb.ShieldedContractParameters, error) {
	return c.client.CreateShieldedContractParameters(ctx, req)
}

func (c *FullNodeClient) CreateShieldedContractParametersWithoutAsk(ctx context.Context, req *lindapb.ShieldedContractParametersReq) (*lindapb.ShieldedContractParameters, error) {
	return c.client.CreateShieldedContractParametersWithoutAsk(ctx, req)
}

func (c *FullNodeClient) GetTriggerInputForShieldedLRC20Contract(ctx context.Context, req *lindapb.ShieldedLRC20TriggerInputReq) (*lindapb.ShieldedLRC20TriggerInput, error) {
	return c.client.GetTriggerInputForShieldedLRC20Contract(ctx, req)
}

// Utility methods
func (c *FullNodeClient) ValidateAddress(ctx context.Context, req *lindapb.AddressMessage) (*lindapb.AddressValidateResponse, error) {
	return c.client.ValidateAddress(ctx, req)
}

func (c *FullNodeClient) GenerateAddress(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.AddressPrKeyPairMessage, error) {
	return c.client.GenerateAddress(ctx, req)
}

func (c *FullNodeClient) CreateAddress(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.BytesMessage, error) {
	return c.client.CreateAddress(ctx, req)
}

func (c *FullNodeClient) EasyTransfer(ctx context.Context, req *lindapb.EasyTransferMessage) (*lindapb.EasyTransferResponse, error) {
	return c.client.EasyTransfer(ctx, req)
}

func (c *FullNodeClient) EasyTransferByPrivate(ctx context.Context, req *lindapb.EasyTransferByPrivateMessage) (*lindapb.EasyTransferResponse, error) {
	return c.client.EasyTransferByPrivate(ctx, req)
}

func (c *FullNodeClient) EasyTransferAsset(ctx context.Context, req *lindapb.EasyTransferAssetMessage) (*lindapb.EasyTransferResponse, error) {
	return c.client.EasyTransferAsset(ctx, req)
}

func (c *FullNodeClient) EasyTransferAssetByPrivate(ctx context.Context, req *lindapb.EasyTransferAssetByPrivateMessage) (*lindapb.EasyTransferResponse, error) {
	return c.client.EasyTransferAssetByPrivate(ctx, req)
}

// Pending Pool methods
func (c *FullNodeClient) GetTransactionFromPending(ctx context.Context, req *lindapb.BytesMessage) (*lindapb.Transaction, error) {
	return c.client.GetTransactionFromPending(ctx, req)
}

func (c *FullNodeClient) GetTransactionListFromPending(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.TransactionIdList, error) {
	return c.client.GetTransactionListFromPending(ctx, req)
}

func (c *FullNodeClient) GetPendingSize(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	return c.client.GetPendingSize(ctx, req)
}

// Chain info methods
func (c *FullNodeClient) GetChainParameters(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.ChainParameters, error) {
	return c.client.GetChainParameters(ctx, req)
}

func (c *FullNodeClient) TotalTransaction(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	return c.client.TotalTransaction(ctx, req)
}

func (c *FullNodeClient) GetBurnLind(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.NumberMessage, error) {
	return c.client.GetBurnLind(ctx, req)
}

func (c *FullNodeClient) GetEnergyPrices(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.PriceHistory, error) {
	return c.client.GetEnergyPrices(ctx, req)
}

func (c *FullNodeClient) GetBandwidthPrices(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.PriceHistory, error) {
	return c.client.GetBandwidthPrices(ctx, req)
}

func (c *FullNodeClient) GetMemoFeePrices(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.PriceHistory, error) {
	return c.client.GetMemoFeePrices(ctx, req)
}

// Metrics methods
func (c *FullNodeClient) GetStatsInfo(ctx context.Context, req *lindapb.EmptyMessage) (*lindapb.StatsInfo, error) {
	return c.client.GetStatsInfo(ctx, req)
}