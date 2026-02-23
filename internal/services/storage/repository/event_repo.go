package repository

import (
	"encoding/json"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"gorm.io/gorm"
)

// EventRepository struct: Repository for event operations
type EventRepository struct {
	db *gorm.DB
}

// NewEventRepository function: Creates a new event repository
func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

// SaveEvent function: Saves an event
func (r *EventRepository) SaveEvent(event *models.EventResponse) error {
		// Marshal the maps to JSON
		resultJSON, err := json.Marshal(event.Result)
		if err != nil {
			return err
		}
		
		resultTypeJSON, err := json.Marshal(event.ResultType)
		if err != nil {
			return err
		}
    // Convert to Event model for storage
    eventModel := &models.Event{
        BlockNumber:     event.BlockNumber,
        BlockTimestamp:  event.BlockTimestamp,
        ContractAddress: event.ContractAddress,
        EventIndex:      event.EventIndex,
        EventName:       event.EventName,
        EventSignature:  event.Event,
        TransactionID:   event.TransactionID,
        Result:          models.JSON(resultJSON),
        ResultType:      models.JSON(resultTypeJSON),
        Unconfirmed:     event.Unconfirmed,
    }
    return r.db.Save(eventModel).Error
}

// GetEvents function: Retrieves events with filters
func (r *EventRepository) GetEvents(contractAddress, eventName, transactionID string, blockNumber int64, fromTimestamp, toTimestamp int64, offset, limit int, sort string, confirmed bool) ([]*models.EventResponse, int64, error) {
	var events []*models.EventResponse
	var total int64

	query := r.db.Model(&models.EventResponse{})

	if contractAddress != "" {
		query = query.Where("contract_address = ?", contractAddress)
	}
	if eventName != "" {
		query = query.Where("event_name = ?", eventName)
	}
	if transactionID != "" {
		query = query.Where("transaction_id = ?", transactionID)
	}
	if blockNumber > 0 {
		query = query.Where("block_number = ?", blockNumber)
	}
	if fromTimestamp > 0 {
		query = query.Where("block_timestamp >= ?", fromTimestamp)
	}
	if toTimestamp > 0 {
		query = query.Where("block_timestamp <= ?", toTimestamp)
	}
	if confirmed {
		query = query.Where("unconfirmed = ?", false)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sort != "" {
		if sort[0] == '-' {
			query = query.Order(sort[1:] + " DESC")
		} else {
			query = query.Order(sort + " ASC")
		}
	} else {
		query = query.Order("block_timestamp DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetEventsByTransactionId function: Retrieves events for a transaction
func (r *EventRepository) GetEventsByTransactionId(transactionID string, offset, limit int) ([]*models.EventResponse, int64, error) {
	return r.GetEvents("", "", transactionID, 0, 0, 0, offset, limit, "", false)
}

// GetEventsByContractAddress function: Retrieves events for a contract
func (r *EventRepository) GetEventsByContractAddress(contractAddress string, eventName string, fromBlock int64, fromTimestamp, toTimestamp int64, offset, limit int, sort string) ([]*models.EventResponse, int64, error) {
	return r.GetEvents(contractAddress, eventName, "", fromBlock, fromTimestamp, toTimestamp, offset, limit, sort, false)
}