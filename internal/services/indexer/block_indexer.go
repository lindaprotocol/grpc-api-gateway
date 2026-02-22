package indexer

import (
	"context"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/models"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
)

type BlockIndexer struct {
	indexer *Indexer
}

func NewBlockIndexer(indexer *Indexer) *BlockIndexer {
	return &BlockIndexer{
		indexer: indexer,
	}
}

// IndexBlock indexes a single block
func (bi *BlockIndexer) IndexBlock(block *lindapb.Block) error {
	blockModel := &models.Block{
		Number:           block.BlockHeader.RawData.Number,
		Hash:             string(block.BlockID),
		ParentHash:       string(block.BlockHeader.RawData.ParentHash),
		Timestamp:        block.BlockHeader.RawData.Timestamp,
		WitnessAddress:   string(block.BlockHeader.RawData.WitnessAddress),
		WitnessID:        int(block.BlockHeader.RawData.WitnessId),
		TxTrieRoot:       string(block.BlockHeader.RawData.TxTrieRoot),
		TransactionCount: len(block.Transactions),
		Size:             calculateBlockSize(block),
		Version:          int(block.BlockHeader.RawData.Version),
		CreatedAt:        time.Now(),
	}

	// Convert witness address to base58 for storage
	witnessBase58, err := utils.HexToBase58(blockModel.WitnessAddress)
	if err == nil {
		blockModel.WitnessAddress = witnessBase58
	}

	return bi.indexer.blockRepo.SaveBlock(blockModel)
}

// IndexBlocksBatch indexes a batch of blocks
func (bi *BlockIndexer) IndexBlocksBatch(ctx context.Context, startNum, endNum int64) error {
	for blockNum := startNum; blockNum <= endNum; blockNum++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Get block from full node
			block, err := bi.indexer.blockchainClient.GetBlockByNum(ctx, &lindapb.NumberMessage{
				Num: blockNum,
			})
			if err != nil {
				bi.indexer.logger.WithError(err).WithField("block", blockNum).Error("Failed to get block")
				continue
			}

			// Index block
			if err := bi.IndexBlock(block); err != nil {
				bi.indexer.logger.WithError(err).WithField("block", blockNum).Error("Failed to index block")
			}
		}
	}
	return nil
}

// GetLatestBlock retrieves the latest block
func (bi *BlockIndexer) GetLatestBlock(ctx context.Context) (*lindapb.Block, error) {
	return bi.indexer.blockchainClient.GetNowBlock(ctx, &lindapb.EmptyMessage{})
}

// GetBlockRange retrieves a range of blocks from the database
func (bi *BlockIndexer) GetBlockRange(start, end int64) ([]*models.Block, error) {
	return bi.indexer.blockRepo.GetBlockRange(start, end)
}

// CalculateBlockStats calculates statistics for a block
func (bi *BlockIndexer) CalculateBlockStats(block *lindapb.Block) *models.BlockStatsResponse {
	stats := &models.BlockStatsResponse{
		FeeStat: &models.FeeStat{},
	}

	if len(block.Transactions) > 0 {
		stats.TxStat = &models.TxStat{
			ContractTypeDistribute: make(map[int]int),
		}
	}

	for _, tx := range block.Transactions {
		// Count transaction types
		if len(tx.RawData.Contract) > 0 {
			contractType := int(tx.RawData.Contract[0].Type)
			stats.TxStat.ContractTypeDistribute[contractType]++
		}

		// Calculate fees
		for _, ret := range tx.Ret {
			stats.FeeStat.OtherFee += ret.Fee
		}

		// Check if failed
		if len(tx.Ret) > 0 && tx.Ret[0].Ret == lindapb.Transaction_Result_FAILED {
			stats.TxStat.FailTxCount++
		}
	}

	return stats
}

// calculateBlockSize calculates the approximate size of a block
func calculateBlockSize(block *lindapb.Block) int {
	// This is a simplified calculation
	size := 0
	size += len(block.BlockID)
	if block.BlockHeader != nil {
		size += len(block.BlockHeader.WitnessSignature)
		if block.BlockHeader.RawData != nil {
			size += len(block.BlockHeader.RawData.TxTrieRoot)
			size += len(block.BlockHeader.RawData.ParentHash)
			size += len(block.BlockHeader.RawData.WitnessAddress)
			size += len(block.BlockHeader.RawData.AccountStateRoot)
		}
	}
	for _, tx := range block.Transactions {
		size += len(tx.TxID)
		size += len(tx.RawDataHex)
		size += len(tx.Signature) * 65 // Approximate signature size
	}
	return size
}