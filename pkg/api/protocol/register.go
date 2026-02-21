package protocol

import (
    "google.golang.org/protobuf/reflect/protoregistry"
)

// RegisterAllTypes ensures all protocol types are registered
func RegisterAllTypes() error {
    // The import of this package already registers all types
    // This function just verifies they're registered
    types := []string{
        "protocol.ProposalCreateContract",
        "protocol.ProposalApproveContract",
        "protocol.TransferContract",
        "protocol.TransferAssetContract",
        "protocol.VoteWitnessContract",
        "protocol.WitnessCreateContract",
        "protocol.WitnessUpdateContract",
        "protocol.AssetIssueContract",
        "protocol.ParticipateAssetIssueContract",
        "protocol.AccountCreateContract",
        "protocol.AccountUpdateContract",
        "protocol.FreezeBalanceContract",
        "protocol.UnfreezeBalanceContract",
        "protocol.WithdrawBalanceContract",
        "protocol.UnfreezeAssetContract",
        "protocol.UpdateAssetContract",
        "protocol.ProposalCreateContract",
        "protocol.ProposalApproveContract",
        "protocol.ProposalDeleteContract",
        "protocol.ExchangeCreateContract",
        "protocol.ExchangeInjectContract",
        "protocol.ExchangeWithdrawContract",
        "protocol.ExchangeTransactionContract",
        "protocol.MarketSellAssetContract",
        "protocol.MarketCancelOrderContract",
    }
    
    for _, typeName := range types {
        _, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(typeName))
        if err != nil {
            return fmt.Errorf("type %s not registered: %v", typeName, err)
        }
    }
    return nil
}