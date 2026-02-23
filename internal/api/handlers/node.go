package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type NodeHandler struct {
	blockchainClient *blockchain.Client
}

func NewNodeHandler(client *blockchain.Client) *NodeHandler {
	return &NodeHandler{
		blockchainClient: client,
	}
}

// ==================== Node Methods ====================

// ListNodes handles POST /wallet/listnodes
func (h *NodeHandler) ListNodes(c *gin.Context) {
	visible := c.DefaultQuery("visible", "false") == "true"

	nodes, err := h.blockchainClient.ListNodes(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list nodes: "+err.Error())
		return
	}

	var response []gin.H
	for _, node := range nodes.Nodes {
		nodeInfo := gin.H{
			"address": gin.H{
				"host": string(node.Address.Host),
				"port": node.Address.Port,
			},
		}
		response = append(response, nodeInfo)
	}

	utils.RespondWithSuccess(c, gin.H{"nodes": response})
}

// GetNodeInfo handles POST /wallet/getnodeinfo
func (h *NodeHandler) GetNodeInfo(c *gin.Context) {
	info, err := h.blockchainClient.GetNodeInfo(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get node info: "+err.Error())
		return
	}

	response := convertNodeInfoToResponse(info)
	utils.RespondWithSuccess(c, response)
}

// GetNodeInfoSolidity handles POST /walletsolidity/getnodeinfo
func (h *NodeHandler) GetNodeInfoSolidity(c *gin.Context) {
	info, err := h.blockchainClient.GetNodeInfoSolidity(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get node info: "+err.Error())
		return
	}

	response := convertNodeInfoToResponse(info)
	utils.RespondWithSuccess(c, response)
}

// ==================== Witness & Voting Methods ====================

// ListWitnesses handles POST /wallet/listwitnesses
func (h *NodeHandler) ListWitnesses(c *gin.Context) {
	visible := c.DefaultQuery("visible", "false") == "true"

	witnesses, err := h.blockchainClient.ListWitnesses(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list witnesses: "+err.Error())
		return
	}

	var response []gin.H
	for _, w := range witnesses.Witnesses {
		witness := gin.H{
			"address":        string(w.Address),
			"voteCount":      w.VoteCount,
			"url":            w.Url,
			"totalProduced":  w.TotalProduced,
			"totalMissed":    w.TotalMissed,
			"latestBlockNum": w.LatestBlockNum,
			"latestSlotNum":  w.LatestSlotNum,
			"isJobs":         w.IsJobs,
		}
		if visible {
			base58Addr, _ := utils.HexToBase58(string(w.Address))
			witness["address"] = base58Addr
		}
		response = append(response, witness)
	}

	utils.RespondWithSuccess(c, gin.H{"witnesses": response})
}

// ListWitnessesSolidity handles GET /walletsolidity/listwitnesses
func (h *NodeHandler) ListWitnessesSolidity(c *gin.Context) {
	visible := c.DefaultQuery("visible", "false") == "true"

	witnesses, err := h.blockchainClient.ListWitnessesSolidity(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list witnesses: "+err.Error())
		return
	}

	var response []gin.H
	for _, w := range witnesses.Witnesses {
		witness := gin.H{
			"address":        string(w.Address),
			"voteCount":      w.VoteCount,
			"url":            w.Url,
			"totalProduced":  w.TotalProduced,
			"totalMissed":    w.TotalMissed,
			"latestBlockNum": w.LatestBlockNum,
			"latestSlotNum":  w.LatestSlotNum,
			"isJobs":         w.IsJobs,
		}
		if visible {
			base58Addr, _ := utils.HexToBase58(string(w.Address))
			witness["address"] = base58Addr
		}
		response = append(response, witness)
	}

	utils.RespondWithSuccess(c, gin.H{"witnesses": response})
}

// CreateWitness handles POST /wallet/createwitness
func (h *NodeHandler) CreateWitness(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		URL          string `json:"url" binding:"required"`
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

	tx, err := h.blockchainClient.CreateWitness(context.Background(), &lindapb.WitnessCreateContract{
		OwnerAddress: []byte(ownerAddress),
		Url:          []byte(req.URL),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create witness: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// UpdateWitness handles POST /wallet/updatewitness
func (h *NodeHandler) UpdateWitness(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		UpdateURL    string `json:"update_url" binding:"required"`
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

	tx, err := h.blockchainClient.UpdateWitness(context.Background(), &lindapb.WitnessUpdateContract{
		OwnerAddress: []byte(ownerAddress),
		UpdateUrl:    []byte(req.UpdateURL),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update witness: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// VoteWitnessAccount handles POST /wallet/votewitnessaccount
func (h *NodeHandler) VoteWitnessAccount(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		Votes        []struct {
			VoteAddress string `json:"vote_address"`
			VoteCount   int64  `json:"vote_count"`
		} `json:"votes" binding:"required"`
		Visible      bool  `json:"visible" default:"false"`
		PermissionID int32 `json:"permission_id"`
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

	votes := make([]*lindapb.VoteWitnessContract_Vote, len(req.Votes))
	for i, v := range req.Votes {
		voteAddress := v.VoteAddress
		if req.Visible {
			hexAddr, _ := utils.Base58ToHex(voteAddress)
			voteAddress = hexAddr
		}
		votes[i] = &lindapb.VoteWitnessContract_Vote{
			VoteAddress: []byte(voteAddress),
			VoteCount:   v.VoteCount,
		}
	}

	tx, err := h.blockchainClient.VoteWitnessAccount(context.Background(), &lindapb.VoteWitnessContract{
		OwnerAddress: []byte(ownerAddress),
		Votes:        votes,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to vote: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetBrokerage handles POST /wallet/getBrokerage
func (h *NodeHandler) GetBrokerage(c *gin.Context) {
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
		hexAddr, _ := utils.Base58ToHex(address)
		address = hexAddr
	}

	brokerage, err := h.blockchainClient.GetBrokerage(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get brokerage: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"brokerage": brokerage.Num})
}

// GetBrokerageSolidity handles POST /walletsolidity/getBrokerage
func (h *NodeHandler) GetBrokerageSolidity(c *gin.Context) {
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
		hexAddr, _ := utils.Base58ToHex(address)
		address = hexAddr
	}

	brokerage, err := h.blockchainClient.GetBrokerageSolidity(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get brokerage: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"brokerage": brokerage.Num})
}

// UpdateBrokerage handles POST /wallet/updateBrokerage
func (h *NodeHandler) UpdateBrokerage(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		Brokerage    int32  `json:"brokerage" binding:"required"`
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

	tx, err := h.blockchainClient.UpdateBrokerage(context.Background(), &lindapb.UpdateBrokerageContract{
		OwnerAddress: []byte(ownerAddress),
		Brokerage:    req.Brokerage,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update brokerage: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetReward handles POST /wallet/getReward
func (h *NodeHandler) GetReward(c *gin.Context) {
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
		hexAddr, _ := utils.Base58ToHex(address)
		address = hexAddr
	}

	reward, err := h.blockchainClient.GetReward(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get reward: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"reward": reward.Num})
}

// GetRewardSolidity handles POST /walletsolidity/getReward
func (h *NodeHandler) GetRewardSolidity(c *gin.Context) {
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
		hexAddr, _ := utils.Base58ToHex(address)
		address = hexAddr
	}

	reward, err := h.blockchainClient.GetRewardSolidity(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get reward: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"reward": reward.Num})
}

// GetNextMaintenanceTime handles GET /wallet/getnextmaintenancetime
func (h *NodeHandler) GetNextMaintenanceTime(c *gin.Context) {
	time, err := h.blockchainClient.GetNextMaintenanceTime(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get next maintenance time: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"num": time.Num})
}

// GetPaginatedNowWitnessList handles POST /wallet/getpaginatednowwitnesslist
func (h *NodeHandler) GetPaginatedNowWitnessList(c *gin.Context) {
	var req struct {
		Offset  int64 `json:"offset"`
		Limit   int64 `json:"limit" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	witnesses, err := h.blockchainClient.GetPaginatedNowWitnessList(context.Background(), &lindapb.PaginatedMessage{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get paginated witness list: "+err.Error())
		return
	}

	var response []gin.H
	for _, w := range witnesses.Witnesses {
		witness := gin.H{
			"address":        string(w.Address),
			"voteCount":      w.VoteCount,
			"url":            w.Url,
			"totalProduced":  w.TotalProduced,
			"totalMissed":    w.TotalMissed,
			"latestBlockNum": w.LatestBlockNum,
			"latestSlotNum":  w.LatestSlotNum,
			"isJobs":         w.IsJobs,
		}
		if req.Visible {
			base58Addr, _ := utils.HexToBase58(string(w.Address))
			witness["address"] = base58Addr
		}
		response = append(response, witness)
	}

	utils.RespondWithSuccess(c, gin.H{"witnesses": response})
}

// GetPaginatedNowWitnessListSolidity handles GET /soliditywallet/getpaginatednowwitnesslist
func (h *NodeHandler) GetPaginatedNowWitnessListSolidity(c *gin.Context) {
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	visible := c.DefaultQuery("visible", "false") == "true"

	if limit <= 0 || limit > 1000 {
		limit = 20
	}

	witnesses, err := h.blockchainClient.GetPaginatedNowWitnessListSolidity(context.Background(), &lindapb.PaginatedMessage{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get paginated witness list: "+err.Error())
		return
	}

	var response []gin.H
	for _, w := range witnesses.Witnesses {
		witness := gin.H{
			"address":        string(w.Address),
			"voteCount":      w.VoteCount,
			"url":            w.Url,
			"totalProduced":  w.TotalProduced,
			"totalMissed":    w.TotalMissed,
			"latestBlockNum": w.LatestBlockNum,
			"latestSlotNum":  w.LatestSlotNum,
			"isJobs":         w.IsJobs,
		}
		if visible {
			base58Addr, _ := utils.HexToBase58(string(w.Address))
			witness["address"] = base58Addr
		}
		response = append(response, witness)
	}

	utils.RespondWithSuccess(c, gin.H{"witnesses": response})
}

// ==================== Proposal Methods ====================

// ListProposals handles GET /wallet/listproposals
func (h *NodeHandler) ListProposals(c *gin.Context) {
	visible := c.DefaultQuery("visible", "false") == "true"

	proposals, err := h.blockchainClient.ListProposals(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list proposals: "+err.Error())
		return
	}

	var response []gin.H
	for _, p := range proposals.Proposals {
		response = append(response, convertProposalToResponse(p, visible))
	}

	utils.RespondWithSuccess(c, gin.H{"proposals": response})
}

// GetProposalById handles POST /wallet/getproposalbyid
func (h *NodeHandler) GetProposalById(c *gin.Context) {
	var req struct {
		ID      int64 `json:"id" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	proposal, err := h.blockchainClient.GetProposalById(context.Background(), &lindapb.NumberMessage{
		Num: req.ID,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get proposal: "+err.Error())
		return
	}

	response := convertProposalToResponse(proposal, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ProposalCreate handles POST /wallet/proposalcreate
func (h *NodeHandler) ProposalCreate(c *gin.Context) {
	var req struct {
		OwnerAddress string           `json:"owner_address" binding:"required"`
		Parameters   map[int64]int64 `json:"parameters" binding:"required"`
		Visible      bool             `json:"visible" default:"false"`
		PermissionID int32            `json:"permission_id"`
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

	tx, err := h.blockchainClient.ProposalCreate(context.Background(), &lindapb.ProposalCreateContract{
		OwnerAddress: []byte(ownerAddress),
		Parameters:   req.Parameters,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create proposal: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ProposalApprove handles POST /wallet/proposalapprove
func (h *NodeHandler) ProposalApprove(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		ProposalID   int64  `json:"proposal_id" binding:"required"`
		IsAddApproval bool  `json:"is_add_approval" binding:"required"`
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

	tx, err := h.blockchainClient.ProposalApprove(context.Background(), &lindapb.ProposalApproveContract{
		OwnerAddress:  []byte(ownerAddress),
		ProposalId:    req.ProposalID,
		IsAddApproval: req.IsAddApproval,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to approve proposal: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ProposalDelete handles POST /wallet/proposaldelete
func (h *NodeHandler) ProposalDelete(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		ProposalID   int64  `json:"proposal_id" binding:"required"`
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

	tx, err := h.blockchainClient.ProposalDelete(context.Background(), &lindapb.ProposalDeleteContract{
		OwnerAddress: []byte(ownerAddress),
		ProposalId:   req.ProposalID,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete proposal: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetPaginatedProposalList handles POST /wallet/getpaginatedproposallist
func (h *NodeHandler) GetPaginatedProposalList(c *gin.Context) {
	var req struct {
		Offset  int64 `json:"offset"`
		Limit   int64 `json:"limit" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	proposals, err := h.blockchainClient.GetPaginatedProposalList(context.Background(), &lindapb.PaginatedMessage{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get paginated proposals: "+err.Error())
		return
	}

	var response []gin.H
	for _, p := range proposals.Proposals {
		response = append(response, convertProposalToResponse(p, req.Visible))
	}

	utils.RespondWithSuccess(c, gin.H{"proposals": response})
}

// ==================== Chain Info Methods ====================

// GetChainParameters handles GET /wallet/getchainparameters
func (h *NodeHandler) GetChainParameters(c *gin.Context) {
	params, err := h.blockchainClient.GetChainParameters(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get chain parameters: "+err.Error())
		return
	}

	var response []gin.H
	for _, p := range params.ChainParameter {
		response = append(response, gin.H{
			"key":   p.Key,
			"value": p.Value,
		})
	}

	utils.RespondWithSuccess(c, gin.H{"chainParameter": response})
}

// GetChainParametersV2 handles GET /api/chainparameters
func (h *NodeHandler) GetChainParametersV2(c *gin.Context) {
	params, err := h.blockchainClient.GetChainParameters(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get chain parameters: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, params.ChainParameter)
}

// TotalTransaction handles GET /wallet/totaltransaction
func (h *NodeHandler) TotalTransaction(c *gin.Context) {
	total, err := h.blockchainClient.TotalTransaction(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get total transactions: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"num": total.Num})
}

// GetBurnLind handles GET /wallet/getburnlind
func (h *NodeHandler) GetBurnLind(c *gin.Context) {
	burn, err := h.blockchainClient.GetBurnLind(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get burn LIND: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"burnLindAmount": burn.Num})
}

// GetBurnLindSolidity handles GET /walletsolidity/getburnlind
func (h *NodeHandler) GetBurnLindSolidity(c *gin.Context) {
	burn, err := h.blockchainClient.GetBurnLindSolidity(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get burn LIND: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"burnLindAmount": burn.Num})
}

// GetEnergyPrices handles GET /wallet/getenergyprices
func (h *NodeHandler) GetEnergyPrices(c *gin.Context) {
	prices, err := h.blockchainClient.GetEnergyPrices(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get energy prices: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"prices": string(prices.Prices)})
}

// GetBandwidthPrices handles GET /wallet/getbandwidthprices
func (h *NodeHandler) GetBandwidthPrices(c *gin.Context) {
	prices, err := h.blockchainClient.GetBandwidthPrices(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get bandwidth prices: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"prices": string(prices.Prices)})
}

// GetMemoFeePrices handles GET /wallet/getmemofee
func (h *NodeHandler) GetMemoFeePrices(c *gin.Context) {
	prices, err := h.blockchainClient.GetMemoFeePrices(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get memo fee prices: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"prices": string(prices.Prices)})
}

// ==================== Metrics Methods ====================

// GetStatsInfo handles GET /monitor/getstatsinfo
func (h *NodeHandler) GetStatsInfo(c *gin.Context) {
	stats, err := h.blockchainClient.GetStatsInfo(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get stats info: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, stats)
}

// ==================== Vote Info (Lindascan) ====================

// GetVoteInfo handles GET /api/vote
func (h *NodeHandler) GetVoteInfo(c *gin.Context) {
	// Get witnesses with their vote counts
	witnesses, err := h.blockchainClient.ListWitnesses(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get vote info: "+err.Error())
		return
	}

	var response []gin.H
	for _, w := range witnesses.Witnesses {
		response = append(response, gin.H{
			"address":   utils.MustHexToBase58(string(w.Address)),
			"voteCount": w.VoteCount,
			"url":       w.Url,
		})
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Node Map (Lindascan) ====================

// GetNodeMap handles GET /api/nodemap
func (h *NodeHandler) GetNodeMap(c *gin.Context) {
	nodes, err := h.blockchainClient.ListNodes(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get node map: "+err.Error())
		return
	}

	var response []gin.H
	for _, node := range nodes.Nodes {
		// In production, you'd use IP geolocation service
		response = append(response, gin.H{
			"ip":        string(node.Address.Host),
			"host":      string(node.Address.Host),
			"port":      node.Address.Port,
			"country":   "Unknown",
			"city":      "Unknown",
			"latitude":  0,
			"longitude": 0,
			"nodeType":  "FullNode",
		})
	}

	utils.RespondWithSuccess(c, gin.H{"nodes": response})
}

// ==================== Node Upload (V2) ====================

// UploadNodeOverview handles POST /api/v2/node/overview_upload
func (h *NodeHandler) UploadNodeOverview(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Process upload (would store in database)
	utils.RespondWithSuccess(c, gin.H{
		"uptime":      3600,
		"block_height": 1000000,
		"peers":       25,
		"version":     "4.5.0",
	})
}

// UploadNodeInfo handles POST /api/v2/node/info_upload
func (h *NodeHandler) UploadNodeInfo(c *gin.Context) {
	var req struct {
		Address     string `json:"address" binding:"required"`
		Version     string `json:"version"`
		Location    string `json:"location"`
		BlockHeight int64  `json:"block_height"`
		Peers       int    `json:"peers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Process upload (would store in database)
	utils.RespondWithSuccess(c, gin.H{
		"success": true,
		"message": "Node info uploaded successfully",
	})
}

// ==================== Helper Functions ====================

func convertNodeInfoToResponse(info *lindapb.NodeInfo) gin.H {
	return gin.H{
		"beginSyncNum":        info.BeginSyncNum,
		"block":               info.Block,
		"solidityBlock":       info.SolidityBlock,
		"currentConnectCount": info.CurrentConnectCount,
		"activeConnectCount":  info.ActiveConnectCount,
		"passiveConnectCount": info.PassiveConnectCount,
		"totalFlow":           info.TotalFlow,
		"peerInfoList":        info.PeerInfoList,
		"configNodeInfo":      info.ConfigNodeInfo,
		"machineInfo":         info.MachineInfo,
		"cheatWitnessInfoMap": info.CheatWitnessInfoMap,
	}
}

func convertProposalToResponse(p *lindapb.Proposal, visible bool) gin.H {
	state := "PENDING"
	switch p.State {
	case 1:
		state = "DISAPPROVED"
	case 2:
		state = "APPROVED"
	case 3:
		state = "CANCELED"
	}

	resp := gin.H{
		"proposal_id":      p.ProposalId,
		"proposer_address": string(p.ProposerAddress),
		"parameters":       p.Parameters,
		"expiration_time":  p.ExpirationTime,
		"create_time":      p.CreateTime,
		"approvals":        p.Approvals,
		"state":            state,
	}

	if visible {
		base58Addr, _ := utils.HexToBase58(string(p.ProposerAddress))
		resp["proposer_address"] = base58Addr

		approvals := make([]string, len(p.Approvals))
		for i, addr := range p.Approvals {
			base58Addr, _ := utils.HexToBase58(addr)
			approvals[i] = base58Addr
		}
		resp["approvals"] = approvals
	}

	return resp
}