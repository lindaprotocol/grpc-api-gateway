package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/event"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type EventHandler struct {
	blockchainClient *blockchain.Client
	eventRepo        *repository.EventRepository
	eventService     *event.EventService
}

func NewEventHandler(client *blockchain.Client, eventRepo *repository.EventRepository) *EventHandler {
	return &EventHandler{
		blockchainClient: client,
		eventRepo:        eventRepo,
		eventService:     event.NewEventService(eventRepo),
	}
}

// ==================== Event Query Service Methods ====================

// GetEvents handles GET /v1/events
func (h *EventHandler) GetEvents(c *gin.Context) {
	limit, start, sort, fingerprint := utils.ParseV1PaginationParams(c)

	// Parse additional filters
	blockNumber, _ := strconv.ParseInt(c.DefaultQuery("block", "0"), 10, 64)
	since, _ := strconv.ParseInt(c.DefaultQuery("since", "0"), 10, 64)
	to, _ := strconv.ParseInt(c.DefaultQuery("to", "0"), 10, 64)
	contract := c.Query("contract")
	eventName := c.Query("event_name")

	filter := &event.EventFilter{
		ContractAddress: contract,
		EventName:       eventName,
		BlockNumber:     blockNumber,
		FromTimestamp:   since,
		ToTimestamp:     to,
		Offset:          start,
		Limit:           limit,
		Sort:            sort,
		Confirmed:       true,
	}

	events, total, err := h.eventService.GetEvents(context.Background(), filter)
	if err != nil {
		utils.RespondWithV1Error(c, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	// Convert to interface slice
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	utils.RespondWithV1Success(c, data, limit, "")
}

// GetEventsByTransactionId handles GET /v1/events/transaction/{transaction_id}
func (h *EventHandler) GetEventsByTransactionId(c *gin.Context) {
	txID := c.Param("transaction_id")
	if txID == "" {
		utils.RespondWithV1Error(c, http.StatusBadRequest, "Transaction ID is required")
		return
	}

	limit, start, _, _ := utils.ParseV1PaginationParams(c)

	events, total, err := h.eventService.GetEventsByTransactionID(context.Background(), txID, start, limit)
	if err != nil {
		utils.RespondWithV1Error(c, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	// Convert to interface slice
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	utils.RespondWithV1Success(c, data, limit, "")
}

// GetEventsByContractAddress handles GET /v1/events/{contract_address}
func (h *EventHandler) GetEventsByContractAddress(c *gin.Context) {
	contractAddr := c.Param("contract_address")
	if contractAddr == "" {
		utils.RespondWithV1Error(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	limit, start, sort, fingerprint := utils.ParseV1PaginationParams(c)
	eventName := c.DefaultQuery("event_name", "")
	blockNumber, _ := strconv.ParseInt(c.DefaultQuery("block", "0"), 10, 64)
	since, _ := strconv.ParseInt(c.DefaultQuery("since", "0"), 10, 64)

	events, total, err := h.eventService.GetEventsByContractAddress(
		context.Background(),
		contractAddr,
		eventName,
		blockNumber,
		since,
		0,
		start,
		limit,
		sort,
	)
	if err != nil {
		utils.RespondWithV1Error(c, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	// Convert to interface slice
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	utils.RespondWithV1Success(c, data, limit, fingerprint)
}

// GetEventsByContractAndEventName handles GET /v1/events/contract/{contract_address}/{event_name}
func (h *EventHandler) GetEventsByContractAndEventName(c *gin.Context) {
	contractAddr := c.Param("contract_address")
	eventName := c.Param("event_name")

	if contractAddr == "" || eventName == "" {
		utils.RespondWithV1Error(c, http.StatusBadRequest, "Contract address and event name are required")
		return
	}

	limit, start, sort, fingerprint := utils.ParseV1PaginationParams(c)
	since, _ := strconv.ParseInt(c.DefaultQuery("since", "0"), 10, 64)

	events, total, err := h.eventService.GetEventsByContractAddress(
		context.Background(),
		contractAddr,
		eventName,
		0,
		since,
		0,
		start,
		limit,
		sort,
	)
	if err != nil {
		utils.RespondWithV1Error(c, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	// Convert to interface slice
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	utils.RespondWithV1Success(c, data, limit, fingerprint)
}

// GetEventsByContractEventAndBlock handles GET /v1/events/contract/{contract_address}/{event_name}/{block_number}
func (h *EventHandler) GetEventsByContractEventAndBlock(c *gin.Context) {
	contractAddr := c.Param("contract_address")
	eventName := c.Param("event_name")
	blockNumberStr := c.Param("block_number")

	if contractAddr == "" || eventName == "" || blockNumberStr == "" {
		utils.RespondWithV1Error(c, http.StatusBadRequest, "Contract address, event name, and block number are required")
		return
	}

	blockNumber, err := strconv.ParseInt(blockNumberStr, 10, 64)
	if err != nil {
		utils.RespondWithV1Error(c, http.StatusBadRequest, "Invalid block number")
		return
	}

	events, total, err := h.eventService.GetEventsByContractAddress(
		context.Background(),
		contractAddr,
		eventName,
		blockNumber,
		0,
		0,
		0,
		100,
		"",
	)
	if err != nil {
		utils.RespondWithV1Error(c, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	// Convert to interface slice
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	utils.RespondWithV1Success(c, data, 100, "")
}

// GetEventsByTimestamp handles GET /v1/events/timestamp
func (h *EventHandler) GetEventsByTimestamp(c *gin.Context) {
	limit, start, sort, fingerprint := utils.ParseV1PaginationParams(c)
	since, _ := strconv.ParseInt(c.DefaultQuery("since", "0"), 10, 64)
	contract := c.Query("contract")

	filter := &event.EventFilter{
		ContractAddress: contract,
		FromTimestamp:   since,
		Offset:          start,
		Limit:           limit,
		Sort:            sort,
		Confirmed:       true,
	}

	events, total, err := h.eventService.GetEvents(context.Background(), filter)
	if err != nil {
		utils.RespondWithV1Error(c, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	// Convert to interface slice
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	utils.RespondWithV1Success(c, data, limit, fingerprint)
}

// GetConfirmedEvents handles GET /v1/events/confirmed
func (h *EventHandler) GetConfirmedEvents(c *gin.Context) {
	limit, start, sort, fingerprint := utils.ParseV1PaginationParams(c)
	since, _ := strconv.ParseInt(c.DefaultQuery("since", "0"), 10, 64)

	filter := &event.EventFilter{
		FromTimestamp: since,
		Offset:        start,
		Limit:         limit,
		Sort:          sort,
		Confirmed:     true,
	}

	events, total, err := h.eventService.GetEvents(context.Background(), filter)
	if err != nil {
		utils.RespondWithV1Error(c, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	// Convert to interface slice
	data := make([]interface{}, len(events))
	for i, e := range events {
		data[i] = e
	}

	utils.RespondWithV1Success(c, data, limit, fingerprint)
}

// ==================== Additional Event Methods for Lindascan ====================

// GetContractEvents handles GET /api/contract/events
func (h *EventHandler) GetContractEvents(c *gin.Context) {
	var req struct {
		Contract string `form:"contract" binding:"required"`
		From     int64  `form:"from"`
		To       int64  `form:"to"`
		Limit    int    `form:"limit" default:"20"`
		Start    int    `form:"start" default:"0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	filter := &event.EventFilter{
		ContractAddress: req.Contract,
		FromTimestamp:   req.From,
		ToTimestamp:     req.To,
		Offset:          req.Start,
		Limit:           req.Limit,
		Sort:            "-timestamp",
		Confirmed:       true,
	}

	events, total, err := h.eventService.GetEvents(context.Background(), filter)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	utils.RespondWithPagination(c, events, total, req.Start/req.Limit+1, req.Limit)
}

// ProcessTransactionEvents processes events from a transaction (called by indexer)
func (h *EventHandler) ProcessTransactionEvents(tx *lindapb.Transaction, txInfo *lindapb.TransactionInfo, block *lindapb.Block) error {
	return h.eventService.ProcessTransactionEvents(tx, txInfo, block)
}