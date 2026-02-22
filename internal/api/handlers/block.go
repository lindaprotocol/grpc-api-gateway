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

type BlockHandler struct {
	blockchainClient *blockchain.Client
	blockRepo        *repository.BlockRepository
}

func NewBlockHandler(client *blockchain.Client, blockRepo *repository.BlockRepository) *BlockHandler {
	return &BlockHandler{
		blockchainClient: client,
		blockRepo:        blockRepo,
	}
}

// ==================== Wallet Service Block Methods ====================

// GetNowBlock handles POST /wallet/getnowblock
// Returns the most recent block
func (h *BlockHandler) GetNowBlock(c *gin.Context) {
	var req struct {
		Visible bool `json:"visible" form:"visible" default:"false"`
	}

	if err := c.ShouldBind(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	block, err := h.blockchainClient.GetNowBlock(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get now block: "+err.Error())
		return
	}

	response := convertBlockToResponse(block, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetBlockByNum handles POST /wallet/getblockbynum
// Returns block by number
func (h *BlockHandler) GetBlockByNum(c *gin.Context) {
	var req struct {
		Num     int64 `json:"num" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	block, err := h.blockchainClient.GetBlockByNum(context.Background(), &lindapb.NumberMessage{
		Num: req.Num,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get block: "+err.Error())
		return
	}

	response := convertBlockToResponse(block, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetBlockById handles POST /wallet/getblockbyid
// Returns block by hash
func (h *BlockHandler) GetBlockById(c *gin.Context) {
	var req struct {
		Value   string `json:"value" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	block, err := h.blockchainClient.GetBlockById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(req.Value),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get block: "+err.Error())
		return
	}

	response := convertBlockToResponse(block, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetBlockByLimitNext handles POST /wallet/getblockbylimitnext
// Returns blocks in range
func (h *BlockHandler) GetBlockByLimitNext(c *gin.Context) {
	var req struct {
		StartNum int64 `json:"startNum" binding:"required"`
		EndNum   int64 `json:"endNum" binding:"required"`
		Visible  bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	blocks, err := h.blockchainClient.GetBlockByLimitNext(context.Background(), &lindapb.BlockLimit{
		StartNum: req.StartNum,
		EndNum:   req.EndNum,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get blocks: "+err.Error())
		return
	}

	var blockResponses []*models.BlockResponse
	for _, block := range blocks.Block {
		blockResponses = append(blockResponses, convertBlockToResponse(block, req.Visible))
	}

	utils.RespondWithSuccess(c, gin.H{"block": blockResponses})
}

// GetBlockByLatestNum handles POST /wallet/getblockbylatestnum
// Returns latest N blocks
func (h *BlockHandler) GetBlockByLatestNum(c *gin.Context) {
	var req struct {
		Num     int64 `json:"num" binding:"required"`
		Visible bool  `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	blocks, err := h.blockchainClient.GetBlockByLatestNum(context.Background(), &lindapb.NumberMessage{
		Num: req.Num,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get blocks: "+err.Error())
		return
	}

	var blockResponses []*models.BlockResponse
	for _, block := range blocks.Block {
		blockResponses = append(blockResponses, convertBlockToResponse(block, req.Visible))
	}

	utils.RespondWithSuccess(c, gin.H{"block": blockResponses})
}

// GetBlock handles POST /wallet/getblock
// Returns block by id or number
func (h *BlockHandler) GetBlock(c *gin.Context) {
	var req struct {
		IdOrNum string `json:"id_or_num"`
		Detail  bool   `json:"detail" default:"false"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	var block *lindapb.BlockExtention
	var err error

	if req.IdOrNum == "" {
		// Get latest block
		b, _ := h.blockchainClient.GetNowBlock(context.Background(), &lindapb.EmptyMessage{})
		block = &lindapb.BlockExtention{
			BlockHeader:  b.BlockHeader,
			Transactions: b.Transactions,
		}
	} else {
		// Try as number first
		if num, err := strconv.ParseInt(req.IdOrNum, 10, 64); err == nil {
			b, err := h.blockchainClient.GetBlockByNum(context.Background(), &lindapb.NumberMessage{Num: num})
			if err == nil {
				block = &lindapb.BlockExtention{
					BlockHeader:  b.BlockHeader,
					Transactions: b.Transactions,
				}
			}
		} else {
			// Try as hash
			b, err := h.blockchainClient.GetBlockById(context.Background(), &lindapb.BytesMessage{Value: []byte(req.IdOrNum)})
			if err == nil {
				block = &lindapb.BlockExtention{
					BlockHeader:  b.BlockHeader,
					Transactions: b.Transactions,
				}
			}
		}
	}

	if err != nil || block == nil {
		utils.RespondWithError(c, http.StatusNotFound, "Block not found")
		return
	}

	response := convertBlockExtentionToResponse(block, req.Detail, req.Visible)
	utils.RespondWithSuccess(c, response)
}

// GetBlockBalance handles POST /wallet/getblockbalance
// Returns balance changes in block
func (h *BlockHandler) GetBlockBalance(c *gin.Context) {
	var req struct {
		Hash    string `json:"hash" binding:"required"`
		Number  int64  `json:"number" binding:"required"`
		Visible bool   `json:"visible" default:"false"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := h.blockchainClient.GetBlockBalance(context.Background(), &lindapb.BlockBalanceReq{
		Hash:   []byte(req.Hash),
		Number: req.Number,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get block balance: "+err.Error())
		return
	}

	response := gin.H{
		"timestamp":        resp.Timestamp,
		"block_identifier": resp.BlockIdentifier,
		"transaction_balance_trace": resp.TransactionBalanceTrace,
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Solidity Node Block Methods ====================

// GetNowBlockSolidity handles POST /walletsolidity/getnowblock
// Returns most recent confirmed block
func (h *BlockHandler) GetNowBlockSolidity(c *gin.Context) {
	block, err := h.blockchainClient.GetNowBlockSolidity(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get now block: "+err.Error())
		return
	}

	response := convertBlockToResponse(block, false)
	utils.RespondWithSuccess(c, response)
}

// GetBlockByNumSolidity handles POST /walletsolidity/getblockbynum
// Returns confirmed block by number
func (h *BlockHandler) GetBlockByNumSolidity(c *gin.Context) {
	var req struct {
		Num int64 `json:"num" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	block, err := h.blockchainClient.GetBlockByNumSolidity(context.Background(), &lindapb.NumberMessage{
		Num: req.Num,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get block: "+err.Error())
		return
	}

	response := convertBlockToResponse(block, false)
	utils.RespondWithSuccess(c, response)
}

// ==================== Event Query Block Methods ====================

// GetBlockByHashEvent handles GET /v1/blocks/{hash}
// Returns block by hash from event service
func (h *BlockHandler) GetBlockByHashEvent(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Hash is required")
		return
	}

	// Get from blockchain
	block, err := h.blockchainClient.GetBlockById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(hash),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Block not found")
		return
	}

	response := &models.EventBlockResponse{
		Hash:             string(block.BlockID),
		Number:           block.BlockHeader.RawData.Number,
		Timestamp:        block.BlockHeader.RawData.Timestamp,
		ParentHash:       string(block.BlockHeader.RawData.ParentHash),
		WitnessAddress:   string(block.BlockHeader.RawData.WitnessAddress),
		TransactionCount: len(block.Transactions),
	}

	// Convert witness address to base58
	base58Addr, _ := utils.HexToBase58(string(block.BlockHeader.RawData.WitnessAddress))
	response.WitnessAddress = base58Addr

	utils.RespondWithSuccess(c, gin.H{
		"data":    []*models.EventBlockResponse{response},
		"success": true,
		"meta": gin.H{
			"at": utils.GetCurrentTimestamp(),
		},
	})
}

// GetBlocksEvent handles GET /v1/blocks
// Returns paginated blocks from event service
func (h *BlockHandler) GetBlocksEvent(c *gin.Context) {
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

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	// Get from database
	blocks, total, err := h.blockRepo.GetBlocks(req.Block, req.Start, req.Limit, req.Sort)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get blocks: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{
		"data":    blocks,
		"success": true,
		"meta": gin.H{
			"at":        utils.GetCurrentTimestamp(),
			"page_size": req.Limit,
		},
	})
}

// GetLatestSolidifiedBlockNumber handles GET /v1/blocks/latestSolidifiedBlockNumber
// Returns latest solidified block number
func (h *BlockHandler) GetLatestSolidifiedBlockNumber(c *gin.Context) {
	// Get from solidity node
	block, err := h.blockchainClient.GetNowBlockSolidity(context.Background(), &lindapb.EmptyMessage{})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get latest block: "+err.Error())
		return
	}

	utils.RespondWithSuccess(c, gin.H{
		"number": block.BlockHeader.RawData.Number,
	})
}

// GetBlockStats handles GET /v1/blocks/{blockNum}/stats
// Returns block statistics
func (h *BlockHandler) GetBlockStats(c *gin.Context) {
	blockNumStr := c.Param("blockNum")
	blockNum, err := strconv.ParseInt(blockNumStr, 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid block number")
		return
	}

	getTxDetail := c.DefaultQuery("get_tx_detail", "false") == "true"

	// Get block info
	block, err := h.blockchainClient.GetBlockByNum(context.Background(), &lindapb.NumberMessage{
		Num: blockNum,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Block not found")
		return
	}

	// Calculate statistics
	stats := calculateBlockStats(block, getTxDetail)

	utils.RespondWithSuccess(c, stats)
}

// ==================== Lindascan Custom Block Methods ====================

// GetBlocksV2 handles GET /api/block
// Returns paginated blocks for lindascan
func (h *BlockHandler) GetBlocksV2(c *gin.Context) {
	var req struct {
		Number int64  `form:"number"`
		Hash   string `form:"hash"`
		Limit  int    `form:"limit" default:"20"`
		Start  int    `form:"start" default:"0"`
		Sort   string `form:"sort" default:"-number"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	var blocks []*models.BlockResponse
	var total int64

	if req.Number > 0 {
		// Get by number
		block, err := h.blockchainClient.GetBlockByNum(context.Background(), &lindapb.NumberMessage{
			Num: req.Number,
		})
		if err == nil {
			blocks = append(blocks, convertBlockToResponse(block, true))
			total = 1
		}
	} else if req.Hash != "" {
		// Get by hash
		block, err := h.blockchainClient.GetBlockById(context.Background(), &lindapb.BytesMessage{
			Value: []byte(req.Hash),
		})
		if err == nil {
			blocks = append(blocks, convertBlockToResponse(block, true))
			total = 1
		}
	} else {
		// Get paginated
		// This would come from database
		blocks, total, _ = h.blockRepo.GetBlocks(0, req.Start, req.Limit, req.Sort)
	}

	response := gin.H{
		"blocks": blocks,
		"total":  total,
	}

	utils.RespondWithSuccess(c, response)
}

// ==================== Helper Functions ====================

func convertBlockToResponse(block *lindapb.Block, visible bool) *models.BlockResponse {
	resp := &models.BlockResponse{
		BlockID: string(block.BlockID),
	}

	if block.BlockHeader != nil {
		resp.BlockHeader = &models.BlockHeader{
			WitnessSignature: string(block.BlockHeader.WitnessSignature),
		}

		if block.BlockHeader.RawData != nil {
			raw := block.BlockHeader.RawData
			resp.BlockHeader.RawData = &models.BlockRawData{
				Timestamp:        raw.Timestamp,
				TxTrieRoot:       string(raw.TxTrieRoot),
				ParentHash:       string(raw.ParentHash),
				Number:           raw.Number,
				WitnessID:        raw.WitnessId,
				WitnessAddress:   string(raw.WitnessAddress),
				Version:          raw.Version,
				AccountStateRoot: string(raw.AccountStateRoot),
			}

			if visible {
				base58Addr, _ := utils.HexToBase58(string(raw.WitnessAddress))
				resp.BlockHeader.RawData.WitnessAddress = base58Addr
			}
		}
	}

	// Convert transactions
	for _, tx := range block.Transactions {
		resp.Transactions = append(resp.Transactions, *convertTransactionToResponse(tx, visible))
	}

	return resp
}

func convertBlockExtentionToResponse(block *lindapb.BlockExtention, detail bool, visible bool) gin.H {
	resp := gin.H{
		"blockID": string(block.BlockID),
	}

	if block.BlockHeader != nil {
		header := gin.H{
			"witness_signature": string(block.BlockHeader.WitnessSignature),
		}

		if block.BlockHeader.RawData != nil {
			raw := block.BlockHeader.RawData
			rawData := gin.H{
				"timestamp":         raw.Timestamp,
				"txTrieRoot":        string(raw.TxTrieRoot),
				"parentHash":        string(raw.ParentHash),
				"number":            raw.Number,
				"witness_id":        raw.WitnessId,
				"witness_address":   string(raw.WitnessAddress),
				"version":           raw.Version,
				"accountStateRoot":  string(raw.AccountStateRoot),
			}

			if visible {
				base58Addr, _ := utils.HexToBase58(string(raw.WitnessAddress))
				rawData["witness_address"] = base58Addr
			}

			header["raw_data"] = rawData
		}

		resp["block_header"] = header
	}

	if detail {
		var txs []gin.H
		for _, tx := range block.Transactions {
			txs = append(txs, convertTransactionExtentionToResponse(tx, visible))
		}
		resp["transactions"] = txs
	}

	return resp
}

func convertTransactionExtentionToResponse(tx *lindapb.Transaction, visible bool) gin.H {
	resp := gin.H{
		"txID":        string(tx.TxID),
		"raw_data_hex": string(tx.RawDataHex),
		"signature":   tx.Signature,
	}

	var rets []gin.H
	for _, r := range tx.Ret {
		rets = append(rets, gin.H{
			"contractRet": r.Ret.String(),
			"fee":         r.Fee,
		})
	}
	resp["ret"] = rets

	if tx.RawData != nil {
		rawData := gin.H{
			"ref_block_bytes": string(tx.RawData.RefBlockBytes),
			"ref_block_num":   tx.RawData.RefBlockNum,
			"ref_block_hash":  string(tx.RawData.RefBlockHash),
			"expiration":      tx.RawData.Expiration,
			"timestamp":       tx.RawData.Timestamp,
		}

		var contracts []gin.H
		for _, c := range tx.RawData.Contract {
			contract := gin.H{
				"type":      c.Type.String(),
				"parameter": c.Parameter.Value,
			}
			contracts = append(contracts, contract)
		}
		rawData["contract"] = contracts

		resp["raw_data"] = rawData
	}

	return resp
}

func calculateBlockStats(block *lindapb.Block, getTxDetail bool) *models.BlockStatsResponse {
	stats := &models.BlockStatsResponse{
		FeeStat: &models.FeeStat{},
	}

	if getTxDetail {
		stats.TxStat = &models.TxStat{
			ContractTypeDistribute: make(map[int]int),
		}
	}

	// Calculate statistics from transactions
	for _, tx := range block.Transactions {
		// Count transaction types
		if getTxDetail && len(tx.RawData.Contract) > 0 {
			contractType := int(tx.RawData.Contract[0].Type)
			stats.TxStat.ContractTypeDistribute[contractType]++
		}

		// Calculate fees
		for _, ret := range tx.Ret {
			stats.FeeStat.OtherFee += ret.Fee
		}

		// Check if failed
		if getTxDetail && len(tx.Ret) > 0 && tx.Ret[0].Ret == lindapb.Transaction_Result_FAILED {
			stats.TxStat.FailTxCount++
		}
	}

	return stats
}