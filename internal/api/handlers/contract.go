package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type ContractHandler struct {
	blockchainClient *blockchain.Client
}

func NewContractHandler(client *blockchain.Client) *ContractHandler {
	return &ContractHandler{
		blockchainClient: client,
	}
}

// ==================== Smart Contract Methods ====================

// DeployContract handles POST /wallet/deploycontract
func (h *ContractHandler) DeployContract(c *gin.Context) {
	var req struct {
		OwnerAddress               string          `json:"owner_address" binding:"required"`
		ABI                        json.RawMessage `json:"abi" binding:"required"`
		Bytecode                   string          `json:"bytecode" binding:"required"`
		Name                       string          `json:"name"`
		FeeLimit                   int64           `json:"fee_limit" binding:"required"`
		Parameter                  string          `json:"parameter"`
		OriginEnergyLimit          int64           `json:"origin_energy_limit"`
		CallValue                  int64           `json:"call_value"`
		ConsumeUserResourcePercent int32           `json:"consume_user_resource_percent"`
		TokenID                    int64           `json:"token_id"`
		TokenValue                 int64           `json:"token_value"`
		Visible                    bool            `json:"visible" default:"false"`
		PermissionID               int32           `json:"permission_id"`
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

	// Create smart contract object
	newContract := &lindapb.SmartContract{
		OriginAddress:              []byte(ownerAddress),
		Abi:                        req.ABI,
		Bytecode:                   []byte(req.Bytecode),
		Name:                       req.Name,
		ConsumeUserResourcePercent: req.ConsumeUserResourcePercent,
		OriginEnergyLimit:          req.OriginEnergyLimit,
		CallValue:                  req.CallValue,
	}

	// Create contract
	contract := &lindapb.CreateSmartContract{
		OwnerAddress: []byte(ownerAddress),
		NewContract:  newContract,
		FeeLimit:     req.FeeLimit,
	}

	if req.TokenID > 0 {
		contract.TokenId = req.TokenID
		contract.CallTokenValue = req.TokenValue
	}

	result, err := h.blockchainClient.DeployContract(context.Background(), contract)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to deploy contract: "+err.Error())
		return
	}

	response := gin.H{
		"transaction": convertTransactionToResponse(result.Transaction, req.Visible),
		"txid":        string(result.Txid),
		"result":      result.Result,
	}

	utils.RespondWithSuccess(c, response)
}

// TriggerSmartContract handles POST /wallet/triggersmartcontract
func (h *ContractHandler) TriggerSmartContract(c *gin.Context) {
	var req struct {
		OwnerAddress     string `json:"owner_address" binding:"required"`
		ContractAddress  string `json:"contract_address" binding:"required"`
		FunctionSelector string `json:"function_selector"`
		Parameter        string `json:"parameter"`
		Data             string `json:"data"`
		FeeLimit         int64  `json:"fee_limit" binding:"required"`
		CallValue        int64  `json:"call_value"`
		CallTokenValue   int64  `json:"call_token_value"`
		TokenID          int64  `json:"token_id"`
		Visible          bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	contractAddress := req.ContractAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexContract, _ := utils.Base58ToHex(contractAddress)
		ownerAddress = hexOwner
		contractAddress = hexContract
	}

	triggerReq := &lindapb.TriggerSmartContractReq{
		OwnerAddress:     []byte(ownerAddress),
		ContractAddress:  []byte(contractAddress),
		FeeLimit:         req.FeeLimit,
		CallValue:        req.CallValue,
		CallTokenValue:   req.CallTokenValue,
		TokenId:          req.TokenID,
	}

	// Use either function_selector+parameter or data
	if req.FunctionSelector != "" {
		triggerReq.FunctionSelector = req.FunctionSelector
		triggerReq.Parameter = req.Parameter
	} else {
		triggerReq.Data = []byte(req.Data)
	}

	result, err := h.blockchainClient.TriggerSmartContract(context.Background(), triggerReq)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to trigger contract: "+err.Error())
		return
	}

	response := gin.H{
		"transaction": convertTransactionToResponse(result.Transaction, req.Visible),
		"txid":        string(result.Txid),
		"result":      result.Result,
	}

	utils.RespondWithSuccess(c, response)
}

// TriggerConstantContract handles POST /wallet/triggerconstantcontract
func (h *ContractHandler) TriggerConstantContract(c *gin.Context) {
	var req struct {
		OwnerAddress     string `json:"owner_address" binding:"required"`
		ContractAddress  string `json:"contract_address" binding:"required"`
		FunctionSelector string `json:"function_selector"`
		Parameter        string `json:"parameter"`
		Data             string `json:"data"`
		CallValue        int64  `json:"call_value"`
		CallTokenValue   int64  `json:"call_token_value"`
		TokenID          int64  `json:"token_id"`
		Visible          bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	contractAddress := req.ContractAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexContract, _ := utils.Base58ToHex(contractAddress)
		ownerAddress = hexOwner
		contractAddress = hexContract
	}

	constantReq := &lindapb.TriggerSmartContractReq{
		OwnerAddress:     []byte(ownerAddress),
		ContractAddress:  []byte(contractAddress),
		CallValue:        req.CallValue,
		CallTokenValue:   req.CallTokenValue,
		TokenId:          req.TokenID,
	}

	if req.FunctionSelector != "" {
		constantReq.FunctionSelector = req.FunctionSelector
		constantReq.Parameter = req.Parameter
	} else {
		constantReq.Data = []byte(req.Data)
	}

	result, err := h.blockchainClient.TriggerConstantContract(context.Background(), constantReq)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to trigger constant contract: "+err.Error())
		return
	}

	response := gin.H{
		"result":          result.Result,
		"energy_used":     result.EnergyUsed,
		"energy_penalty":  result.EnergyPenalty,
		"constant_result": result.ConstantResult,
		"transaction":     convertTransactionToResponse(result.Transaction, req.Visible),
	}

	utils.RespondWithSuccess(c, response)
}

// TriggerConstantContractSolidity handles POST /walletsolidity/triggerconstantcontract
func (h *ContractHandler) TriggerConstantContractSolidity(c *gin.Context) {
	var req struct {
		OwnerAddress     string `json:"owner_address" binding:"required"`
		ContractAddress  string `json:"contract_address" binding:"required"`
		FunctionSelector string `json:"function_selector"`
		Parameter        string `json:"parameter"`
		Data             string `json:"data"`
		CallValue        int64  `json:"call_value"`
		CallTokenValue   int64  `json:"call_token_value"`
		TokenID          int64  `json:"token_id"`
		Visible          bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	contractAddress := req.ContractAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexContract, _ := utils.Base58ToHex(contractAddress)
		ownerAddress = hexOwner
		contractAddress = hexContract
	}

	constantReq := &lindapb.TriggerSmartContractReq{
		OwnerAddress:     []byte(ownerAddress),
		ContractAddress:  []byte(contractAddress),
		CallValue:        req.CallValue,
		CallTokenValue:   req.CallTokenValue,
		TokenId:          req.TokenID,
	}

	if req.FunctionSelector != "" {
		constantReq.FunctionSelector = req.FunctionSelector
		constantReq.Parameter = req.Parameter
	} else {
		constantReq.Data = []byte(req.Data)
	}

	result, err := h.blockchainClient.TriggerConstantContractSolidity(context.Background(), constantReq)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to trigger constant contract: "+err.Error())
		return
	}

	response := gin.H{
		"result":          result.Result,
		"energy_used":     result.EnergyUsed,
		"energy_penalty":  result.EnergyPenalty,
		"constant_result": result.ConstantResult,
		"transaction":     convertTransactionToResponse(result.Transaction, req.Visible),
	}

	utils.RespondWithSuccess(c, response)
}

// EstimateEnergy handles POST /wallet/estimateenergy
func (h *ContractHandler) EstimateEnergy(c *gin.Context) {
	var req struct {
		OwnerAddress     string `json:"owner_address" binding:"required"`
		ContractAddress  string `json:"contract_address" binding:"required"`
		FunctionSelector string `json:"function_selector"`
		Parameter        string `json:"parameter"`
		Data             string `json:"data"`
		CallValue        int64  `json:"call_value"`
		CallTokenValue   int64  `json:"call_token_value"`
		TokenID          int64  `json:"token_id"`
		Visible          bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	contractAddress := req.ContractAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexContract, _ := utils.Base58ToHex(contractAddress)
		ownerAddress = hexOwner
		contractAddress = hexContract
	}

	estimateReq := &lindapb.TriggerSmartContractReq{
		OwnerAddress:     []byte(ownerAddress),
		ContractAddress:  []byte(contractAddress),
		CallValue:        req.CallValue,
		CallTokenValue:   req.CallTokenValue,
		TokenId:          req.TokenID,
	}

	if req.FunctionSelector != "" {
		estimateReq.FunctionSelector = req.FunctionSelector
		estimateReq.Parameter = req.Parameter
	} else {
		estimateReq.Data = []byte(req.Data)
	}

	result, err := h.blockchainClient.EstimateEnergy(context.Background(), estimateReq)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to estimate energy: "+err.Error())
		return
	}

	response := gin.H{
		"result":          result.Result,
		"energy_required": result.EnergyRequired,
	}

	utils.RespondWithSuccess(c, response)
}

// EstimateEnergySolidity handles POST /walletsolidity/estimateenergy
func (h *ContractHandler) EstimateEnergySolidity(c *gin.Context) {
	var req struct {
		OwnerAddress     string `json:"owner_address" binding:"required"`
		ContractAddress  string `json:"contract_address" binding:"required"`
		FunctionSelector string `json:"function_selector"`
		Parameter        string `json:"parameter"`
		Data             string `json:"data"`
		CallValue        int64  `json:"call_value"`
		CallTokenValue   int64  `json:"call_token_value"`
		TokenID          int64  `json:"token_id"`
		Visible          bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	contractAddress := req.ContractAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexContract, _ := utils.Base58ToHex(contractAddress)
		ownerAddress = hexOwner
		contractAddress = hexContract
	}

	estimateReq := &lindapb.TriggerSmartContractReq{
		OwnerAddress:     []byte(ownerAddress),
		ContractAddress:  []byte(contractAddress),
		CallValue:        req.CallValue,
		CallTokenValue:   req.CallTokenValue,
		TokenId:          req.TokenID,
	}

	if req.FunctionSelector != "" {
		estimateReq.FunctionSelector = req.FunctionSelector
		estimateReq.Parameter = req.Parameter
	} else {
		estimateReq.Data = []byte(req.Data)
	}

	result, err := h.blockchainClient.EstimateEnergySolidity(context.Background(), estimateReq)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to estimate energy: "+err.Error())
		return
	}

	response := gin.H{
		"result":          result.Result,
		"energy_required": result.EnergyRequired,
	}

	utils.RespondWithSuccess(c, response)
}

// GetContract handles POST /wallet/getcontract
func (h *ContractHandler) GetContract(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	contractAddr := req.Value
	if req.Visible {
		hexAddr, _ := utils.Base58ToHex(contractAddr)
		contractAddr = hexAddr
	}

	contract, err := h.blockchainClient.GetContract(context.Background(), &lindapb.BytesMessage{
		Value: []byte(contractAddr),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get contract: "+err.Error())
		return
	}

	response := convertSmartContractToResponse(contract, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetContractInfo handles POST /wallet/getcontractinfo
func (h *ContractHandler) GetContractInfo(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	contractAddr := req.Value
	if req.Visible {
		hexAddr, _ := utils.Base58ToHex(contractAddr)
		contractAddr = hexAddr
	}

	info, err := h.blockchainClient.GetContractInfo(context.Background(), &lindapb.BytesMessage{
		Value: []byte(contractAddr),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get contract info: "+err.Error())
		return
	}

	response := gin.H{
		"runtimecode":   string(info.Runtimecode),
		"smart_contract": convertSmartContractToResponse(info.SmartContract, req.Visible),
		"contract_state": info.ContractState,
	}

	utils.RespondWithSuccess(c, response)
}

// UpdateSetting handles POST /wallet/updatesetting
func (h *ContractHandler) UpdateSetting(c *gin.Context) {
	var req struct {
		OwnerAddress                string `json:"owner_address" binding:"required"`
		ContractAddress             string `json:"contract_address" binding:"required"`
		ConsumeUserResourcePercent int32  `json:"consume_user_resource_percent" binding:"required"`
		Visible                     bool   `json:"visible" default:"false"`
		PermissionID                int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	contractAddress := req.ContractAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexContract, _ := utils.Base58ToHex(contractAddress)
		ownerAddress = hexOwner
		contractAddress = hexContract
	}

	tx, err := h.blockchainClient.UpdateSetting(context.Background(), &lindapb.UpdateSettingContract{
		OwnerAddress:                []byte(ownerAddress),
		ContractAddress:             []byte(contractAddress),
		ConsumeUserResourcePercent: req.ConsumeUserResourcePercent,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update setting: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// UpdateEnergyLimit handles POST /wallet/updateenergylimit
func (h *ContractHandler) UpdateEnergyLimit(c *gin.Context) {
	var req struct {
		OwnerAddress      string `json:"owner_address" binding:"required"`
		ContractAddress   string `json:"contract_address" binding:"required"`
		OriginEnergyLimit int64  `json:"origin_energy_limit" binding:"required"`
		Visible           bool   `json:"visible" default:"false"`
		PermissionID      int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	contractAddress := req.ContractAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexContract, _ := utils.Base58ToHex(contractAddress)
		ownerAddress = hexOwner
		contractAddress = hexContract
	}

	tx, err := h.blockchainClient.UpdateEnergyLimit(context.Background(), &lindapb.UpdateEnergyLimitContract{
		OwnerAddress:      []byte(ownerAddress),
		ContractAddress:   []byte(contractAddress),
		OriginEnergyLimit: req.OriginEnergyLimit,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update energy limit: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ClearAbi handles POST /wallet/clearabi
func (h *ContractHandler) ClearAbi(c *gin.Context) {
	var req struct {
		OwnerAddress    string `json:"owner_address" binding:"required"`
		ContractAddress string `json:"contract_address" binding:"required"`
		Visible         bool   `json:"visible" default:"false"`
		PermissionID    int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	contractAddress := req.ContractAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexContract, _ := utils.Base58ToHex(contractAddress)
		ownerAddress = hexOwner
		contractAddress = hexContract
	}

	tx, err := h.blockchainClient.ClearAbi(context.Background(), &lindapb.ClearAbiContract{
		OwnerAddress:    []byte(ownerAddress),
		ContractAddress: []byte(contractAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to clear ABI: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ==================== Resource Methods (Stake 1.0) ====================

// FreezeBalance handles POST /wallet/freezebalance
func (h *ContractHandler) FreezeBalance(c *gin.Context) {
	var req struct {
		OwnerAddress    string `json:"owner_address" binding:"required"`
		FrozenBalance   int64  `json:"frozen_balance" binding:"required"`
		FrozenDuration  int32  `json:"frozen_duration" default:"3"`
		Resource        string `json:"resource" binding:"required"`
		ReceiverAddress string `json:"receiver_address"`
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

	var receiverAddress []byte
	if req.ReceiverAddress != "" {
		addr := req.ReceiverAddress
		if req.Visible {
			hexAddr, _ := utils.Base58ToHex(addr)
			addr = hexAddr
		}
		receiverAddress = []byte(addr)
	}

	resourceType := lindapb.ResourceCode_BANDWIDTH
	if req.Resource == "ENERGY" {
		resourceType = lindapb.ResourceCode_ENERGY
	}

	tx, err := h.blockchainClient.FreezeBalance(context.Background(), &lindapb.FreezeBalanceContract{
		OwnerAddress:    []byte(ownerAddress),
		FrozenBalance:   req.FrozenBalance,
		FrozenDuration:  req.FrozenDuration,
		Resource:        resourceType,
		ReceiverAddress: receiverAddress,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to freeze balance: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// UnfreezeBalance handles POST /wallet/unfreezebalance
func (h *ContractHandler) UnfreezeBalance(c *gin.Context) {
	var req struct {
		OwnerAddress    string `json:"owner_address" binding:"required"`
		Resource        string `json:"resource" binding:"required"`
		ReceiverAddress string `json:"receiver_address"`
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

	var receiverAddress []byte
	if req.ReceiverAddress != "" {
		addr := req.ReceiverAddress
		if req.Visible {
			hexAddr, _ := utils.Base58ToHex(addr)
			addr = hexAddr
		}
		receiverAddress = []byte(addr)
	}

	resourceType := lindapb.ResourceCode_BANDWIDTH
	if req.Resource == "ENERGY" {
		resourceType = lindapb.ResourceCode_ENERGY
	}

	tx, err := h.blockchainClient.UnfreezeBalance(context.Background(), &lindapb.UnfreezeBalanceContract{
		OwnerAddress:    []byte(ownerAddress),
		Resource:        resourceType,
		ReceiverAddress: receiverAddress,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to unfreeze balance: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// WithdrawBalance handles POST /wallet/withdrawbalance
func (h *ContractHandler) WithdrawBalance(c *gin.Context) {
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

	tx, err := h.blockchainClient.WithdrawBalance(context.Background(), &lindapb.WithdrawBalanceContract{
		OwnerAddress: []byte(ownerAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to withdraw balance: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ==================== Resource Methods (Stake 2.0) ====================

// FreezeBalanceV2 handles POST /wallet/freezebalancev2
func (h *ContractHandler) FreezeBalanceV2(c *gin.Context) {
	var req struct {
		OwnerAddress  string `json:"owner_address" binding:"required"`
		FrozenBalance int64  `json:"frozen_balance" binding:"required"`
		Resource      string `json:"resource" binding:"required"`
		Visible       bool   `json:"visible" default:"false"`
		PermissionID  int32  `json:"permission_id"`
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

	resourceType := lindapb.ResourceCode_BANDWIDTH
	if req.Resource == "ENERGY" {
		resourceType = lindapb.ResourceCode_ENERGY
	}

	tx, err := h.blockchainClient.FreezeBalanceV2(context.Background(), &lindapb.FreezeBalanceV2Contract{
		OwnerAddress:  []byte(ownerAddress),
		FrozenBalance: req.FrozenBalance,
		Resource:      resourceType,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to freeze balance v2: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// UnfreezeBalanceV2 handles POST /wallet/unfreezebalancev2
func (h *ContractHandler) UnfreezeBalanceV2(c *gin.Context) {
	var req struct {
		OwnerAddress    string `json:"owner_address" binding:"required"`
		UnfreezeBalance int64  `json:"unfreeze_balance" binding:"required"`
		Resource        string `json:"resource" binding:"required"`
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

	resourceType := lindapb.ResourceCode_BANDWIDTH
	if req.Resource == "ENERGY" {
		resourceType = lindapb.ResourceCode_ENERGY
	}

	tx, err := h.blockchainClient.UnfreezeBalanceV2(context.Background(), &lindapb.UnfreezeBalanceV2Contract{
		OwnerAddress:     []byte(ownerAddress),
		UnfreezeBalance: req.UnfreezeBalance,
		Resource:         resourceType,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to unfreeze balance v2: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// WithdrawExpireUnfreeze handles POST /wallet/withdrawexpireunfreeze
func (h *ContractHandler) WithdrawExpireUnfreeze(c *gin.Context) {
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

	tx, err := h.blockchainClient.WithdrawExpireUnfreeze(context.Background(), &lindapb.WithdrawExpireUnfreezeContract{
		OwnerAddress: []byte(ownerAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to withdraw expire unfreeze: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// DelegateResource handles POST /wallet/delegateresource
func (h *ContractHandler) DelegateResource(c *gin.Context) {
	var req struct {
		OwnerAddress    string `json:"owner_address" binding:"required"`
		ReceiverAddress string `json:"receiver_address" binding:"required"`
		Balance         int64  `json:"balance" binding:"required"`
		Resource        string `json:"resource" binding:"required"`
		Lock            bool   `json:"lock"`
		LockPeriod      int64  `json:"lock_period"`
		Visible         bool   `json:"visible" default:"false"`
		PermissionID    int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	receiverAddress := req.ReceiverAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexReceiver, _ := utils.Base58ToHex(receiverAddress)
		ownerAddress = hexOwner
		receiverAddress = hexReceiver
	}

	resourceType := lindapb.ResourceCode_BANDWIDTH
	if req.Resource == "ENERGY" {
		resourceType = lindapb.ResourceCode_ENERGY
	}

	tx, err := h.blockchainClient.DelegateResource(context.Background(), &lindapb.DelegateResourceContract{
		OwnerAddress:    []byte(ownerAddress),
		ReceiverAddress: []byte(receiverAddress),
		Balance:         req.Balance,
		Resource:        resourceType,
		Lock:            req.Lock,
		LockPeriod:      req.LockPeriod,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delegate resource: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// UnDelegateResource handles POST /wallet/undelegateresource
func (h *ContractHandler) UnDelegateResource(c *gin.Context) {
	var req struct {
		OwnerAddress    string `json:"owner_address" binding:"required"`
		ReceiverAddress string `json:"receiver_address" binding:"required"`
		Balance         int64  `json:"balance" binding:"required"`
		Resource        string `json:"resource" binding:"required"`
		Visible         bool   `json:"visible" default:"false"`
		PermissionID    int32  `json:"permission_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ownerAddress := req.OwnerAddress
	receiverAddress := req.ReceiverAddress

	if req.Visible {
		hexOwner, _ := utils.Base58ToHex(ownerAddress)
		hexReceiver, _ := utils.Base58ToHex(receiverAddress)
		ownerAddress = hexOwner
		receiverAddress = hexReceiver
	}

	resourceType := lindapb.ResourceCode_BANDWIDTH
	if req.Resource == "ENERGY" {
		resourceType = lindapb.ResourceCode_ENERGY
	}

	tx, err := h.blockchainClient.UnDelegateResource(context.Background(), &lindapb.UnDelegateResourceContract{
		OwnerAddress:    []byte(ownerAddress),
		ReceiverAddress: []byte(receiverAddress),
		Balance:         req.Balance,
		Resource:        resourceType,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to undelegate resource: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// CancelAllUnfreezeV2 handles POST /wallet/cancelallunfreezev2
func (h *ContractHandler) CancelAllUnfreezeV2(c *gin.Context) {
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

	tx, err := h.blockchainClient.CancelAllUnfreezeV2(context.Background(), &lindapb.CancelAllUnfreezeV2Contract{
		OwnerAddress: []byte(ownerAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to cancel all unfreeze v2: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetAvailableUnfreezeCount handles POST /wallet/getavailableunfreezecount
func (h *ContractHandler) GetAvailableUnfreezeCount(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		Visible      bool   `json:"visible" default:"false"`
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

	count, err := h.blockchainClient.GetAvailableUnfreezeCount(context.Background(), &lindapb.Account{
		Address: []byte(ownerAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get available unfreeze count: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"count": count.Num})
}

// GetCanWithdrawUnfreezeAmount handles POST /wallet/getcanwithdrawunfreezeamount
func (h *ContractHandler) GetCanWithdrawUnfreezeAmount(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		Timestamp    int64  `json:"timestamp" binding:"required"`
		Visible      bool   `json:"visible" default:"false"`
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

	amount, err := h.blockchainClient.GetCanWithdrawUnfreezeAmount(context.Background(), &lindapb.CanWithdrawUnfreezeAmountReq{
		OwnerAddress: []byte(ownerAddress),
		Timestamp:    req.Timestamp,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get can withdraw unfreeze amount: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"amount": amount.Num})
}

// GetDelegatedResource handles POST /wallet/getdelegatedresource
func (h *ContractHandler) GetDelegatedResource(c *gin.Context) {
	var req struct {
		FromAddress string `json:"fromAddress" binding:"required"`
		ToAddress   string `json:"toAddress" binding:"required"`
		Visible     bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	fromAddress := req.FromAddress
	toAddress := req.ToAddress

	if req.Visible {
		hexFrom, _ := utils.Base58ToHex(fromAddress)
		hexTo, _ := utils.Base58ToHex(toAddress)
		fromAddress = hexFrom
		toAddress = hexTo
	}

	resources, err := h.blockchainClient.GetDelegatedResource(context.Background(), &lindapb.DelegatedResourceReq{
		FromAddress: []byte(fromAddress),
		ToAddress:   []byte(toAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get delegated resource: "+err.Error())
		return
	}

	var response []gin.H
	for _, r := range resources.DelegatedResource {
		item := gin.H{
			"from": string(r.From),
			"to":   string(r.To),
		}
		if r.FrozenBalanceForBandwidth > 0 {
			item["frozen_balance_for_bandwidth"] = r.FrozenBalanceForBandwidth
			item["expire_time_for_bandwidth"] = r.ExpireTimeForBandwidth
		}
		if r.FrozenBalanceForEnergy > 0 {
			item["frozen_balance_for_energy"] = r.FrozenBalanceForEnergy
			item["expire_time_for_energy"] = r.ExpireTimeForEnergy
		}
		if req.Visible {
			base58From, _ := utils.HexToBase58(string(r.From))
			base58To, _ := utils.HexToBase58(string(r.To))
			item["from"] = base58From
			item["to"] = base58To
		}
		response = append(response, item)
	}

	utils.RespondWithSuccess(c, response)
}

// GetDelegatedResourceV2 handles POST /wallet/getdelegatedresourcev2
func (h *ContractHandler) GetDelegatedResourceV2(c *gin.Context) {
	var req struct {
		FromAddress string `json:"fromAddress" binding:"required"`
		ToAddress   string `json:"toAddress" binding:"required"`
		Visible     bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	fromAddress := req.FromAddress
	toAddress := req.ToAddress

	if req.Visible {
		hexFrom, _ := utils.Base58ToHex(fromAddress)
		hexTo, _ := utils.Base58ToHex(toAddress)
		fromAddress = hexFrom
		toAddress = hexTo
	}

	resources, err := h.blockchainClient.GetDelegatedResourceV2(context.Background(), &lindapb.DelegatedResourceReq{
		FromAddress: []byte(fromAddress),
		ToAddress:   []byte(toAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get delegated resource v2: "+err.Error())
		return
	}

	var response []gin.H
	for _, r := range resources.DelegatedResource {
		item := gin.H{
			"from": string(r.From),
			"to":   string(r.To),
		}
		if r.FrozenBalanceForBandwidth > 0 {
			item["frozen_balance_for_bandwidth"] = r.FrozenBalanceForBandwidth
			item["expire_time_for_bandwidth"] = r.ExpireTimeForBandwidth
		}
		if r.FrozenBalanceForEnergy > 0 {
			item["frozen_balance_for_energy"] = r.FrozenBalanceForEnergy
			item["expire_time_for_energy"] = r.ExpireTimeForEnergy
		}
		if req.Visible {
			base58From, _ := utils.HexToBase58(string(r.From))
			base58To, _ := utils.HexToBase58(string(r.To))
			item["from"] = base58From
			item["to"] = base58To
		}
		response = append(response, item)
	}

	utils.RespondWithSuccess(c, response)
}

// GetDelegatedResourceAccountIndex handles POST /wallet/getdelegatedresourceaccountindex
func (h *ContractHandler) GetDelegatedResourceAccountIndex(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	address := req.Value
	if req.Visible {
		hexAddr, _ := utils.Base58ToHex(address)
		address = hexAddr
	}

	index, err := h.blockchainClient.GetDelegatedResourceAccountIndex(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get delegated resource account index: "+err.Error())
		return
	}

	response := gin.H{
		"account": string(index.Account),
	}

	if len(index.FromAccounts) > 0 {
		fromAccounts := make([]string, len(index.FromAccounts))
		for i, addr := range index.FromAccounts {
			if req.Visible {
				base58Addr, _ := utils.HexToBase58(addr)
				fromAccounts[i] = base58Addr
			} else {
				fromAccounts[i] = addr
			}
		}
		response["fromAccounts"] = fromAccounts
	}

	if len(index.ToAccounts) > 0 {
		toAccounts := make([]string, len(index.ToAccounts))
		for i, addr := range index.ToAccounts {
			if req.Visible {
				base58Addr, _ := utils.HexToBase58(addr)
				toAccounts[i] = base58Addr
			} else {
				toAccounts[i] = addr
			}
		}
		response["toAccounts"] = toAccounts
	}

	if req.Visible {
		base58Addr, _ := utils.HexToBase58(string(index.Account))
		response["account"] = base58Addr
	}

	utils.RespondWithSuccess(c, response)
}

// GetDelegatedResourceAccountIndexV2 handles POST /wallet/getdelegatedresourceaccountindexv2
func (h *ContractHandler) GetDelegatedResourceAccountIndexV2(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	address := req.Value
	if req.Visible {
		hexAddr, _ := utils.Base58ToHex(address)
		address = hexAddr
	}

	index, err := h.blockchainClient.GetDelegatedResourceAccountIndexV2(context.Background(), &lindapb.Account{
		Address: []byte(address),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get delegated resource account index v2: "+err.Error())
		return
	}

	response := gin.H{
		"account": string(index.Account),
	}

	if len(index.FromAccounts) > 0 {
		fromAccounts := make([]string, len(index.FromAccounts))
		for i, addr := range index.FromAccounts {
			if req.Visible {
				base58Addr, _ := utils.HexToBase58(addr)
				fromAccounts[i] = base58Addr
			} else {
				fromAccounts[i] = addr
			}
		}
		response["fromAccounts"] = fromAccounts
	}

	if len(index.ToAccounts) > 0 {
		toAccounts := make([]string, len(index.ToAccounts))
		for i, addr := range index.ToAccounts {
			if req.Visible {
				base58Addr, _ := utils.HexToBase58(addr)
				toAccounts[i] = base58Addr
			} else {
				toAccounts[i] = addr
			}
		}
		response["toAccounts"] = toAccounts
	}

	if req.Visible {
		base58Addr, _ := utils.HexToBase58(string(index.Account))
		response["account"] = base58Addr
	}

	utils.RespondWithSuccess(c, response)
}

// GetCanDelegatedMaxSize handles POST /wallet/getcandelegatedmaxsize
func (h *ContractHandler) GetCanDelegatedMaxSize(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		Type         int32  `json:"type" binding:"required"`
		Visible      bool   `json:"visible" default:"false"`
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

	maxSize, err := h.blockchainClient.GetCanDelegatedMaxSize(context.Background(), &lindapb.CanDelegatedMaxSizeReq{
		OwnerAddress: []byte(ownerAddress),
		Type:         req.Type,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get can delegated max size: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"max_size": maxSize.Num})
}

// ==================== Exchange Methods ====================

// ExchangeCreate handles POST /wallet/exchangecreate
func (h *ContractHandler) ExchangeCreate(c *gin.Context) {
	var req struct {
		OwnerAddress      string `json:"owner_address" binding:"required"`
		FirstTokenID      string `json:"first_token_id" binding:"required"`
		FirstTokenBalance int64  `json:"first_token_balance" binding:"required"`
		SecondTokenID     string `json:"second_token_id" binding:"required"`
		SecondTokenBalance int64 `json:"second_token_balance" binding:"required"`
		Visible           bool   `json:"visible" default:"false"`
		PermissionID      int32  `json:"permission_id"`
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

	tx, err := h.blockchainClient.ExchangeCreate(context.Background(), &lindapb.ExchangeCreateContract{
		OwnerAddress:        []byte(ownerAddress),
		FirstTokenId:        []byte(req.FirstTokenID),
		FirstTokenBalance:   req.FirstTokenBalance,
		SecondTokenId:       []byte(req.SecondTokenID),
		SecondTokenBalance:  req.SecondTokenBalance,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create exchange: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ExchangeInject handles POST /wallet/exchangeinject
func (h *ContractHandler) ExchangeInject(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		ExchangeID   int64  `json:"exchange_id" binding:"required"`
		TokenID      string `json:"token_id" binding:"required"`
		Quant        int64  `json:"quant" binding:"required"`
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

	tx, err := h.blockchainClient.ExchangeInject(context.Background(), &lindapb.ExchangeInjectContract{
		OwnerAddress: []byte(ownerAddress),
		ExchangeId:   req.ExchangeID,
		TokenId:      []byte(req.TokenID),
		Quant:        req.Quant,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to inject exchange: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ExchangeWithdraw handles POST /wallet/exchangewithdraw
func (h *ContractHandler) ExchangeWithdraw(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		ExchangeID   int64  `json:"exchange_id" binding:"required"`
		TokenID      string `json:"token_id" binding:"required"`
		Quant        int64  `json:"quant" binding:"required"`
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

	tx, err := h.blockchainClient.ExchangeWithdraw(context.Background(), &lindapb.ExchangeWithdrawContract{
		OwnerAddress: []byte(ownerAddress),
		ExchangeId:   req.ExchangeID,
		TokenId:      []byte(req.TokenID),
		Quant:        req.Quant,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to withdraw from exchange: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ExchangeTransaction handles POST /wallet/exchangetransaction
func (h *ContractHandler) ExchangeTransaction(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		ExchangeID   int64  `json:"exchange_id" binding:"required"`
		TokenID      string `json:"token_id" binding:"required"`
		Quant        int64  `json:"quant" binding:"required"`
		Expected     int64  `json:"expected" binding:"required"`
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

	tx, err := h.blockchainClient.ExchangeTransaction(context.Background(), &lindapb.ExchangeTransactionContract{
		OwnerAddress: []byte(ownerAddress),
		ExchangeId:   req.ExchangeID,
		TokenId:      []byte(req.TokenID),
		Quant:        req.Quant,
		Expected:     req.Expected,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to exchange transaction: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetExchangeById handles POST /wallet/getexchangebyid
func (h *ContractHandler) GetExchangeById(c *gin.Context) {
	var req struct {
		ExchangeID int64 `json:"exchange_id" binding:"required"`
		Visible    bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	exchange, err := h.blockchainClient.GetExchangeById(context.Background(), &lindapb.NumberMessage{
		Num: req.ExchangeID,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get exchange: "+err.Error())
		return
	}

	response := convertExchangeToResponse(exchange, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// ListExchanges handles GET /wallet/listexchanges
func (h *ContractHandler) ListExchanges(c *gin.Context) {
	visible := c.DefaultQuery("visible", "false") == "true"

	exchanges, err := h.blockchainClient.ListExchanges(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list exchanges: "+err.Error())
		return
	}

	var response []gin.H
	for _, exchange := range exchanges.Exchanges {
		response = append(response, convertExchangeToResponse(exchange, visible))
	}

	utils.RespondWithSuccess(c, response)
}

// GetPaginatedExchangeList handles POST /wallet/getpaginatedexchangelist
func (h *ContractHandler) GetPaginatedExchangeList(c *gin.Context) {
	var req struct {
		Offset  int64 `json:"offset"`
		Limit   int64 `json:"limit" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	exchanges, err := h.blockchainClient.GetPaginatedExchangeList(context.Background(), &lindapb.PaginatedMessage{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get paginated exchanges: "+err.Error())
		return
	}

	var response []gin.H
	for _, exchange := range exchanges.Exchanges {
		response = append(response, convertExchangeToResponse(exchange, req.Visible))
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Market Methods (DEX) ====================

// MarketSellAsset handles POST /wallet/marketsellasset
func (h *ContractHandler) MarketSellAsset(c *gin.Context) {
	var req struct {
		OwnerAddress   string `json:"owner_address" binding:"required"`
		SellTokenID    string `json:"sell_token_id" binding:"required"`
		SellTokenValue int64  `json:"sell_token_value" binding:"required"`
		BuyTokenID     string `json:"buy_token_id" binding:"required"`
		BuyTokenValue  int64  `json:"buy_token_value" binding:"required"`
		Visible        bool   `json:"visible" default:"false"`
		PermissionID   int32  `json:"permission_id"`
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

	tx, err := h.blockchainClient.MarketSellAsset(context.Background(), &lindapb.MarketSellAssetContract{
		OwnerAddress:   []byte(ownerAddress),
		SellTokenId:    []byte(req.SellTokenID),
		SellTokenValue: req.SellTokenValue,
		BuyTokenId:     []byte(req.BuyTokenID),
		BuyTokenValue:  req.BuyTokenValue,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to market sell asset: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// MarketCancelOrder handles POST /wallet/marketcancelorder
func (h *ContractHandler) MarketCancelOrder(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		OrderID      string `json:"order_id" binding:"required"`
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

	tx, err := h.blockchainClient.MarketCancelOrder(context.Background(), &lindapb.MarketCancelOrderContract{
		OwnerAddress: []byte(ownerAddress),
		OrderId:      []byte(req.OrderID),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to cancel market order: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetMarketOrderByAccount handles POST /wallet/getmarketorderbyaccount
func (h *ContractHandler) GetMarketOrderByAccount(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		Visible      bool   `json:"visible" default:"false"`
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

	orders, err := h.blockchainClient.GetMarketOrderByAccount(context.Background(), &lindapb.MarketOrderReq{
		OwnerAddress: []byte(ownerAddress),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get market orders: "+err.Error())
		return
	}

	var response []gin.H
	for _, order := range orders.Orders {
		response = append(response, convertMarketOrderToResponse(order, req.Visible))
	}

	utils.RespondWithSuccess(c, response)
}

// GetMarketOrderById handles POST /wallet/getmarketorderbyid
func (h *ContractHandler) GetMarketOrderById(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	order, err := h.blockchainClient.GetMarketOrderById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.OrderID),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get market order: "+err.Error())
		return
	}

	response := convertMarketOrderToResponse(order, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetMarketPriceByPair handles POST /wallet/getmarketpricebypair
func (h *ContractHandler) GetMarketPriceByPair(c *gin.Context) {
	var req struct {
		SellTokenID string `json:"sell_token_id" binding:"required"`
		BuyTokenID  string `json:"buy_token_id" binding:"required"`
		Visible     bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	prices, err := h.blockchainClient.GetMarketPriceByPair(context.Background(), &lindapb.MarketPriceReq{
		SellTokenId: []byte(req.SellTokenID),
		BuyTokenId:  []byte(req.BuyTokenID),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get market price: "+err.Error())
		return
	}

	var response []gin.H
	for _, price := range prices.Prices {
		response = append(response, gin.H{
			"sell_token_id": string(price.SellTokenId),
			"buy_token_id":  string(price.BuyTokenId),
			"price":         price.Price,
		})
	}

	utils.RespondWithSuccess(c, response)
}

// GetMarketOrderListByPair handles POST /wallet/getmarketorderlistbypair
func (h *ContractHandler) GetMarketOrderListByPair(c *gin.Context) {
	var req struct {
		SellTokenID string `json:"sell_token_id" binding:"required"`
		BuyTokenID  string `json:"buy_token_id" binding:"required"`
		Visible     bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	orders, err := h.blockchainClient.GetMarketOrderListByPair(context.Background(), &lindapb.MarketOrderListReq{
		SellTokenId: []byte(req.SellTokenID),
		BuyTokenId:  []byte(req.BuyTokenID),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get market orders: "+err.Error())
		return
	}

	var response []gin.H
	for _, order := range orders.Orders {
		response = append(response, convertMarketOrderToResponse(order, req.Visible))
	}

	utils.RespondWithSuccess(c, response)
}

// GetMarketPairList handles GET /wallet/getmarketpairlist
func (h *ContractHandler) GetMarketPairList(c *gin.Context) {
	visible := c.DefaultQuery("visible", "false") == "true"

	pairs, err := h.blockchainClient.GetMarketPairList(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get market pairs: "+err.Error())
		return
	}

	var response []gin.H
	for _, pair := range pairs.Pairs {
		response = append(response, gin.H{
			"sell_token_id": string(pair.SellTokenId),
			"buy_token_id":  string(pair.BuyTokenId),
		})
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Contract Logs (Event Service) ====================

// GetContractLogs handles GET /v1/contractlogs
func (h *ContractHandler) GetContractLogs(c *gin.Context) {
	var req struct {
		Limit int    `form:"limit" default:"20"`
		Sort  string `form:"sort" default:"-timestamp"`
		Start int    `form:"start" default:"0"`
		Block int64  `form:"block" default:"0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	// Get logs from database
	// This would be implemented with the event repository

	utils.RespondWithV1Success(c, []interface{}{}, req.Limit, "")
}

// GetContractLogsByTransactionId handles GET /v1/contractlogs/transaction/{transaction_id}
func (h *ContractHandler) GetContractLogsByTransactionId(c *gin.Context) {
	txID := c.Param("transaction_id")
	if txID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Transaction ID is required")
		return
	}

	// Get logs from database
	// This would be implemented with the event repository

	utils.RespondWithV1Success(c, []interface{}{}, 0, "")
}

// GetContractLogsByContractAddress handles GET /v1/contractlogs/contract/{contract_address}
func (h *ContractHandler) GetContractLogsByContractAddress(c *gin.Context) {
	contractAddr := c.Param("contract_address")
	if contractAddr == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	// Get logs from database
	// This would be implemented with the event repository

	utils.RespondWithV1Success(c, []interface{}{}, 0, "")
}

// GetContractWithAbi handles POST /v1/contract/transaction/{transaction_id}
func (h *ContractHandler) GetContractWithAbi(c *gin.Context) {
	txID := c.Param("transaction_id")
	if txID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Transaction ID is required")
		return
	}

	var req struct {
		ABI json.RawMessage `json:"abi"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get transaction and parse logs with ABI
	// This would be implemented with the event service

	utils.RespondWithSuccess(c, gin.H{
		"contractAddress": "",
		"abi":             req.ABI,
		"logs":            []interface{}{},
	})
}

// GetContractByAddressWithAbi handles POST /v1/contract/contractAddress/{contract_address}
func (h *ContractHandler) GetContractByAddressWithAbi(c *gin.Context) {
	contractAddr := c.Param("contract_address")
	if contractAddr == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Contract address is required")
		return
	}

	var req struct {
		ABI json.RawMessage `json:"abi"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get contract and parse logs with ABI
	// This would be implemented with the event service

	utils.RespondWithSuccess(c, gin.H{
		"contractAddress": contractAddr,
		"abi":             req.ABI,
		"logs":            []interface{}{},
	})
}

// GetContractAccountHistory handles GET /api/contract_account_history
func (h *ContractHandler) GetContractAccountHistory(c *gin.Context) {
	var req struct {
		Contract string `form:"contract" binding:"required"`
		Address  string `form:"address" binding:"required"`
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

	// Get history from database
	// This would be implemented with the transaction repository

	utils.RespondWithSuccess(c, gin.H{
		"history": []interface{}{},
		"total":   0,
	})
}

// GetSmartContractTriggersBatch handles GET /api/contracts/smart-contract-triggers-batch
func (h *ContractHandler) GetSmartContractTriggersBatch(c *gin.Context) {
	var req struct {
		Contracts []string `form:"contracts" binding:"required"`
		From      int64    `form:"from"`
		To        int64    `form:"to"`
		Limit     int      `form:"limit" default:"20"`
		Start     int      `form:"start" default:"0"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	// Get triggers from database
	// This would be implemented with the transaction repository

	utils.RespondWithSuccess(c, gin.H{
		"triggers": []interface{}{},
		"total":    0,
	})
}

// ==================== Helper Functions ====================

func convertSmartContractToResponse(contract *lindapb.SmartContract, visible bool) gin.H {
	resp := gin.H{
		"origin_address":               string(contract.OriginAddress),
		"contract_address":              string(contract.ContractAddress),
		"abi":                           contract.Abi,
		"bytecode":                      string(contract.Bytecode),
		"name":                          contract.Name,
		"consume_user_resource_percent": contract.ConsumeUserResourcePercent,
		"origin_energy_limit":           contract.OriginEnergyLimit,
		"code_hash":                     string(contract.CodeHash),
	}

	if visible {
		base58Origin, _ := utils.HexToBase58(string(contract.OriginAddress))
		base58Contract, _ := utils.HexToBase58(string(contract.ContractAddress))
		resp["origin_address"] = base58Origin
		resp["contract_address"] = base58Contract
	}

	return resp
}

func convertExchangeToResponse(exchange *lindapb.Exchange, visible bool) gin.H {
	resp := gin.H{
		"exchange_id":            exchange.ExchangeId,
		"creator_address":        string(exchange.CreatorAddress),
		"create_time":            exchange.CreateTime,
		"first_token_id":         string(exchange.FirstTokenId),
		"first_token_balance":    exchange.FirstTokenBalance,
		"second_token_id":        string(exchange.SecondTokenId),
		"second_token_balance":   exchange.SecondTokenBalance,
	}

	if visible {
		base58Creator, _ := utils.HexToBase58(string(exchange.CreatorAddress))
		resp["creator_address"] = base58Creator
	}

	return resp
}

func convertMarketOrderToResponse(order *lindapb.MarketOrder, visible bool) gin.H {
	orderStatus := "PENDING"
	switch order.OrderStatus {
	case 1:
		orderStatus = "CANCELLED"
	case 2:
		orderStatus = "COMPLETED"
	}

	resp := gin.H{
		"order_id":         string(order.OrderId),
		"owner_address":    string(order.OwnerAddress),
		"create_time":      order.CreateTime,
		"sell_token_id":    string(order.SellTokenId),
		"sell_token_value": order.SellTokenValue,
		"buy_token_id":     string(order.BuyTokenId),
		"buy_token_value":  order.BuyTokenValue,
		"order_status":     orderStatus,
	}

	if visible {
		base58Owner, _ := utils.HexToBase58(string(order.OwnerAddress))
		resp["owner_address"] = base58Owner
	}

	return resp
}