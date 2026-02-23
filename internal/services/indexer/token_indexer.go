// internal/services/indexer/token_indexer.go
package indexer

import (
	"context"
	"math/big"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

// TokenIndexer struct: Indexer for token operations
type TokenIndexer struct {
	indexer *Indexer
}

// NewTokenIndexer creates a new token indexer
func NewTokenIndexer(indexer *Indexer) *TokenIndexer {
	return &TokenIndexer{
		indexer: indexer,
	}
}

// IndexLRC20Token indexes or updates an LRC20 token
func (ti *TokenIndexer) IndexLRC20Token(ctx context.Context, contractAddr string) error {
	// Get contract info
	contract, err := ti.indexer.blockchainClient.GetContract(ctx, &lindapb.BytesMessage{
		Value: []byte(contractAddr),
	})
	if err != nil {
		return err
	}

	token := &models.LRC20TokenInfo{
		Contract:    utils.MustHexToBase58(contractAddr),
		Name:        "Unknown",
		Symbol:      "UNKNOWN",
		Decimals:    18,
		TotalSupply: "0",
		Owner:       utils.MustHexToBase58(string(contract.OriginAddress)),
		IssueTime:   time.Now().Unix(),
		Holders:     0,
		Transfers:   0,
	}

	return ti.indexer.tokenRepo.SaveLRC20Token(token)
}

// IndexTokenTransfer indexes a token transfer event
func (ti *TokenIndexer) IndexTokenTransfer(ctx context.Context, event *models.EventResponse) error {
	// Parse transfer event
	if event.EventName != "Transfer" {
		return nil
	}

	from, ok := event.Result["from"].(string)
	if !ok {
		return nil
	}
	to, ok := event.Result["to"].(string)
	if !ok {
		return nil
	}
	value, ok := event.Result["value"].(string)
	if !ok {
		return nil
	}

	// Update token holders
	if err := ti.updateTokenHolder(event.ContractAddress, from, to, value); err != nil {
		ti.indexer.logger.WithError(err).Error("Failed to update token holders")
		return err
	}

	// Save transfer to repository
	transfer := &models.TokenTransferResponse{
		TransactionID:  event.TransactionID,
		BlockNumber:    event.BlockNumber,
		BlockTimestamp: event.BlockTimestamp,
		From:           utils.MustHexToBase58(from),
		To:             utils.MustHexToBase58(to),
		Value:          value,
		TokenAddress:   utils.MustHexToBase58(event.ContractAddress),
		TokenSymbol:    "",
		TokenDecimals:  18,
	}
	
	// Note: You'll need to implement SaveTokenTransfer in tokenRepo
	return ti.indexer.tokenRepo.SaveTokenTransfer(transfer)
}

// updateTokenHolder updates token holder balances
func (ti *TokenIndexer) updateTokenHolder(contractAddr, from, to, value string) error {
	valueBig, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil
	}

	// Decrease from balance
	if from != "0x0000000000000000000000000000000000000000" {
		fromBase58 := utils.MustHexToBase58(from)
		if err := ti.indexer.tokenRepo.UpdateHolderBalance(contractAddr, fromBase58, new(big.Int).Neg(valueBig)); err != nil {
			return err
		}
	}

	// Increase to balance
	if to != "0x0000000000000000000000000000000000000000" {
		toBase58 := utils.MustHexToBase58(to)
		if err := ti.indexer.tokenRepo.UpdateHolderBalance(contractAddr, toBase58, valueBig); err != nil {
			return err
		}
	}

	return nil
}

// IndexLRC10Token indexes a LRC10 token
func (ti *TokenIndexer) IndexLRC10Token(ctx context.Context, tokenID string) error {
    asset, err := ti.indexer.blockchainClient.GetAssetIssueById(ctx, &lindapb.BytesMessage{
        Value: []byte(tokenID),
    })
    if err != nil {
        return err
    }

    // Convert to token model
    token := &models.TokenInfo{
        ID:          tokenID, // Use the tokenID parameter, not asset.Id
        Name:        string(asset.Name),
        Symbol:      string(asset.Abbr),
        TotalSupply: asset.TotalSupply,
        Owner:       utils.MustHexToBase58(string(asset.OwnerAddress)),
        Decimals:    int(asset.Precision),
        StartTime:   asset.StartTime,
        EndTime:     asset.EndTime,
        URL:         string(asset.Url),
        Description: string(asset.Description),
    }

    return ti.indexer.tokenRepo.SaveLRC10Token(token)
}

// UpdateTokenHolderCounts updates holder counts for all tokens
func (ti *TokenIndexer) UpdateTokenHolderCounts(ctx context.Context) error {
	tokens, _, err := ti.indexer.tokenRepo.GetLRC20Tokens(0, 1000, "")
	if err != nil {
		return err
	}

	for _, token := range tokens {
		count, err := ti.indexer.tokenRepo.GetHolderCount(token.Contract)
		if err != nil {
			continue
		}
		if err := ti.indexer.tokenRepo.UpdateHolderCount(token.Contract, count); err != nil {
			ti.indexer.logger.WithError(err).WithField("token", token.Contract).Error("Failed to update holder count")
		}
	}

	return nil
}

// CalculateTokenPercentages calculates and updates holder percentages
func (ti *TokenIndexer) CalculateTokenPercentages(ctx context.Context, contractAddr string) error {
	holders, _, err := ti.indexer.tokenRepo.GetTokenHolders(contractAddr, 0, 1000, "-balance")
	if err != nil {
		return err
	}

	token, err := ti.indexer.tokenRepo.GetLRC20TokenByContract(contractAddr)
	if err != nil {
		return err
	}

	totalSupply, ok := new(big.Int).SetString(token.TotalSupply, 10)
	if !ok || totalSupply.Sign() == 0 {
		return nil
	}

	for _, holder := range holders {
		balance, ok := new(big.Int).SetString(holder.Balance, 10)
		if !ok {
			continue
		}
		percentage := new(big.Float).Quo(
			new(big.Float).SetInt(balance),
			new(big.Float).SetInt(totalSupply),
		)
		percentageFloat, _ := percentage.Float64()
		holder.Percentage = percentageFloat * 100

		if err := ti.indexer.tokenRepo.UpdateHolderPercentage(contractAddr, holder.Address, holder.Percentage); err != nil {
			ti.indexer.logger.WithError(err).Error("Failed to update holder percentage")
		}
	}

	return nil
}