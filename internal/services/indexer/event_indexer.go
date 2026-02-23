package indexer

import (
	"context"
	"encoding/hex"
	"math/big"

	
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/event"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type EventIndexer struct {
	indexer *Indexer
	parser  *event.EventParser
}

func NewEventIndexer(indexer *Indexer) *EventIndexer {
	return &EventIndexer{
		indexer: indexer,
		parser:  event.NewEventParser(),
	}
}

// IndexEvents indexes events from a transaction
func (ei *EventIndexer) IndexEvents(ctx context.Context, tx *lindapb.Transaction, txInfo *lindapb.TransactionInfo, block *lindapb.Block) error {
	if txInfo == nil || len(txInfo.Log) == 0 {
		return nil
	}

	for i, log := range txInfo.Log {
		event := &models.EventResponse{
			BlockNumber:          txInfo.BlockNumber,
			BlockTimestamp:       txInfo.BlockTimeStamp,
			ContractAddress:      utils.MustHexToBase58(string(log.Address)),
			EventIndex:           string(rune(i)),
			TransactionID:        string(txInfo.Id),
			Result:               make(map[string]interface{}),
			ResultType:           make(map[string]string),
			Unconfirmed:          false,
		}

		// Parse event name and parameters
		if len(log.Topics) > 0 {
			eventName, params, paramTypes := ei.parser.ParseEvent(log.Topics, log.Data)
			event.EventName = eventName
			event.Event = eventName
			event.Result = params
			event.ResultType = paramTypes
		}

		// Save to repository
		if err := ei.indexer.eventRepo.SaveEvent(event); err != nil {
			return err
		}

		// If this is a token transfer, also index as transfer
		if event.EventName == "Transfer" {
			if err := ei.indexer.tokenIndexer.IndexTokenTransfer(ctx, event); err != nil {
				ei.indexer.logger.WithError(err).Error("Failed to index token transfer")
			}
		}
	}

	return nil
}

// IndexEventsFromBlock indexes all events in a block
func (ei *EventIndexer) IndexEventsFromBlock(ctx context.Context, block *lindapb.Block) error {
	// Get transaction infos for the block
	txInfos, err := ei.indexer.blockchainClient.GetTransactionInfoByBlockNum(ctx, &lindapb.NumberMessage{
		Num: block.BlockHeader.RawData.Number,
	})
	if err != nil {
		return err
	}

	// Map transactions by ID for quick lookup
	txMap := make(map[string]*lindapb.Transaction)
	for _, tx := range block.Transactions {
		txMap[string(tx.TxID)] = tx
	}

	// Index each transaction's events
	for _, info := range txInfos.TransactionInfo {
		tx, ok := txMap[string(info.Id)]
		if !ok {
			continue
		}
		if err := ei.IndexEvents(ctx, tx, info, block); err != nil {
			ei.indexer.logger.WithError(err).WithField("tx", string(info.Id)).Error("Failed to index events")
		}
	}

	return nil
}

// GetEventsByFilter retrieves events based on filters
func (ei *EventIndexer) GetEventsByFilter(filter *event.EventFilter) ([]*models.EventResponse, int64, error) {
	return ei.indexer.eventRepo.GetEvents(
		filter.ContractAddress,
		filter.EventName,
		filter.TransactionID,
		filter.BlockNumber,
		filter.FromTimestamp,
		filter.ToTimestamp,
		filter.Offset,
		filter.Limit,
		filter.Sort,
		filter.Confirmed,
	)
}

// ParseEventLog parses a raw event log into a structured event
func (ei *EventIndexer) ParseEventLog(log *lindapb.TransactionInfo_Log, blockNumber, blockTimestamp int64, txID string) (*models.EventResponse, error) {
	event := &models.EventResponse{
		BlockNumber:     blockNumber,
		BlockTimestamp:  blockTimestamp,
		ContractAddress: utils.MustHexToBase58(string(log.Address)),
		TransactionID:   txID,
		Result:          make(map[string]interface{}),
		ResultType:      make(map[string]string),
		Unconfirmed:     false,
	}

	// Parse event name from first topic
	if len(log.Topics) > 0 {
		eventName, params, paramTypes := ei.parser.ParseEvent(log.Topics, log.Data)
		event.EventName = eventName
		event.Event = eventName
		event.Result = params
		event.ResultType = paramTypes
	}

	return event, nil
}

// ParseTransferEvent parses a transfer event specifically
func (ei *EventIndexer) ParseTransferEvent(event *models.EventResponse) (from, to, value string, err error) {
	if event.EventName != "Transfer" {
		return "", "", "", nil
	}

	fromVal, ok := event.Result["from"]
	if !ok {
		return "", "", "", nil
	}
	from, ok = fromVal.(string)
	if !ok {
		return "", "", "", nil
	}

	toVal, ok := event.Result["to"]
	if !ok {
		return "", "", "", nil
	}
	to, ok = toVal.(string)
	if !ok {
		return "", "", "", nil
	}

	valueVal, ok := event.Result["value"]
	if !ok {
		return "", "", "", nil
	}
	value, ok = valueVal.(string)
	if !ok {
		return "", "", "", nil
	}

	return from, to, value, nil
}

// IsERC20Transfer checks if an event is an ERC20 transfer
func (ei *EventIndexer) IsERC20Transfer(event *models.EventResponse) bool {
	if event.EventName != "Transfer" {
		return false
	}
	_, hasFrom := event.Result["from"]
	_, hasTo := event.Result["to"]
	_, hasValue := event.Result["value"]
	return hasFrom && hasTo && hasValue
}

// GetTokenTransferValue returns the transfer value as *big.Int
func (ei *EventIndexer) GetTokenTransferValue(event *models.EventResponse) (*big.Int, error) {
	if !ei.IsERC20Transfer(event) {
		return nil, nil
	}
	valueStr, ok := event.Result["value"].(string)
	if !ok {
		return nil, nil
	}
	return new(big.Int).SetString(valueStr, 10)
}