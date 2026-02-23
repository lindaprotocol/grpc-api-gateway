package handlers

import (
	"context"
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

type SearchHandler struct {
	blockchainClient *blockchain.Client
	accountRepo      *repository.AccountRepository
	blockRepo        *repository.BlockRepository
	txRepo           *repository.TransactionRepository
	tokenRepo        *repository.TokenRepository
}

func NewSearchHandler(
	client *blockchain.Client,
	accountRepo *repository.AccountRepository,
	blockRepo *repository.BlockRepository,
	txRepo *repository.TransactionRepository,
	tokenRepo *repository.TokenRepository,
) *SearchHandler {
	return &SearchHandler{
		blockchainClient: client,
		accountRepo:      accountRepo,
		blockRepo:        blockRepo,
		txRepo:           txRepo,
		tokenRepo:        tokenRepo,
	}
}

// Search handles GET /api/search
func (h *SearchHandler) Search(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Query == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Query is required")
		return
	}

	results := make([]models.SearchResult, 0)

	// Determine search type if not specified
	searchType := req.Type
	if searchType == "" {
		searchType = h.determineSearchType(req.Query)
	}

	switch searchType {
	case "block":
		result := h.searchBlock(req.Query)
		if result != nil {
			results = append(results, *result)
		}
	case "transaction":
		result := h.searchTransaction(req.Query)
		if result != nil {
			results = append(results, *result)
		}
	case "address":
		result := h.searchAddress(req.Query)
		if result != nil {
			results = append(results, *result)
		}
	case "token":
		tokenResults := h.searchToken(req.Query)
		results = append(results, tokenResults...)
	default:
		// Search all types
		if blockResult := h.searchBlock(req.Query); blockResult != nil {
			results = append(results, *blockResult)
		}
		if txResult := h.searchTransaction(req.Query); txResult != nil {
			results = append(results, *txResult)
		}
		if addrResult := h.searchAddress(req.Query); addrResult != nil {
			results = append(results, *addrResult)
		}
		tokenResults := h.searchToken(req.Query)
		results = append(results, tokenResults...)
	}

	response := &models.SearchResponse{
		Results: results,
	}

	utils.RespondWithSuccess(c, response)
}

// searchBlock searches for a block by number or hash
func (h *SearchHandler) searchBlock(query string) *models.SearchResult {
	// Try as number first
	if num, err := strconv.ParseInt(query, 10, 64); err == nil {
		block, err := h.blockchainClient.GetBlockByNum(context.Background(), &lindapb.NumberMessage{
			Num: num,
		})
		if err == nil && block != nil {
			return &models.SearchResult{
				Type:        "block",
				ID:          strconv.FormatInt(num, 10),
				Name:        "Block #" + strconv.FormatInt(num, 10),
				URL:         "/#/block/" + strconv.FormatInt(num, 10),
				Description: "Block at height " + strconv.FormatInt(num, 10),
			}
		}
	}

	// Try as hash
	if len(query) == 64 {
		block, err := h.blockchainClient.GetBlockById(context.Background(), &lindapb.BytesMessage{
			Value: []byte(query),
		})
		if err == nil && block != nil {
			num := block.BlockHeader.RawData.Number
			return &models.SearchResult{
				Type:        "block",
				ID:          strconv.FormatInt(num, 10),
				Name:        "Block #" + strconv.FormatInt(num, 10),
				URL:         "/#/block/" + strconv.FormatInt(num, 10),
				Description: "Block with hash " + query[:16] + "...",
			}
		}
	}

	return nil
}

// searchTransaction searches for a transaction by hash
func (h *SearchHandler) searchTransaction(query string) *models.SearchResult {
	if len(query) != 64 {
		return nil
	}

	tx, err := h.blockchainClient.GetTransactionById(context.Background(), &lindapb.BytesMessage{
		Value: []byte(query),
	})
	if err == nil && tx != nil {
		return &models.SearchResult{
			Type:        "transaction",
			ID:          query,
			Name:        "Transaction " + query[:8] + "...",
			URL:         "/#/transaction/" + query,
			Description: "Transaction with hash " + query[:16] + "...",
		}
	}

	return nil
}

// searchAddress searches for an address (account or contract)
func (h *SearchHandler) searchAddress(query string) *models.SearchResult {
	// Validate address format
	if !utils.IsValidBase58Address(query) && !utils.IsValidHexAddress(query) {
		return nil
	}

	// Convert to hex for blockchain query
	hexAddr := query
	if !utils.IsValidHexAddress(query) {
		var err error
		hexAddr, err = utils.Base58ToHex(query)
		if err != nil {
			return nil
		}
	}

	// Try as account
	account, err := h.blockchainClient.GetAccount(context.Background(), &lindapb.Account{
		Address: []byte(hexAddr),
	})
	if err == nil && account != nil {
		accountType := "address"
		if account.ContractAddress != nil {
			accountType = "contract"
		}
		return &models.SearchResult{
			Type:        accountType,
			ID:          query,
			Name:        utils.TruncateString(query, 12),
			URL:         "/#/address/" + query,
			Description: accountType + " with address " + query[:8] + "...",
		}
	}

	return nil
}

// searchToken searches for tokens by name or symbol
func (h *SearchHandler) searchToken(query string) []models.SearchResult {
	results := make([]models.SearchResult, 0)

	// Search in database
	tokens, err := h.tokenRepo.SearchTokens(query, 10)
	if err != nil {
		return results
	}

	for _, token := range tokens {
		url := "/#/token/" + token.Contract
		if token.Symbol == "LIND" {
			url = "/#/token/" + token.Contract
		}
		results = append(results, models.SearchResult{
			Type:        "token",
			ID:          token.Contract,
			Name:        token.Symbol + " (" + token.Name + ")",
			URL:         url,
			Description: token.Name + " token with symbol " + token.Symbol,
		})
	}

	return results
}

// determineSearchType determines the type of search based on query format
func (h *SearchHandler) determineSearchType(query string) string {
	// Check if it's a number (block)
	if _, err := strconv.ParseInt(query, 10, 64); err == nil {
		return "block"
	}

	// Check if it's a transaction hash (64 hex chars)
	if len(query) == 64 && h.isHexString(query) {
		return "transaction"
	}

	// Check if it's an address (base58 or hex with prefix)
	if utils.IsValidBase58Address(query) || utils.IsValidHexAddress(query) {
		return "address"
	}

	// Default to token search
	return "token"
}

// isHexString checks if a string is a valid hex string
func (h *SearchHandler) isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}