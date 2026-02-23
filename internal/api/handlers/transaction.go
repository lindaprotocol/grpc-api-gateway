package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/blockchain"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/storage/repository"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type TransactionHandler struct {
	blockchainClient *blockchain.Client
	txRepo           *repository.TransactionRepository
}

func NewTransactionHandler(client *blockchain.Client, txRepo *repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		blockchainClient: client,
		txRepo:           txRepo,
	}
}

// ==================== Wallet Service Transaction Methods ====================

// CreateTransaction handles POST /wallet/createtransaction
// Creates a LIND transfer transaction
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req struct {
		OwnerAddress string `json:"owner_address" binding:"required"`
		ToAddress    string `json:"to_address" binding:"required"`
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
		hexOwner, err := utils.Base58ToHex(ownerAddress)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid owner address format")
			return
		}
		hexTo, err := utils.Base58ToHex(toAddress)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid to address format")
			return
		}
		ownerAddress = hexOwner
		toAddress = hexTo
	}

	tx, err := h.blockchainClient.CreateTransaction(context.Background(), &lindapb.TransferContract{
		OwnerAddress: []byte(ownerAddress),
		ToAddress:    []byte(toAddress),
		Amount:       req.Amount,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create transaction: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// BroadcastTransaction handles POST /wallet/broadcasttransaction
// Broadcasts a signed transaction
func (h *TransactionHandler) BroadcastTransaction(c *gin.Context) {
	var req lindapb.Transaction
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.blockchainClient.BroadcastTransaction(context.Background(), &req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to broadcast: "+err.Error())
		return
	}

	response := gin.H{
		"result": resp.Result,
		"code":   resp.Code.String(),
		"txid":   string(resp.Txid),
	}

	if !resp.Result && resp.Message != nil {
		response["message"] = string(resp.Message)
	}

	utils.RespondWithSuccess(c, response)
}

// BroadcastHex handles POST /wallet/broadcasthex
// Broadcasts a transaction from hex string
func (h *TransactionHandler) BroadcastHex(c *gin.Context) {
	var req struct {
		Transaction string `json:"transaction" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Decode hex
	txBytes, err := hex.DecodeString(req.Transaction)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid hex format")
		return
	}

	// Parse transaction
	var tx lindapb.Transaction
	if err := tx.Unmarshal(txBytes); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Failed to parse transaction: "+err.Error())
		return
	}

	resp, err := h.blockchainClient.BroadcastTransaction(context.Background(), &tx)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to broadcast: "+err.Error())
		return
	}

	response := gin.H{
		"result": resp.Result,
		"code":   resp.Code.String(),
		"txid":   string(resp.Txid),
	}

	if !resp.Result && resp.Message != nil {
		response["message"] = string(resp.Message)
	}

	utils.RespondWithSuccess(c, response)
}

// GetTransactionById handles POST /wallet/gettransactionbyid
// Returns transaction by ID
func (h *TransactionHandler) GetTransactionById(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	tx, err := h.blockchainClient.GetTransactionById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transaction: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetTransactionInfoById handles POST /wallet/gettransactioninfobyid
// Returns transaction info (fee, block, logs)
func (h *TransactionHandler) GetTransactionInfoById(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	info, err := h.blockchainClient.GetTransactionInfoById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transaction info: "+err.Error())
		return
	}

	response := convertTransactionInfoToResponse(info, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetTransactionReceiptById handles POST /wallet/gettransactionreceiptbyid
// Returns transaction receipt
func (h *TransactionHandler) GetTransactionReceiptById(c *gin.Context) {
	var req struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// This is similar to GetTransactionInfoById but returns receipt format
	info, err := h.blockchainClient.GetTransactionInfoById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transaction receipt: "+err.Error())
		return
	}

	receipt := gin.H{
		"id":               string(info.Id),
		"fee":              info.Fee,
		"blockNumber":      info.BlockNumber,
		"blockTimeStamp":   info.BlockTimeStamp,
		"receipt":          info.Receipt,
		"contract_address": string(info.ContractAddress),
		"result":           info.Result,
	}

	if info.ResMessage != nil {
		receipt["resMessage"] = string(info.ResMessage)
	}

	utils.RespondWithSuccess(c, receipt)
}

// GetTransactionCountByBlockNum handles POST /wallet/gettransactioncountbyblocknum
// Returns transaction count in a block
func (h *TransactionHandler) GetTransactionCountByBlockNum(c *gin.Context) {
	var req struct {
		Num int64 `json:"num" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	count, err := h.blockchainClient.GetTransactionCountByBlockNum(context.Background(), &lindapb.NumberMessage{
		Num: req.Num,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transaction count: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{"count": count.Num})
}

// GetTransactionInfoByBlockNum handles POST /wallet/gettransactioninfobyblocknum
// Returns all transaction infos in a block
func (h *TransactionHandler) GetTransactionInfoByBlockNum(c *gin.Context) {
	var req struct {
		Num     int64 `json:"num" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	infos, err := h.blockchainClient.GetTransactionInfoByBlockNum(context.Background(), &lindapb.NumberMessage{
		Num: req.Num,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transaction infos: "+err.Error())
		return
	}

	var responses []*models.TransactionInfoResponse
	for _, info := range infos.TransactionInfo {
		responses = append(responses, convertTransactionInfoToResponse(info, req.Visible))
	}

	utils.RespondWithSuccess(c, responses)
}

// GetTransactionSign handles POST /wallet/gettransactionsign
// Signs a transaction (warning: private key required)
func (h *TransactionHandler) GetTransactionSign(c *gin.Context) {
	var req struct {
		Transaction json.RawMessage `json:"transaction" binding:"required"`
		PrivateKey  string          `json:"privateKey" binding:"required"`
		Visible     bool            `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Parse transaction
	var tx lindapb.Transaction
	if err := json.Unmarshal(req.Transaction, &tx); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid transaction format")
		return
	}

	// Sign transaction
	signedTx, err := h.blockchainClient.GetTransactionSign(context.Background(), &lindapb.TransactionSign{
		Transaction: &tx,
		PrivateKey:  []byte(req.PrivateKey),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to sign transaction: "+err.Error())
		return
	}

	response := convertTransactionToResponse(signedTx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetTransactionSignWeight handles POST /wallet/getsignweight
// Returns sign weight of a transaction
func (h *TransactionHandler) GetTransactionSignWeight(c *gin.Context) {
	var req struct {
		Transaction json.RawMessage `json:"transaction" binding:"required"`
		Visible     bool            `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	var tx lindapb.Transaction
	if err := json.Unmarshal(req.Transaction, &tx); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid transaction format")
		return
	}

	resp, err := h.blockchainClient.GetTransactionSignWeight(context.Background(), &tx)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get sign weight: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, resp)
}

// GetTransactionApprovedList handles POST /wallet/getapprovedlist
// Returns list of addresses that signed the transaction
func (h *TransactionHandler) GetTransactionApprovedList(c *gin.Context) {
	var req struct {
		Signature   []string        `json:"signature" binding:"required"`
		RawData     json.RawMessage `json:"raw_data" binding:"required"`
		Visible     bool            `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Build transaction for verification
	tx := &lindapb.Transaction{
		Signature: req.Signature,
	}

	// Parse raw_data
	var rawData lindapb.TransactionRaw
	if err := json.Unmarshal(req.RawData, &rawData); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid raw_data format")
		return
	}
	tx.RawData = &rawData

	resp, err := h.blockchainClient.GetTransactionApprovedList(context.Background(), tx)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get approved list: "+err.Error())
		return
	}

	// Convert addresses if visible
	approvedList := resp.ApprovedList
	if req.Visible {
		for i, addr := range approvedList {
			base58Addr, _ := utils.HexToBase58(addr)
			approvedList[i] = base58Addr
		}
	}

	response := gin.H{
		"approved_list": approvedList,
		"result":        resp.Result,
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Solidity Node Transaction Methods ====================

// GetTransactionByIdSolidity handles POST /walletsolidity/gettransactionbyid
// Returns confirmed transaction by ID
func (h *TransactionHandler) GetTransactionByIdSolidity(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	tx, err := h.blockchainClient.GetTransactionByIdSolidity(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transaction: "+err.Error())
		return
	}

	response := convertTransactionToResponse(tx, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetTransactionInfoByIdSolidity handles POST /walletsolidity/gettransactioninfobyid
// Returns confirmed transaction info
func (h *TransactionHandler) GetTransactionInfoByIdSolidity(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	info, err := h.blockchainClient.GetTransactionInfoByIdSolidity(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transaction info: "+err.Error())
		return
	}

	response := convertTransactionInfoToResponse(info, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetTransactionInfoByBlockNumSolidity handles POST /walletsolidity/gettransactioninfobyblocknum
// Returns all transaction infos in a confirmed block
func (h *TransactionHandler) GetTransactionInfoByBlockNumSolidity(c *gin.Context) {
	var req struct {
		Num     int64 `json:"num" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	infos, err := h.blockchainClient.GetTransactionInfoByBlockNumSolidity(context.Background(), &lindapb.NumberMessage{
		Num: req.Num,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transaction infos: "+err.Error())
		return
	}

	var responses []*models.TransactionInfoResponse
	for _, info := range infos.TransactionInfo {
		responses = append(responses, convertTransactionInfoToResponse(info, req.Visible))
	}

	utils.RespondWithSuccess(c, responses)
}

// ==================== Event Query Transaction Methods ====================

// GetTransactionsEvent handles GET /v1/transactions
// Returns paginated transactions
func (h *TransactionHandler) GetTransactionsEvent(c *gin.Context) {
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

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	// Get from database
	txs, total, err := h.txRepo.GetTransactions(req.Block, req.Start, req.Limit, req.Sort)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transactions: "+err.Error())
		return
	}

	// Convert to event response format
	var data []*models.EventTransactionResponse
	for _, tx := range txs {
		data = append(data, &models.EventTransactionResponse{
			ID:             tx.TxID,
			BlockNumber:    tx.BlockNumber,
			BlockTimestamp: tx.BlockTimestamp,
			From:           tx.FromAddress,
			To:             tx.ToAddress,
			Value:          strconv.FormatInt(tx.Amount, 10),
			Fee:            tx.Fee,
		})
	}

	utils.RespondWithSuccess(c, gin.H{
		"data":    data,
		"success": true,
		"meta": gin.H{
			"at":        utils.GetCurrentTimestamp(),
			"page_size": req.Limit,
		},
	})
}

// GetTransactionByHashEvent handles GET /v1/transactions/{hash}
// Returns transaction by hash from event service
func (h *TransactionHandler) GetTransactionByHashEvent(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Hash is required")
		return
	}

	// Get from blockchain
	tx, err := h.blockchainClient.GetTransactionById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(hash),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Transaction not found")
		return
	}

	// Get additional info
	info, _ := h.blockchainClient.GetTransactionInfoById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(hash),
	})

	response := &models.EventTransactionResponse{
		ID:          string(tx.TxID),
		BlockNumber: 0,
	}

	// Extract from address and to address from contract
	if len(tx.RawData.Contract) > 0 {
		// Parse contract based on type
		// This is simplified - actual implementation would handle different contract types
		var contract map[string]interface{}
		if err := json.Unmarshal(tx.RawData.Contract[0].Parameter.Value, &contract); err == nil {
			if owner, ok := contract["owner_address"]; ok {
				response.From = string(owner.([]uint8))
			}
			if to, ok := contract["to_address"]; ok {
				response.To = string(to.([]uint8))
			}
			if amount, ok := contract["amount"]; ok {
				if amt, ok := amount.(float64); ok {
					response.Value = strconv.FormatInt(int64(amt), 10)
				}
			}
		}
	}

	if info != nil {
		response.BlockNumber = info.BlockNumber
		response.BlockTimestamp = info.BlockTimeStamp
		response.Fee = info.Fee
	}

	utils.RespondWithSuccess(c, gin.H{
		"data":    []*models.EventTransactionResponse{response},
		"success": true,
		"meta": gin.H{
			"at": utils.GetCurrentTimestamp(),
		},
	})
}

// ==================== Lindascan Custom Transaction Methods ====================

// GetTransactionsV2 handles GET /api/transaction
// Returns paginated transactions for lindascan
func (h *TransactionHandler) GetTransactionsV2(c *gin.Context) {
	var req struct {
		Hash    string `form:"hash"`
		Block   int64  `form:"block"`
		Address string `form:"address"`
		Limit   int    `form:"limit" default:"20"`
		Start   int    `form:"start" default:"0"`
		Sort    string `form:"sort" default:"-timestamp"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	var txs []*models.TransactionResponse
	var total int64
	var err error

	if req.Hash != "" {
		// Get by hash
		tx, err := h.blockchainClient.GetTransactionById(context.Background(), &lindapb.BytesMessage{
			Value: []byte(req.Hash),
		})
		if err == nil {
			txs = append(txs, convertTransactionToResponse(tx, true))
			total = 1
		}
	} else if req.Block > 0 {
		// Get by block
		block, err := h.blockchainClient.GetBlockByNum(context.Background(), &lindapb.NumberMessage{
			Num: req.Block,
		})
		if err == nil {
			for _, tx := range block.Transactions {
				txs = append(txs, convertTransactionToResponse(tx, true))
			}
			total = int64(len(txs))
		}
	} else if req.Address != "" {
		// Get by address
		txs, total, err = h.txRepo.GetTransactionsByAddress(req.Address, req.Start, req.Limit, req.Sort)
	} else {
		// Get paginated
		txs, total, err = h.txRepo.GetTransactions(0, req.Start, req.Limit, req.Sort)
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get transactions: "+err.Error())
		return
	}

	response := gin.H{
		"transactions": txs,
		"total":        total,
	}

	utils.RespondWithSuccess(c, response)
}

// GetInternalTransactions handles GET /api/internal-transaction
// Returns internal transactions
func (h *TransactionHandler) GetInternalTransactions(c *gin.Context) {
	var req struct {
		TxHash  string `form:"tx_hash"`
		Address string `form:"address"`
		Limit   int    `form:"limit" default:"20"`
		Start   int    `form:"start" default:"0"`
		Sort    string `form:"sort" default:"-timestamp"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	var internalTxs []*models.InternalTransaction
	var total int64

	if req.TxHash != "" {
		// Get internal txs for a specific transaction
		info, err := h.blockchainClient.GetTransactionInfoById(context.Background(), &lindapb.BytesMessage{
			Value: []byte(req.TxHash),
		})
		if err == nil {
			internalTxs = info.InternalTransactions
			total = int64(len(internalTxs))
		}
	} else if req.Address != "" {
		// Get from database
		internalTxs, total, err = h.txRepo.GetInternalTransactionsByAddress(req.Address, req.Start, req.Limit, req.Sort)
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get internal transactions: "+err.Error())
		return
	}

	response := gin.H{
		"internal_transactions": internalTxs,
		"total":                 total,
	}

	utils.RespondWithSuccess(c, response)
}

// GetContractTransactions handles GET /api/contracts/transaction
// Returns transactions for a contract
func (h *TransactionHandler) GetContractTransactions(c *gin.Context) {
	var req struct {
		Contract string `form:"contract" binding:"required"`
		Limit    int    `form:"limit" default:"20"`
		Start    int    `form:"start" default:"0"`
		Sort     string `form:"sort" default:"-timestamp"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	txs, total, err := h.txRepo.GetTransactionsByContract(req.Contract, req.Start, req.Limit, req.Sort)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get contract transactions: "+err.Error())
		return
	}

	response := gin.H{
		"transactions": txs,
		"total":        total,
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Helper Functions ====================

func convertTransactionInfoToResponse(info *lindapb.TransactionInfo, visible bool) *models.TransactionInfoResponse {
	resp := &models.TransactionInfoResponse{
		ID:                        string(info.Id),
		Fee:                       info.Fee,
		BlockNumber:               info.BlockNumber,
		BlockTimeStamp:            info.BlockTimeStamp,
		ContractResult:            info.ContractResult,
		ContractAddress:           string(info.ContractAddress),
		Receipt:                   convertReceiptToResponse(info.Receipt),
		Result:                    int(info.Result),
		AssetIssueID:              info.AssetIssueID,
		WithdrawAmount:            info.WithdrawAmount,
		UnfreezeAmount:            info.UnfreezeAmount,
		InternalTransactions:      info.InternalTransactions,
		WithdrawExpireAmount:      info.WithdrawExpireAmount,
		CancelUnfreezeV2Amount:    info.CancelUnfreezeV2Amount,
		ExchangeReceivedAmount:    info.ExchangeReceivedAmount,
		ExchangeInjectAnotherAmount: info.ExchangeInjectAnotherAmount,
		ExchangeWithdrawAnotherAmount: info.ExchangeWithdrawAnotherAmount,
		ExchangeID:                info.ExchangeId,
		ShieldedTransactionFee:    info.ShieldedTransactionFee,
	}

	// Convert logs
	for _, log := range info.Log {
		resp.Log = append(resp.Log, &models.EventLog{
			Address: string(log.Address),
			Topics:  log.Topics,
			Data:    string(log.Data),
		})
	}

	if info.ResMessage != nil {
		resp.ResMessage = string(info.ResMessage)
	}

	if visible {
		base58Addr, _ := utils.HexToBase58(string(info.ContractAddress))
		resp.ContractAddress = base58Addr
	}

	return resp
}

func convertReceiptToResponse(receipt *lindapb.ResourceReceipt) *models.ResourceReceipt {
	if receipt == nil {
		return nil
	}

	return &models.ResourceReceipt{
		EnergyUsage:        receipt.EnergyUsage,
		EnergyFee:          receipt.EnergyFee,
		OriginEnergyUsage:  receipt.OriginEnergyUsage,
		EnergyUsageTotal:   receipt.EnergyUsageTotal,
		NetUsage:           receipt.NetUsage,
		NetFee:             receipt.NetFee,
		Result:             receipt.Result.String(),
		EnergyPenaltyTotal: receipt.EnergyPenaltyTotal,
	}
}