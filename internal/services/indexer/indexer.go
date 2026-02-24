// internal/services/indexer/indexer.go
package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/sirupsen/logrus"
)

// Indexer struct: Main indexer service
type Indexer struct {
	config           *config.IndexerConfig
	blockchainClient *blockchain.Client
	blockRepo        *repository.BlockRepository
	txRepo           *repository.TransactionRepository
	tokenRepo        *repository.TokenRepository
	eventRepo        *repository.EventRepository
	statsRepo        *repository.StatsRepository
	
	logger           *logrus.Logger
	stopChan         chan struct{}
	wg               sync.WaitGroup
	currentBlock     int64
	
	// Indexer components
	blockIndexer     *BlockIndexer
	txIndexer        *TransactionIndexer
	tokenIndexer     *TokenIndexer
	eventIndexer     *EventIndexer
}

// NewIndexer creates a new indexer instance
func NewIndexer(
	cfg *config.IndexerConfig,
	client *blockchain.Client,
	blockRepo *repository.BlockRepository,
	txRepo *repository.TransactionRepository,
	tokenRepo *repository.TokenRepository,
	eventRepo *repository.EventRepository,
	statsRepo *repository.StatsRepository,
) *Indexer {
	idx := &Indexer{
		config:           cfg,
		blockchainClient: client,
		blockRepo:        blockRepo,
		txRepo:           txRepo,
		tokenRepo:        tokenRepo,
		eventRepo:        eventRepo,
		statsRepo:        statsRepo,
		logger:           logrus.New(),
		stopChan:         make(chan struct{}),
		currentBlock:     0,
	}
	
	// Initialize indexers
	idx.blockIndexer = NewBlockIndexer(idx)
	idx.txIndexer = NewTransactionIndexer(idx)
	idx.tokenIndexer = NewTokenIndexer(idx)
	idx.eventIndexer = NewEventIndexer(idx)
	
	return idx
}

// Start begins the indexing process
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

// Stop halts the indexing process
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
	nowBlock, err := i.blockchainClient.GetNowBlock(ctx, &lindapb.EmptyMessage{})
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

// syncBlock fetches and indexes a single block
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
	if err := i.blockIndexer.IndexBlock(blockSolidity); err != nil {
		return err
	}

	// Get block timestamp
	blockTimestamp := blockSolidity.BlockHeader.RawData.Timestamp

	// Index transactions
	for _, tx := range blockSolidity.Transactions {
		// FIX: Add blockTimestamp as the 4th argument
		if err := i.txIndexer.IndexTransaction(ctx, tx, blockNum, blockTimestamp); err != nil {
			i.logger.WithError(err).WithField("tx", string(tx.TxID)).Error("Failed to index transaction")
		}
	}

	// Get transaction infos
	txInfos, err := i.blockchainClient.GetTransactionInfoByBlockNumSolidity(ctx, &lindapb.NumberMessage{
		Num: blockNum,
	})
	if err == nil {
		for _, info := range txInfos.TransactionInfo {
			if err := i.txIndexer.IndexTransactionInfo(info); err != nil {
				i.logger.WithError(err).Error("Failed to index transaction info")
			}
			
			// Index events from transaction info
			if err := i.eventIndexer.IndexEvents(ctx, nil, info, blockSolidity); err != nil {
				i.logger.WithError(err).Error("Failed to index events")
			}
		}
	}

	return nil
}