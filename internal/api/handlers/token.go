package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type TokenHandler struct {
	blockchainClient *blockchain.Client
	tokenRepo        *repository.TokenRepository
}

func NewTokenHandler(client *blockchain.Client, tokenRepo *repository.TokenRepository) *TokenHandler {
	return &TokenHandler{
		blockchainClient: client,
		tokenRepo:        tokenRepo,
	}
}

// ==================== Wallet Service Asset Methods (LRC-10) ====================

// GetAssetIssueByAccount handles POST /wallet/getassetissuebyaccount
// Returns LRC-10 tokens issued by an account
func (h *TokenHandler) GetAssetIssueByAccount(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	address := req.Address
	if req.Visible {
		hexAddr, err := utils.Base58ToHex(address)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid address format")
			return
		}
		address = hexAddr
	}

	assets, err := h.blockchainClient.GetAssetIssueByAccount(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get assets: "+err.Error())
		return
	}

	response := convertAssetIssueListToResponse(assets, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetAssetIssueById handles POST /wallet/getassetissuebyid
// Returns LRC-10 token by ID
func (h *TokenHandler) GetAssetIssueById(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	asset, err := h.blockchainClient.GetAssetIssueById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get asset: "+err.Error())
		return
	}

	response := convertAssetIssueToResponse(asset, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetAssetIssueByName handles POST /wallet/getassetissuebyname
// Returns LRC-10 token by name
func (h *TokenHandler) GetAssetIssueByName(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	asset, err := h.blockchainClient.GetAssetIssueByName(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get asset: "+err.Error())
		return
	}

	response := convertAssetIssueToResponse(asset, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetAssetIssueList handles GET /wallet/getassetissuelist
// Returns all LRC-10 tokens
func (h *TokenHandler) GetAssetIssueList(c *gin.Context) {
	visible := c.DefaultQuery("visible", "false") == "true"

	assets, err := h.blockchainClient.GetAssetIssueList(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get assets: "+err.Error())
		return
	}

	response := convertAssetIssueListToResponse(assets, visible)
	utils.RespondWithSuccess(c, response)
}

// GetAssetIssueListByName handles POST /wallet/getassetissuelistbyname
// Returns all LRC-10 tokens with given name
func (h *TokenHandler) GetAssetIssueListByName(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	assets, err := h.blockchainClient.GetAssetIssueListByName(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get assets: "+err.Error())
		return
	}

	response := convertAssetIssueListToResponse(assets, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetPaginatedAssetIssueList handles POST /wallet/getpaginatedassetissuelist
// Returns paginated LRC-10 tokens
func (h *TokenHandler) GetPaginatedAssetIssueList(c *gin.Context) {
	var req struct {
		Offset  int64 `json:"offset" default:"0"`
		Limit   int64 `json:"limit" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	assets, err := h.blockchainClient.GetPaginatedAssetIssueList(context.Background(), &lindapb.PaginatedMessage{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get assets: "+err.Error())
		return
	}

	response := convertAssetIssueListToResponse(assets, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// CreateAssetIssue handles POST /wallet/createassetissue
// Creates a new LRC-10 token
func (h *TokenHandler) CreateAssetIssue(c *gin.Context) {
	var req struct {
		OwnerAddress             string          `json:"owner_address" binding:"required"`
		Name                     string          `json:"name" binding:"required"`
		Abbr                     string          `json:"abbr" binding:"required"`
		TotalSupply              int64           `json:"total_supply" binding:"required"`
		LindNum                   int32           `json:"lind_num" binding:"required"`
		Num                      int32           `json:"num" binding:"required"`
		StartTime                int64           `json:"start_time" binding:"required"`
		EndTime                  int64           `json:"end_time" binding:"required"`
		Precision                int32           `json:"precision"`
		Description              string          `json:"description"`
		URL                      string          `json:"url"`
		FreeAssetNetLimit        int64           `json:"free_asset_net_limit"`
		PublicFreeAssetNetLimit  int64           `json:"public_free_asset_net_limit"`
		FrozenSupply             json.RawMessage `json:"frozen_supply"`
		Visible                  bool            `json:"visible" default:"false"`
		PermissionID             int32           `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	if req.Visible {
		hexAddr, err := utils.Base58ToHex(ownerAddress)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid address format")
			return
		}
		ownerAddress = hexAddr
	}

	// Build contract
	contract := &lindapb.AssetIssueContract{
		OwnerAddress:             []byte(ownerAddress),
		Name:                     []byte(req.Name),
		Abbr:                     []byte(req.Abbr),
		TotalSupply:              req.TotalSupply,
		LindNum:                   req.LindNum,
		Num:                      req.Num,
		StartTime:                req.StartTime,
		EndTime:                  req.EndTime,
		Precision:                req.Precision,
		Description:              []byte(req.Description),
		Url:                      []byte(req.URL),
		FreeAssetNetLimit:        req.FreeAssetNetLimit,
		PublicFreeAssetNetLimit:  req.PublicFreeAssetNetLimit,
	}

	// Parse frozen supply if present
	if len(req.FrozenSupply) > 0 {
		var frozenSupplies []struct {
			FrozenAmount int64 `json:"frozen_amount"`
			FrozenDays   int64 `json:"frozen_days"`
		}
		if err := json.Unmarshal(req.FrozenSupply, &frozenSupplies); err == nil {
			for _, fs := range frozenSupplies {
				contract.FrozenSupply = append(contract.FrozenSupply, &lindapb.AssetIssueContract_FrozenSupply{
					FrozenAmount: fs.FrozenAmount,
					FrozenDays:   fs.FrozenDays,
				})
			}
		}
	}

	tx, err := h.blockchainClient.CreateAssetIssue(context.Background(), contract)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create asset: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// TransferAsset handles POST /wallet/transferasset
// Transfers LRC-10 token
func (h *TokenHandler) TransferAsset(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		ToAddress    string `json:"to_address" binding:"required"`
		AssetName    string `json:"asset_name" binding:"required"`
		Amount       int64  `json:"amount" binding:"required"`
		Visible      bool   `json:"visible" default:"false"`
		PermissionID int32  `json:"permission_id"`
		ExtraData    string `json:"extra_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	toAddress := req.ToAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexTo, _ := utils.Base58ToHex(toAddress)
		ownerAddress = hexOwner
		toAddress = hexTo
	}

	tx, err := h.blockchainClient.TransferAsset(context.Background(), &lindapb.TransferAssetContract{
		OwnerAddress: []byte(ownerAddress),
		ToAddress:    []byte(toAddress),
		AssetName:    []byte(req.AssetName),
		Amount:       req.Amount,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to transfer asset: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ParticipateAssetIssue handles POST /wallet/participateassetissue
// Participates in token issuance
func (h *TokenHandler) ParticipateAssetIssue(c *gin.Context) {
	var req struct {
		ToAddress    string `json:"to_address" binding:"required"`
		OwnerAddress string `json:"owner_address" binding:"required"`
		Amount       int64  `json:"amount" binding:"required"`
		AssetName    string `json:"asset_name" binding:"required"`
		Visible      bool   `json:"visible" default:"false"`
		PermissionID int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	toAddress := req.ToAddress
	ownerAddress := req.OwnerAddress

	if req.Visible {
		hexTo, _ := utils.Base58ToHex(toAddress)
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		toAddress = hexTo
		ownerAddress = hexOwner
	}

	tx, err := h.blockchainClient.ParticipateAssetIssue(context.Background(), &lindapb.ParticipateAssetIssueContract{
		ToAddress:    []byte(toAddress),
		OwnerAddress: []byte(ownerAddress),
		AssetName:    []byte(req.AssetName),
		Amount:       req.Amount,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to participate: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// UnfreezeAsset handles POST /wallet/unfreezeasset
// Unfreezes frozen LRC-10 token
func (h *TokenHandler) UnfreezeAsset(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		Visible      bool   `json:"visible" default:"false"`
		PermissionID int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	if req.Visible {
		hexAddr, _ := utils.Base58ToHex(ownerAddress)
		ownerAddress = hexAddr
	}

	tx, err := h.blockchainClient.UnfreezeAsset(context.Background(), &lindapb.UnfreezeAssetContract{
		OwnerAddress: []byte(ownerAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to unfreeze asset: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// UpdateAsset handles POST /wallet/updateasset
// Updates LRC-10 token information
func (h *TokenHandler) UpdateAsset(c *gin.Context) {
	var req struct {
		OwnerAddress    string `json:"owner_address" binding:"required"`
		Description     string `json:"description"`
		URL             string `json:"url"`
		NewLimit        int32  `json:"new_limit"`
		NewPublicLimit  int32  `json:"new_public_limit"`
		Visible         bool   `json:"visible" default:"false"`
		PermissionID    int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	if req.Visible {
		hexAddr, _ := utils.Base58ToHex(ownerAddress)
		ownerAddress = hexAddr
	}

	tx, err := h.blockchainClient.UpdateAsset(context.Background(), &lindapb.UpdateAssetContract{
		OwnerAddress:    []byte(ownerAddress),
		Description:     []byte(req.Description),
		Url:             []byte(req.URL),
		NewLimit:        req.NewLimit,
		NewPublicLimit:  req.NewPublicLimit,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update asset: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ==================== Solidity Node Asset Methods ====================

// GetAssetIssueByIdSolidity handles POST /walletsolidity/getassetissuebyid
// Returns confirmed LRC-10 token by ID
func (h *TokenHandler) GetAssetIssueByIdSolidity(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	id, err := strconv.ParseInt(req.Value, 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid token ID")
		return
	}

	asset, err := h.blockchainClient.GetAssetIssueByIdSolidity(context.Background(), &lindapb.NumberMessage{
		Num: id,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get asset: "+err.Error())
		return
	}

	response := convertAssetIssueToResponse(asset, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetAssetIssueByNameSolidity handles POST /walletsolidity/getassetissuebyname
// Returns confirmed LRC-10 token by name
func (h *TokenHandler) GetAssetIssueByNameSolidity(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	asset, err := h.blockchainClient.GetAssetIssueByNameSolidity(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get asset: "+err.Error())
		return
	}

	response := convertAssetIssueToResponse(asset, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetAssetIssueListSolidity handles GET /walletsolidity/getassetissuelist
// Returns all confirmed LRC-10 tokens
func (h *TokenHandler) GetAssetIssueListSolidity(c *gin.Context) {
	visible := c.DefaultQuery("visible", "false") == "true"

	assets, err := h.blockchainClient.GetAssetIssueListSolidity(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get assets: "+err.Error())
		return
	}

	response := convertAssetIssueListToResponse(assets, visible)
	utils.RespondWithSuccess(c, response)
}

// ==================== LRC20 Token Methods ====================

// GetLRC20Tokens handles GET /api/token_lrc20
// Returns paginated LRC20 tokens
func (h *TokenHandler) GetLRC20Tokens(c *gin.Context) {
	var req models.LRC20TokenRequest
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

	var tokens []*models.LRC20TokenInfo
	var total int64
	var err error

	if req.Contract != "" {
		// Get single token
		token, err := h.tokenRepo.GetLRC20TokenByContract(req.Contract)
		if err == nil && token != nil {
			tokens = append(tokens, token)
			total = 1
		}
	} else {
		// Get paginated list
		tokens, total, err = h.tokenRepo.GetLRC20Tokens(req.Start, req.Limit, req.Sort)
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get tokens: "+err.Error())
		return
	}

	response := &models.LRC20TokenListResponse{
		Tokens: tokens,
		Total:  total,
	}

	utils.RespondWithSuccess(c, response)
}

// GetLRC20TokenByContract handles GET /api/token_lrc20/{contract}
// Returns LRC20 token by contract address
func (h *TokenHandler) GetLRC20TokenByContract(c *gin.Context) {
	contract := c.Param("contract")
	if contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	token, err := h.tokenRepo.GetLRC20TokenByContract(contract)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Token not found")
		return
	}

	utils.RespondWithSuccess(c, token)
}

// GetTokenHolders handles GET /api/tokenholders
// Returns token holders with pagination
func (h *TokenHandler) GetTokenHolders(c *gin.Context) {
	var req models.TokenHoldersRequest
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

	holders, total, err := h.tokenRepo.GetTokenHolders(req.Contract, req.Start, req.Limit, req.Sort)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get holders: "+err.Error())
		return
	}

	// Handle CSV export
	if req.Format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment;filename=holders.csv")

		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		// Write header
		writer.Write([]string{"Address", "Balance", "Percentage", "Rank"})

		// Write data
		for _, holder := range holders {
			writer.Write([]string{
				holder.Address,
				holder.Balance,
				strconv.FormatFloat(holder.Percentage, 'f', 4, 64),
				strconv.FormatInt(holder.Rank, 10),
			})
		}
		return
	}

	response := &models.TokenHoldersResponse{
		Holders: holders,
		Total:   total,
	}

	utils.RespondWithSuccess(c, response)
}

// GetTokenTransfers handles GET /api/token_lrc20/transfers
// Returns token transfers
func (h *TokenHandler) GetTokenTransfers(c *gin.Context) {
	var req models.TokenTransfersRequest
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

	transfers, total, err := h.tokenRepo.GetTokenTransfers(req.Contract, req.From, req.To, req.Start, req.Limit, req.Sort)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transfers: "+err.Error())
		return
	}

	response := &models.TokenTransfersResponse{
		Transfers: transfers,
		Total:     total,
	}

	utils.RespondWithSuccess(c, response)
}

// GetTokens handles GET /api/token
// Returns tokens (both LRC-10 and LRC20)
func (h *TokenHandler) GetTokens(c *gin.Context) {
	var req models.TokenListRequest
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

	var response gin.H

	if req.ID != "" {
		// Get LRC-10 by ID
		asset, err := h.blockchainClient.GetAssetIssueById(context.Background(), &lindapb.BytesMessage{
			Value: []byte(req.ID),
		})
		if err == nil {
			response = gin.H{
				"tokens": []interface{}{convertAssetIssueToResponse(asset, true)},
				"total":  1,
			}
		}
	} else if req.Contract != "" {
		// Get LRC20 by contract
		token, err := h.tokenRepo.GetLRC20TokenByContract(req.Contract)
		if err == nil {
			response = gin.H{
				"tokens": []interface{}{token},
				"total":  1,
			}
		}
	} else {
		// Get paginated based on filter
		if req.Filter == "lrc20" {
			tokens, total, _ := h.tokenRepo.GetLRC20Tokens(req.Start, req.Limit, req.Sort)
			response = gin.H{
				"tokens": tokens,
				"total":  total,
			}
		} else if req.Filter == "lrc10" {
			// Get LRC-10 from blockchain (simplified)
			assets, _ := h.blockchainClient.GetAssetIssueList(context.Background(), &lindapb.EmptyMessage{})
			var tokens []interface{}
			for i, asset := range assets.AssetIssue {
				if i >= req.Start && i < req.Start+req.Limit {
					tokens = append(tokens, convertAssetIssueToResponse(asset, true))
				}
			}
			response = gin.H{
				"tokens": tokens,
				"total":  int64(len(assets.AssetIssue)),
			}
		}
	}

	if response == nil {
		response = gin.H{"tokens": []interface{}{}, "total": 0}
	}

	utils.RespondWithSuccess(c, response)
}

// GetTokensOverview handles GET /api/tokens/overview
// Returns token overview statistics
func (h *TokenHandler) GetTokensOverview(c *gin.Context) {
	var req struct {
		UUID   string `form:"uuid"`
		Start  int    `form:"start" default:"0"`
		Limit  int    `form:"limit" default:"20"`
		Filter string `form:"filter"` // lrc20
		Sort   string `form:"sort"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	// Get token statistics from database
	stats, total, err := h.tokenRepo.GetTokensOverview(req.Filter, req.Start, req.Limit, req.Sort)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get overview: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{
		"data":  stats,
		"total": total,
	})
}

// GetTokenPrice handles GET /api/token/price
// Returns token price information
func (h *TokenHandler) GetTokenPrice(c *gin.Context) {
	// Get from cache or external API
	price, err := h.tokenRepo.GetTokenPrice()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get price: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, price)
}

// GetParticipateAssetIssue handles GET /api/tokens/participateassetissue
// Returns participation information
func (h *TokenHandler) GetParticipateAssetIssue(c *gin.Context) {
	var req models.ParticipateAssetIssueRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	// Get from database
	participations, total, err := h.tokenRepo.GetParticipations(req.Start, req.Limit)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get participations: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{
		"participations": participations,
		"total":          total,
	})
}

// GetTokenPositionDistribution handles GET /api/tokens/position-distribution
// Returns token position distribution
func (h *TokenHandler) GetTokenPositionDistribution(c *gin.Context) {
	var req models.TokenPositionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Contract == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	if req.Limit <= 0 {
		req.Limit = 100
	}

	positions, err := h.tokenRepo.GetTokenPositionDistribution(req.Contract, req.Limit)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get positions: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, positions)
}

// GetWinkFund handles GET /api/wink/fund
// Returns WINK fund information
func (h *TokenHandler) GetWinkFund(c *gin.Context) {
	fund, err := h.tokenRepo.GetWinkFund()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get WINK fund: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, fund)
}

// GetWinkGraphic handles GET /api/wink/graphic
// Returns WINK graphic data
func (h *TokenHandler) GetWinkGraphic(c *gin.Context) {
	data, err := h.tokenRepo.GetWinkGraphic()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get WINK graphic: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, data)
}

// GetJSTFund handles GET /api/jst/fund
// Returns JST fund information
func (h *TokenHandler) GetJSTFund(c *gin.Context) {
	fund, err := h.tokenRepo.GetJSTFund()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get JST fund: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, fund)
}

// GetJSTGraphic handles GET /api/jst/graphic
// Returns JST graphic data
func (h *TokenHandler) GetJSTGraphic(c *gin.Context) {
	data, err := h.tokenRepo.GetJSTGraphic()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get JST graphic: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, data)
}

// GetBitTorrentGraphic handles GET /api/bittorrent/graphic
// Returns BitTorrent graphic data
func (h *TokenHandler) GetBitTorrentGraphic(c *gin.Context) {
	data, err := h.tokenRepo.GetBitTorrentGraphic()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get BitTorrent graphic: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, data)
}

// GetAssetTransfer handles GET /api/asset/transfer
// Returns asset transfers
func (h *TokenHandler) GetAssetTransfer(c *gin.Context) {
	var req models.AssetTransferRequest
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

	transfers, total, err := h.tokenRepo.GetAssetTransfers(req.AssetName, req.Start, req.Limit, req.Sort)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transfers: "+err.Error())
		return
	}

	response := &models.AssetTransferResponse{
		Transfers: transfers,
		Total:     total,
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Helper Functions ====================

func convertAssetIssueToResponse(asset *lindapb.AssetIssueContract, visible bool) *models.AssetIssueResponse {
	resp := &models.AssetIssueResponse{
		ID:                       string(asset.Id),
		OwnerAddress:             string(asset.OwnerAddress),
		Name:                     string(asset.Name),
		Abbr:                     string(asset.Abbr),
		TotalSupply:              asset.TotalSupply,
		LindNum:                   asset.LindNum,
		Num:                      asset.Num,
		Precision:                asset.Precision,
		StartTime:                asset.StartTime,
		EndTime:                  asset.EndTime,
		VoteScore:                asset.VoteScore,
		Description:              string(asset.Description),
		URL:                      string(asset.Url),
		FreeAssetNetLimit:        asset.FreeAssetNetLimit,
		PublicFreeAssetNetLimit:  asset.PublicFreeAssetNetLimit,
		PublicFreeAssetNetUsage:  asset.PublicFreeAssetNetUsage,
		PublicLatestFreeNetTime:  asset.PublicLatestFreeNetTime,
	}

	for _, fs := range asset.FrozenSupply {
		resp.FrozenSupply = append(resp.FrozenSupply, models.FrozenSupply{
			FrozenAmount: fs.FrozenAmount,
			FrozenDays:   fs.FrozenDays,
		})
	}

	if visible {
		base58Addr, _ := utils.HexToBase58(string(asset.OwnerAddress))
		resp.OwnerAddress = base58Addr
	}

	return resp
}

func convertAssetIssueListToResponse(assets *lindapb.AssetIssueList, visible bool) *models.AssetIssueListResponse {
	resp := &models.AssetIssueListResponse{}

	for _, asset := range assets.AssetIssue {
		resp.AssetIssue = append(resp.AssetIssue, *convertAssetIssueToResponse(asset, visible))
	}

	return resp
}