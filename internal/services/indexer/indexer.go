package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/sirupsen/logrus"
)

type Indexer struct {
	config          *config.IndexerConfig
	blockchainClient *blockchain.Client
	blockRepo       *repository.BlockRepository
	txRepo          *repository.TransactionRepository
	tokenRepo       *repository.TokenRepository
	eventRepo       *repository.EventRepository
	statsRepo       *repository.StatsRepository
	
	logger          *logrus.Logger
	stopChan        chan struct{}
	wg              sync.WaitGroup
	currentBlock    int64
}

func NewIndexer(
	cfg *config.IndexerConfig,
	client *blockchain.Client,
	blockRepo *repository.BlockRepository,
	txRepo *repository.TransactionRepository,
	tokenRepo *repository.TokenRepository,
	eventRepo *repository.EventRepository,
	statsRepo *repository.StatsRepository,
) *Indexer {
	return &Indexer{
		config:          cfg,
		blockchainClient: client,
		blockRepo:       blockRepo,
		txRepo:          txRepo,
		tokenRepo:       tokenRepo,
		eventRepo:       eventRepo,
		statsRepo:       statsRepo,
		logger:          logrus.New(),
		stopChan:        make(chan struct{}),
		currentBlock:    0,
	}
}

func (i *Indexer) Start() error {
	i.logger.Info("Starting blockchain indexer")
	
	// Get last indexed block
	lastBlock, err := i.blockRepo.GetLastIndexedBlock()
	if err == nil {
		i.currentBlock = lastBlock
	} else {
		i.currentBlock = i.config.StartBlock
	}

	// Start workers
	for w := 0; w < i.config.MaxWorkers; w++ {
		i.wg.Add(1)
		go i.worker(w)
	}

	// Start sync ticker
	ticker := time.NewTicker(i.config.SyncInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				i.sync()
			case <-i.stopChan:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (i *Indexer) Stop() error {
	i.logger.Info("Stopping blockchain indexer")
	close(i.stopChan)
	i.wg.Wait()
	return nil
}

func (i *Indexer) worker(id int) {
	defer i.wg.Done()
	i.logger.WithField("worker", id).Info("Indexer worker started")

	for {
		select {
		case <-i.stopChan:
			i.logger.WithField("worker", id).Info("Indexer worker stopped")
			return
		default:
			// Workers pick up tasks from a queue
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (i *Indexer) sync() {
	ctx := context.Background()

	// Get latest block
	nowBlock, err := i.blockchainClient.GetNowBlock(ctx, nil)
	if err != nil {
		i.logger.WithError(err).Error("Failed to get now block")
		return
	}

	latestBlock := nowBlock.BlockHeader.RawData.Number
	i.logger.WithFields(logrus.Fields{
		"current": i.currentBlock,
		"latest":  latestBlock,
	}).Info("Syncing blocks")

	// Sync blocks in batches
	for i.currentBlock < latestBlock {
		select {
		case <-i.stopChan:
			return
		default:
			endBlock := i.currentBlock + int64(i.config.BlockBatchSize)
			if endBlock > latestBlock {
				endBlock = latestBlock
			}

			if err := i.syncBlockRange(ctx, i.currentBlock+1, endBlock); err != nil {
				i.logger.WithError(err).Error("Failed to sync block range")
				time.Sleep(5 * time.Second)
				break
			}

			i.currentBlock = endBlock
		}
	}
}

func (i *Indexer) syncBlockRange(ctx context.Context, start, end int64) error {
	i.logger.WithFields(logrus.Fields{
		"start": start,
		"end":   end,
	}).Info("Syncing block range")

	for blockNum := start; blockNum <= end; blockNum++ {
		if err := i.syncBlock(ctx, blockNum); err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) syncBlock(ctx context.Context, blockNum int64) error {
	// Get block from full node
	block, err := i.blockchainClient.GetBlockByNum(ctx, &lindapb.NumberMessage{
		Num: blockNum,
	})
	if err != nil {
		return err
	}

	// Get block info from solidity node (confirmed)
	blockSolidity, err := i.blockchainClient.GetBlockByNumSolidity(ctx, &lindapb.NumberMessage{
		Num: blockNum,
	})
	if err != nil {
		// Not confirmed yet, use full node data
		blockSolidity = block
	}

	// Index block
	if err := i.indexBlock(blockSolidity); err != nil {
		return err
	}

	// Index transactions
	for _, tx := range blockSolidity.Transactions {
		if err := i.indexTransaction(ctx, tx, blockNum); err != nil {
			i.logger.WithError(err).WithField("tx", string(tx.TxID)).Error("Failed to index transaction")
		}
	}

	// Get transaction infos
	txInfos, err := i.blockchainClient.GetTransactionInfoByBlockNumSolidity(ctx, &lindapb.NumberMessage{
		Num: blockNum,
	})
	if err == nil {
		for _, info := range txInfos.TransactionInfo {
			if err := i.indexTransactionInfo(info); err != nil {
				i.logger.WithError(err).Error("Failed to index transaction info")
			}
		}
	}

	return nil
}

func (i *Indexer) indexBlock(block *lindapb.Block) error {
	blockModel := &models.Block{
		Number:           block.BlockHeader.RawData.Number,
		Hash:             string(block.BlockID),
		ParentHash:       string(block.BlockHeader.RawData.ParentHash),
		Timestamp:        block.BlockHeader.RawData.Timestamp,
		WitnessAddress:   string(block.BlockHeader.RawData.WitnessAddress),
		WitnessID:        int(block.BlockHeader.RawData.WitnessId),
		TxTrieRoot:       string(block.BlockHeader.RawData.TxTrieRoot),
		TransactionCount: len(block.Transactions),
		Size:             0, // Would need actual size
		Version:          int(block.BlockHeader.RawData.Version),
	}

	return i.blockRepo.SaveBlock(blockModel)
}

func (i *Indexer) indexTransaction(ctx context.Context, tx *lindapb.Transaction, blockNum int64) error {
	txModel := &models.Transaction{
		Hash:          string(tx.TxID),
		BlockNumber:   blockNum,
		RawData:       string(tx.RawDataHex),
		Signature:     tx.Signature,
	}

	// Parse from address and to address
	if len(tx.RawData.Contract) > 0 {
		// Parse contract based on type
		// This is simplified - actual implementation would handle different contract types
	}

	return i.txRepo.SaveTransaction(txModel)
}

func (i *Indexer) indexTransactionInfo(info *lindapb.TransactionInfo) error {
	// Update transaction with info
	return i.txRepo.UpdateTransactionWithInfo(info)
}