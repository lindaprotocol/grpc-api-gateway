package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type AccountHandler struct {
	blockchainClient *blockchain.Client
	accountRepo      *repository.AccountRepository
	tagRepo          *repository.TagRepository
}

func NewAccountHandler(client *blockchain.Client, accountRepo *repository.AccountRepository, tagRepo *repository.TagRepository) *AccountHandler {
	return &AccountHandler{
		blockchainClient: client,
		accountRepo:      accountRepo,
		tagRepo:          tagRepo,
	}
}

// ==================== Wallet Service Account Methods ====================

// GetAccount handles GET /wallet/getaccount
// Returns account information including balance, resources, permissions
func (h *AccountHandler) GetAccount(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Convert address if needed
	address := req.Address
	if !req.Visible {
		// Already in hex format
	} else {
		// Convert from base58 to hex
		hexAddr, err := utils.Base58ToHex(address)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid address format")
			return
		}
		address = hexAddr
	}

	// Call blockchain
	account, err := h.blockchainClient.GetAccount(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get account: "+err.Error())
		return
	}

	// Convert to response format
	response := convertAccountToResponse(account, req.Visible)

	// Get tags for this address
	tags, _ := h.tagRepo.GetTagsByAddress(address, 0, 100)
	if len(tags) > 0 {
		response.Tags = tags
	}

	utils.RespondWithSuccess(c, response)
}

// GetAccountBalance handles POST /wallet/getaccountbalance
// Returns historical balance at specific block
func (h *AccountHandler) GetAccountBalance(c *gin.Context) {
	var req struct {
		AccountIdentifier struct {
			Address string `json:"address"`
		} `json:"account_identifier" binding:"required"`
		BlockIdentifier struct {
			Hash   string `json:"hash"`
			Number int64  `json:"number"`
		} `json:"block_identifier" binding:"required"`
		Visible bool `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Convert address if needed
	address := req.AccountIdentifier.Address
	if req.Visible {
		hexAddr, err := utils.Base58ToHex(address)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid address format")
			return
		}
		address = hexAddr
	}

	// Call blockchain
	resp, err := h.blockchainClient.GetAccountBalance(context.Background(), &lindapb.AccountBalanceRequest{
		AccountIdentifier: &lindapb.AccountIdentifier{
			Address: []byte(address),
		},
		BlockIdentifier: &lindapb.BlockIdentifier{
			Hash:   []byte(req.BlockIdentifier.Hash),
			Number: req.BlockIdentifier.Number,
		},
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get account balance: "+err.Error())
		return
	}

	response := gin.H{
		"balance": resp.Balance,
		"block_identifier": gin.H{
			"hash":   string(resp.BlockIdentifier.Hash),
			"number": resp.BlockIdentifier.Number,
		},
	}

	utils.RespondWithSuccess(c, response)
}

// GetAccountResource handles POST /wallet/getaccountresource
// Returns resource information (bandwidth, energy, etc)
func (h *AccountHandler) GetAccountResource(c *gin.Context) {
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

	resp, err := h.blockchainClient.GetAccountResource(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get account resource: "+err.Error())
		return
	}

	response := &models.AccountResourceResponse{
		FreeNetUsed:     resp.FreeNetUsed,
		FreeNetLimit:    resp.FreeNetLimit,
		NetUsed:         resp.NetUsed,
		NetLimit:        resp.NetLimit,
		TotalNetLimit:   resp.TotalNetLimit,
		TotalNetWeight:  resp.TotalNetWeight,
		TotalLindaPowerWeight: resp.TotalLindaPowerWeight,
		LindaPowerLimit:  resp.LindaPowerLimit,
		LindaPowerUsed:   resp.LindaPowerUsed,
		EnergyUsed:      resp.EnergyUsed,
		EnergyLimit:     resp.EnergyLimit,
		TotalEnergyLimit: resp.TotalEnergyLimit,
		TotalEnergyWeight: resp.TotalEnergyWeight,
		AssetNetUsed:    resp.AssetNetUsed,
		AssetNetLimit:   resp.AssetNetLimit,
	}

	utils.RespondWithSuccess(c, response)
}

// GetAccountNet handles POST /wallet/getaccountnet
// Returns bandwidth information
func (h *AccountHandler) GetAccountNet(c *gin.Context) {
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

	resp, err := h.blockchainClient.GetAccountNet(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get account net: "+err.Error())
		return
	}

	response := gin.H{
		"freeNetUsed":     resp.FreeNetUsed,
		"freeNetLimit":    resp.FreeNetLimit,
		"NetUsed":         resp.NetUsed,
		"NetLimit":        resp.NetLimit,
		"TotalNetLimit":   resp.TotalNetLimit,
		"TotalNetWeight":  resp.TotalNetWeight,
		"assetNetUsed":    resp.AssetNetUsed,
		"assetNetLimit":   resp.AssetNetLimit,
	}

	utils.RespondWithSuccess(c, response)
}

// CreateAccount handles POST /wallet/createaccount
// Activates a new account
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req struct {
		OwnerAddress   string `json:"owner_address" binding:"required"`
		AccountAddress string `json:"account_address" binding:"required"`
		Visible        bool   `json:"visible" default:"false"`
		PermissionID   int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	accountAddress := req.AccountAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexAccount, _ := utils.Base58ToHex(accountAddress)
		ownerAddress = hexOwner
		accountAddress = hexAccount
	}

	tx, err := h.blockchainClient.CreateAccount(context.Background(), &lindapb.AccountCreateContract{
		OwnerAddress:   []byte(ownerAddress),
		AccountAddress: []byte(accountAddress),
		Type:           lindapb.AccountType_Normal,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create account: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// UpdateAccount handles POST /wallet/updateaccount
// Updates account name
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		AccountName  string `json:"account_name" binding:"required"`
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

	tx, err := h.blockchainClient.UpdateAccount(context.Background(), &lindapb.AccountUpdateContract{
		OwnerAddress:  []byte(ownerAddress),
		AccountName:   []byte(req.AccountName),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update account: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// AccountPermissionUpdate handles POST /wallet/accountpermissionupdate
// Updates account permissions
func (h *AccountHandler) AccountPermissionUpdate(c *gin.Context) {
	var req struct {
		OwnerAddress string          `json:"owner_address" binding:"required"`
		Owner        json.RawMessage `json:"owner"`
		Witness      json.RawMessage `json:"witness"`
		Actives      []json.RawMessage `json:"actives"`
		Visible      bool            `json:"visible" default:"false"`
		PermissionID int32           `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Parse permissions (simplified - full implementation would parse all fields)
	ownerAddress := req.OwnerAddress
	if req.Visible {
		hexAddr, _ := utils.Base58ToHex(ownerAddress)
		ownerAddress = hexAddr
	}

	// Build contract
	contract := &lindapb.AccountPermissionUpdateContract{
		OwnerAddress: []byte(ownerAddress),
		// Parse and set owner, witness, actives
	}

	tx, err := h.blockchainClient.AccountPermissionUpdate(context.Background(), contract)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update permissions: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ==================== Solidity Node Account Methods ====================

// GetAccountSolidity handles POST /walletsolidity/getaccount
// Returns confirmed account information
func (h *AccountHandler) GetAccountSolidity(c *gin.Context) {
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

	account, err := h.blockchainClient.GetAccountSolidity(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get account: "+err.Error())
		return
	}

	response := convertAccountToResponse(account, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetAccountByIdSolidity handles POST /walletsolidity/getaccountbyid
// Returns account by ID
func (h *AccountHandler) GetAccountByIdSolidity(c *gin.Context) {
	var req struct {
		AccountID string `json:"account_id" binding:"required"`
		Visible   bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	accountId := req.AccountID
	if req.Visible {
		hexAddr, err := utils.Base58ToHex(accountId)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid account ID format")
			return
		}
		accountId = hexAddr
	}

	account, err := h.blockchainClient.GetAccountByIdSolidity(context.Background(), &lindapb.Account{
		AccountId: []byte(accountId),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get account by ID: "+err.Error())
		return
	}

	response := convertAccountToResponse(account, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ==================== Lindascan Custom Account Methods ====================

// GetAccountList handles GET /api/account/list
// Returns paginated list of accounts
func (h *AccountHandler) GetAccountList(c *gin.Context) {
	var req models.AccountListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	var accounts []*models.AccountResponse
	var total int64
	var err error

	if req.Address != "" {
		// Get single account
		account, err := h.accountRepo.GetByAddress(req.Address)
		if err == nil && account != nil {
			accounts = append(accounts, account)
			total = 1
		}
	} else {
		// Get paginated list
		accounts, total, err = h.accountRepo.GetList(req.Start, req.Limit, req.Sort)
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get account list: "+err.Error())
		return
	}

	response := gin.H{
		"accounts": accounts,
		"total":    total,
	}

	utils.RespondWithSuccess(c, response)
}

// GetAccountResourceInfo handles GET /api/account/resource
// Returns resource information for an account
func (h *AccountHandler) GetAccountResourceInfo(c *gin.Context) {
	var req models.AccountResourceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Convert address if needed (assuming base58)
	hexAddr, err := utils.Base58ToHex(req.Address)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid address format")
		return
	}

	resp, err := h.blockchainClient.GetAccountResource(context.Background(), &lindapb.Account{
		Address: []byte(hexAddr),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get account resource: "+err.Error())
		return
	}

	response := &models.AccountResourceResponse{
		FreeNetUsed:     resp.FreeNetUsed,
		FreeNetLimit:    resp.FreeNetLimit,
		NetUsed:         resp.NetUsed,
		NetLimit:        resp.NetLimit,
		TotalNetLimit:   resp.TotalNetLimit,
		TotalNetWeight:  resp.TotalNetWeight,
		TotalLindaPowerWeight: resp.TotalLindaPowerWeight,
		LindaPowerLimit:  resp.LindaPowerLimit,
		LindaPowerUsed:   resp.LindaPowerUsed,
		EnergyUsed:      resp.EnergyUsed,
		EnergyLimit:     resp.EnergyLimit,
		TotalEnergyLimit: resp.TotalEnergyLimit,
		TotalEnergyWeight: resp.TotalEnergyWeight,
		AssetNetUsed:    resp.AssetNetUsed,
		AssetNetLimit:   resp.AssetNetLimit,
	}

	utils.RespondWithSuccess(c, response)
}

// GetAccountProposals handles GET /api/account-proposal
// Returns proposals created by an account
func (h *AccountHandler) GetAccountProposals(c *gin.Context) {
	var req models.AccountProposalRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	// Call blockchain to get proposals by proposer
	proposals, err := h.blockchainClient.ListProposals(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get proposals: "+err.Error())
		return
	}

	// Filter by proposer address
	var filtered []*models.ProposalResponse
	for _, p := range proposals.Proposals {
		proposerHex := string(p.ProposerAddress)
		proposerBase58, _ := utils.HexToBase58(proposerHex)
		
		if req.Address == "" || proposerBase58 == req.Address {
			filtered = append(filtered, convertProposalToResponse(p, true))
		}
	}

	// Apply pagination
	start := req.Start
	end := start + req.Limit
	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	paginated := filtered[start:end]

	response := gin.H{
		"proposals": paginated,
		"total":     len(filtered),
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Tag System ====================

// GetTags handles GET /external/tag
// Returns tags for addresses
func (h *AccountHandler) GetTags(c *gin.Context) {
	var req models.TagRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	var tags []*models.TagResponse
	var total int64
	var err error

	if req.Address != "" {
		// Get tags for specific address
		tags, total, err = h.tagRepo.GetTagsByAddress(req.Address, req.Start, req.Limit)
	} else {
		// Get all tags
		tags, total, err = h.tagRepo.GetAllTags(req.Start, req.Limit, req.Sort)
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get tags: "+err.Error())
		return
	}

	response := gin.H{
		"tags":  tags,
		"total": total,
	}

	utils.RespondWithSuccess(c, response)
}

// InsertTag handles POST /external/tag/insert
// Inserts a new tag
func (h *AccountHandler) InsertTag(c *gin.Context) {
	var req models.TagInsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Verify signature
	valid, err := utils.VerifySignature(req.Address, req.Signature, req.Tag)
	if err != nil || !valid {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid signature")
		return
	}

	// Insert tag
	id, err := h.tagRepo.InsertTag(&models.TagResponse{
		Address:     req.Address,
		Tag:         req.Tag,
		Description: req.Description,
		Owner:       req.Owner,
		CreatedAt:   time.Now().Unix(),
		Votes:       0,
	})

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to insert tag: "+err.Error())
		return
	}

	response := models.TagInsertResponse{
		Success: true,
		Message: "Tag created successfully",
		ID:      id,
	}

	utils.RespondWithSuccess(c, response)
}

// UpdateTag handles POST /external/tag/update
// Updates an existing tag
func (h *AccountHandler) UpdateTag(c *gin.Context) {
	var req models.TagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get existing tag
	tag, err := h.tagRepo.GetTagByID(req.ID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Tag not found")
		return
	}

	// Verify signature
	valid, err := utils.VerifySignature(tag.Owner, req.Signature, tag.Tag)
	if err != nil || !valid {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid signature")
		return
	}

	// Update tag
	err = h.tagRepo.UpdateTag(req.ID, req.Tag, req.Description)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update tag: "+err.Error())
		return
	}

	response := gin.H{
		"success": true,
		"message": "Tag updated successfully",
	}

	utils.RespondWithSuccess(c, response)
}

// DeleteTag handles POST /external/tag/delete
// Deletes a tag
func (h *AccountHandler) DeleteTag(c *gin.Context) {
	var req models.TagDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get existing tag
	tag, err := h.tagRepo.GetTagByID(req.ID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Tag not found")
		return
	}

	// Verify signature
	valid, err := utils.VerifySignature(tag.Owner, req.Signature, tag.Tag)
	if err != nil || !valid {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid signature")
		return
	}

	// Delete tag
	err = h.tagRepo.DeleteTag(req.ID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete tag: "+err.Error())
		return
	}

	response := gin.H{
		"success": true,
		"message": "Tag deleted successfully",
	}

	utils.RespondWithSuccess(c, response)
}

// RecommendTag handles GET /external/tag/recommend
// Returns recommended tags for an address
func (h *AccountHandler) RecommendTag(c *gin.Context) {
	var req models.TagRecommendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Get most popular tags for this address or overall
	tags, err := h.tagRepo.GetRecommendedTags(req.Address, req.Limit)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get recommendations: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, tags)
}

// ==================== Helper Functions ====================

func convertAccountToResponse(account *lindapb.Account, visible bool) *models.AccountResponse {
	resp := &models.AccountResponse{
		Address:            string(account.Address),
		Balance:            account.Balance,
		AccountName:        string(account.AccountName),
		CreateTime:         account.CreateTime,
		IsWitness:          account.IsWitness,
		Allowance:          account.Allowance,
		LatestWithdrawTime: account.LatestWithdrawTime,
		LatestOperationTime: account.LatestOprationTime,
		LatestConsumeTime:  account.LatestConsumeTime,
		LatestConsumeFreeTime: account.LatestConsumeFreeTime,
		NetWindowSize:      account.NetWindowSize,
		NetWindowOptimized: account.NetWindowOptimized,
	}

	// Convert address to base58 if visible
	if visible {
		base58Addr, _ := utils.HexToBase58(string(account.Address))
		resp.Address = base58Addr
	}

	// Convert Frozen (Stake 1.0)
	for _, f := range account.Frozen {
		resp.Frozen = append(resp.Frozen, models.Frozen{
			FrozenBalance: f.FrozenBalance,
			ExpireTime:    f.ExpireTime,
		})
	}

	// Convert FrozenV2 (Stake 2.0)
	for _, f := range account.FrozenV2 {
		resp.FrozenV2 = append(resp.FrozenV2, models.FreezeV2{
			Amount: f.Amount,
			Type:   f.Type.String(),
		})
	}

	// Convert UnfrozenV2
	for _, u := range account.UnfrozenV2 {
		resp.UnfrozenV2 = append(resp.UnfrozenV2, models.UnFreezeV2{
			Type:               u.Type.String(),
			UnfreezeAmount:     u.UnfreezeAmount,
			UnfreezeExpireTime: u.UnfreezeExpireTime,
		})
	}

	// Convert AccountResource
	if account.AccountResource != nil {
		resp.AccountResource = &models.AccountResource{
			DelegatedFrozenBalanceForEnergy:         account.AccountResource.DelegatedFrozenBalanceForEnergy,
			AcquiredDelegatedFrozenBalanceForEnergy: account.AccountResource.AcquiredDelegatedFrozenBalanceForEnergy,
			DelegatedFrozenV2BalanceForEnergy:       account.AccountResource.DelegatedFrozenV2BalanceForEnergy,
			AcquiredDelegatedFrozenV2BalanceForEnergy: account.AccountResource.AcquiredDelegatedFrozenV2BalanceForEnergy,
			EnergyUsage:                              account.AccountResource.EnergyUsage,
			EnergyWindowSize:                         account.AccountResource.EnergyWindowSize,
			EnergyWindowOptimized:                     account.AccountResource.EnergyWindowOptimized,
			LatestConsumeTimeForEnergy:                account.AccountResource.LatestConsumeTimeForEnergy,
		}

		if account.AccountResource.FrozenBalanceForEnergy != nil {
			resp.AccountResource.FrozenBalanceForEnergy = &models.Frozen{
				FrozenBalance: account.AccountResource.FrozenBalanceForEnergy.FrozenBalance,
				ExpireTime:    account.AccountResource.FrozenBalanceForEnergy.ExpireTime,
			}
		}
	}

	// Convert Votes
	for _, v := range account.Votes {
		voteAddr := string(v.VoteAddress)
		if visible {
			voteAddr, _ = utils.HexToBase58(voteAddr)
		}
		resp.Votes = append(resp.Votes, models.Vote{
			VoteAddress: voteAddr,
			VoteCount:   v.VoteCount,
		})
	}

	// Convert Permissions (simplified)
	if account.OwnerPermission != nil {
		resp.OwnerPermission = convertPermission(account.OwnerPermission, visible)
	}
	if account.WitnessPermission != nil {
		resp.WitnessPermission = convertPermission(account.WitnessPermission, visible)
	}
	for _, p := range account.ActivePermission {
		resp.ActivePermissions = append(resp.ActivePermissions, convertPermission(p, visible))
	}

	// Convert Assets
	resp.Asset = account.Asset
	resp.AssetV2 = account.AssetV2

	// Convert Delegated balances
	resp.DelegatedFrozenBalanceForBandwidth = account.DelegatedFrozenBalanceForBandwidth
	resp.AcquiredDelegatedFrozenBalanceForBandwidth = account.AcquiredDelegatedFrozenBalanceForBandwidth
	resp.DelegatedFrozenV2BalanceForBandwidth = account.DelegatedFrozenV2BalanceForBandwidth
	resp.AcquiredDelegatedFrozenV2BalanceForBandwidth = account.AcquiredDelegatedFrozenV2BalanceForBandwidth

	// Usage
	resp.NetUsage = account.NetUsage
	resp.FreeNetUsage = account.FreeNetUsage
	resp.FreeAssetNetUsageV2 = account.FreeAssetNetUsageV2

	return resp
}

func convertPermission(p *lindapb.Permission, visible bool) *models.Permission {
	resp := &models.Permission{
		Type:           int(p.Type),
		ID:             int(p.Id),
		PermissionName: p.PermissionName,
		Threshold:      p.Threshold,
		ParentID:       int(p.ParentId),
		Operations:     string(p.Operations),
	}

	for _, k := range p.Keys {
		addr := string(k.Address)
		if visible {
			addr, _ = utils.HexToBase58(addr)
		}
		resp.Keys = append(resp.Keys, models.Key{
			Address: addr,
			Weight:  int(k.Weight),
		})
	}

	return resp
}

func convertTransactionToResponse(tx *lindapb.Transaction, visible bool) *models.TransactionResponse {
	resp := &models.TransactionResponse{
		TxID:       string(tx.TxID),
		RawDataHex: string(tx.RawDataHex),
		Signature:  tx.Signature,
	}

	// Convert Ret
	for _, r := range tx.Ret {
		resp.Ret = append(resp.Ret, models.TransactionResult{
			ContractRet: r.Ret.String(),
			Fee:         r.Fee,
		})
	}

	// Convert RawData if present
	if tx.RawData != nil {
		resp.RawData = &models.TransactionRawData{
			RefBlockBytes: string(tx.RawData.RefBlockBytes),
			RefBlockNum:   tx.RawData.RefBlockNum,
			RefBlockHash:  string(tx.RawData.RefBlockHash),
			Expiration:    tx.RawData.Expiration,
			Data:          string(tx.RawData.Data),
			Scripts:       string(tx.RawData.Scripts),
			Timestamp:     tx.RawData.Timestamp,
		}

		// Convert Contracts
		for _, c := range tx.RawData.Contract {
			resp.RawData.Contract = append(resp.RawData.Contract, models.TransactionContract{
				Type:      c.Type.String(),
				Parameter: c.Parameter.Value,
			})
		}
	}

	return resp
}

func convertProposalToResponse(p *lindapb.Proposal, visible bool) *models.ProposalResponse {
	resp := &models.ProposalResponse{
		ProposalID:      p.ProposalId,
		ProposerAddress: string(p.ProposerAddress),
		Parameters:      p.Parameters,
		ExpirationTime:  p.ExpirationTime,
		CreateTime:      p.CreateTime,
		Approvals:       p.Approvals,
		State:           p.State.String(),
	}

	if visible {
		base58Addr, _ := utils.HexToBase58(string(p.ProposerAddress))
		resp.ProposerAddress = base58Addr

		for i, addr := range p.Approvals {
			base58Addr, _ := utils.HexToBase58(addr)
			resp.Approvals[i] = base58Addr
		}
	}

	return resp
}