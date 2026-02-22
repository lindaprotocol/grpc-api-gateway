package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/services/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
)

type EventService struct {
	eventRepo *repository.EventRepository
	parser    *EventParser
}

func NewEventService(eventRepo *repository.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		parser:    NewEventParser(),
	}
}

// ProcessTransactionEvents extracts and stores events from a transaction
func (s *EventService) ProcessTransactionEvents(tx *lindapb.Transaction, txInfo *lindapb.TransactionInfo, block *lindapb.Block) error {
	if txInfo == nil || len(txInfo.Log) == 0 {
		return nil
	}

	for i, log := range txInfo.Log {
		event := &models.EventResponse{
			BlockNumber:          txInfo.BlockNumber,
			BlockTimestamp:       txInfo.BlockTimeStamp,
			ContractAddress:      string(log.Address),
			EventIndex:           string(rune(i)),
			TransactionID:        string(txInfo.Id),
			Result:               make(map[string]interface{}),
			ResultType:           make(map[string]string),
			Unconfirmed:          false,
		}

		// Parse event name and parameters
		if len(log.Topics) > 0 {
			eventName, params, paramTypes := s.parser.ParseEvent(log.Topics, log.Data)
			event.EventName = eventName
			event.Event = eventName
			event.Result = params
			event.ResultType = paramTypes
		}

		// Save to repository
		if err := s.eventRepo.SaveEvent(event); err != nil {
			return err
		}
	}

	return nil
}

// GetEvents retrieves events based on filters
func (s *EventService) GetEvents(ctx context.Context, filter *EventFilter) ([]*models.EventResponse, int64, error) {
	return s.eventRepo.GetEvents(
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

// GetEventsByTransactionID retrieves events for a transaction
func (s *EventService) GetEventsByTransactionID(ctx context.Context, txID string, offset, limit int) ([]*models.EventResponse, int64, error) {
	return s.eventRepo.GetEventsByTransactionId(txID, offset, limit)
}

// GetEventsByContractAddress retrieves events for a contract
func (s *EventService) GetEventsByContractAddress(ctx context.Context, contractAddress string, eventName string, fromBlock int64, fromTimestamp, toTimestamp int64, offset, limit int, sort string) ([]*models.EventResponse, int64, error) {
	return s.eventRepo.GetEventsByContractAddress(contractAddress, eventName, fromBlock, fromTimestamp, toTimestamp, offset, limit, sort)
}

type EventFilter struct {
	ContractAddress string
	EventName       string
	TransactionID   string
	BlockNumber     int64
	FromTimestamp   int64
	ToTimestamp     int64
	Offset          int
	Limit           int
	Sort            string
	Confirmed       bool
}