package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type StatsHandler struct {
	blockchainClient *blockchain.Client
	statsRepo        *repository.StatsRepository
}

func NewStatsHandler(client *blockchain.Client, statsRepo *repository.StatsRepository) *StatsHandler {
	return &StatsHandler{
		blockchainClient: client,
		statsRepo:        statsRepo,
	}
}

// ==================== Homepage Bundle ====================

// GetHomepageBundle handles GET /api/system/homepage-bundle
func (h *StatsHandler) GetHomepageBundle(c *gin.Context) {
	// Get from database or calculate
	bundle, err := h.statsRepo.GetHomepageBundle()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get homepage bundle: "+err.Error())
		return
	}

	// Get market data
	price, marketCap, volume := h.getMarketData()
	bundle.PriceUSD = price
	bundle.MarketCap = marketCap
	bundle.Volume24h = volume

	utils.RespondWithSuccess(c, bundle)
}

// ==================== Top 10 ====================

// GetTop10 handles GET /api/top10
func (h *StatsHandler) GetTop10(c *gin.Context) {
	var req models.Top10Request
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get top witnesses
	witnesses, err := h.blockchainClient.ListWitnesses(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get witnesses: "+err.Error())
		return
	}

	// Get top accounts from database
	accounts, err := h.statsRepo.GetTopAccounts(10)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get accounts: "+err.Error())
		return
	}

	// Get top tokens from database
	tokens, err := h.statsRepo.GetTopTokens(10)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get tokens: "+err.Error())
		return
	}

	response := &models.Top10Response{
		Witnesses: make([]models.WitnessResponse, 0),
		Accounts:  accounts,
		Tokens:    tokens,
	}

	// Convert witnesses
	for i, w := range witnesses.Witnesses {
		if i >= 10 {
			break
		}
		response.Witnesses = append(response.Witnesses, models.WitnessResponse{
			Address:   utils.MustHexToBase58(string(w.Address)),
			VoteCount: w.VoteCount,
			URL:       w.Url,
		})
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Overview Statistics ====================

// GetOverview handles GET /api/stats/overview
func (h *StatsHandler) GetOverview(c *gin.Context) {
	statType := c.DefaultQuery("type", "")

	data, err := h.statsRepo.GetOverview(statType)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get overview: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, data)
}

// ==================== Energy Statistics ====================

// GetEnergyStatistic handles GET /api/energystatistic
func (h *StatsHandler) GetEnergyStatistic(c *gin.Context) {
	var req models.EnergyStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Address == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Address is required")
		return
	}

	stats, err := h.statsRepo.GetEnergyStatistic(req.Address, req.From, req.To)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get energy statistic: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, stats)
}

// GetEnergyDailyStatistic handles GET /api/energydailystatistic
func (h *StatsHandler) GetEnergyDailyStatistic(c *gin.Context) {
	var req models.EnergyDailyStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	stats, err := h.statsRepo.GetEnergyDailyStatistic(req.From, req.To)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get energy daily statistic: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, stats)
}

// ==================== Trigger Statistics ====================

// GetTriggerStatistic handles GET /api/triggerstatistic
func (h *StatsHandler) GetTriggerStatistic(c *gin.Context) {
	var req models.TriggerStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	stats, err := h.statsRepo.GetTriggerStatistic(req.Contract, req.From, req.To)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get trigger statistic: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, stats)
}

// GetTriggerAmountStatistic handles GET /api/triggeramountstatistic
func (h *StatsHandler) GetTriggerAmountStatistic(c *gin.Context) {
	var req models.TriggerAmountStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	// Get from database
	// This would be implemented with the stats repository

	utils.RespondWithSuccess(c, gin.H{
		"contract": req.Contract,
		"triggers": []interface{}{},
	})
}

// ==================== Caller Statistics ====================

// GetCallerAddressStatistic handles GET /api/calleraddressstatistic
func (h *StatsHandler) GetCallerAddressStatistic(c *gin.Context) {
	var req models.CallerAddressStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	stats, err := h.statsRepo.GetCallerAddressStatistic(req.Contract, req.From, req.To)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get caller statistic: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, stats)
}

// ==================== One Contract Statistics ====================

// GetOneContractEnergyStatistic handles GET /api/onecontractenergystatistic
func (h *StatsHandler) GetOneContractEnergyStatistic(c *gin.Context) {
	var req models.OneContractEnergyStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	// Get from database
	utils.RespondWithSuccess(c, gin.H{
		"contract": req.Contract,
		"total":    0,
		"daily":    []interface{}{},
	})
}

// GetOneContractTriggerStatistic handles GET /api/onecontracttriggerstatistic
func (h *StatsHandler) GetOneContractTriggerStatistic(c *gin.Context) {
	var req models.OneContractTriggerStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	// Get from database
	utils.RespondWithSuccess(c, gin.H{
		"contract": req.Contract,
		"count":    0,
		"daily":    []interface{}{},
	})
}

// GetOneContractCallerStatistic handles GET /api/onecontractcallerstatistic
func (h *StatsHandler) GetOneContractCallerStatistic(c *gin.Context) {
	var req models.OneContractCallerStatisticRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	// Get from database
	utils.RespondWithSuccess(c, gin.H{
		"contract": req.Contract,
		"callers":  []interface{}{},
	})
}

// GetOneContractCallers handles GET /api/onecontractcallers
func (h *StatsHandler) GetOneContractCallers(c *gin.Context) {
	var req models.OneContractCallersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	// Get from database
	utils.RespondWithSuccess(c, gin.H{
		"callers": []interface{}{},
		"total":   0,
	})
}

// ==================== Freeze Resource ====================

// GetFreezeResource handles GET /api/freezeresource
func (h *StatsHandler) GetFreezeResource(c *gin.Context) {
	var req models.FreezeResourceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Address == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Address is required")
		return
	}

	data, err := h.statsRepo.GetFreezeResource(req.Address, req.Type)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get freeze resource: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, data)
}

// ExportFreezeResourceToCSV handles CSV export for freeze resource
func (h *StatsHandler) ExportFreezeResourceToCSV(c *gin.Context, req models.CSVExportRequest) {
	// Get data
	data, err := h.statsRepo.GetFreezeResource(req.Address, "")
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get data: "+err.Error())
		return
	}

	headers := []string{"Type", "Amount", "Expire At", "From", "To"}
	var rows [][]string

	for _, f := range data.Frozen {
		rows = append(rows, []string{
			f.Type,
			strconv.FormatInt(f.Amount, 10),
			strconv.FormatInt(f.ExpireAt, 10),
			"",
			"",
		})
	}

	for _, d := range data.Delegated {
		rows = append(rows, []string{
			d.Type,
			strconv.FormatInt(d.Amount, 10),
			strconv.FormatInt(d.ExpireAt, 10),
			d.From,
			d.To,
		})
	}

	utils.RespondWithCSV(c, "freezeresource.csv", headers, rows)
}

// ==================== Turnover ====================

// GetTurnover handles GET /api/turnover
func (h *StatsHandler) GetTurnover(c *gin.Context) {
	var req models.TurnoverRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	data, err := h.statsRepo.GetTurnover(req.From, req.To)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get turnover: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, data)
}

// ExportTurnoverToCSV handles CSV export for turnover
func (h *StatsHandler) ExportTurnoverToCSV(c *gin.Context, req models.CSVExportRequest) {
	data, err := h.statsRepo.GetTurnover(req.From, req.To)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get data: "+err.Error())
		return
	}

	headers := []string{"Date", "Turnover"}
	var rows [][]string

	for _, d := range data.Daily {
		rows = append(rows, []string{
			d.Date,
			strconv.FormatInt(d.Turnover, 10),
		})
	}

	utils.RespondWithCSV(c, "turnover.csv", headers, rows)
}

// ==================== Fund ====================

// GetFund handles GET /api/fund
func (h *StatsHandler) GetFund(c *gin.Context) {
	var req models.FundRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	// Get from database
	utils.RespondWithSuccess(c, gin.H{
		"fund":  []interface{}{},
		"total": 0,
	})
}

// ==================== Ledger ====================

// GetLedger handles GET /api/ledger
func (h *StatsHandler) GetLedger(c *gin.Context) {
	var req models.LedgerRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	// Get from database based on type
	utils.RespondWithSuccess(c, gin.H{
		"ledger": []interface{}{},
		"total":  0,
	})
}

// ==================== Helper Functions ====================

func (h *StatsHandler) getMarketData() (price, marketCap, volume int64) {
	// Get from external API or database
	// This would be implemented with CoinMarketCap API
	return 0, 0, 0
}