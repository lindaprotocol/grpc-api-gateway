package utils

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Response standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// V1Response structure for LindaGrid v1 compatible responses
type V1Response struct {
	Data    []interface{} `json:"data"`
	Success bool          `json:"success"`
	Meta    *V1Meta       `json:"meta,omitempty"`
}

type V1Meta struct {
	At          int64  `json:"at"`
	PageSize    int    `json:"page_size"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Links       *Links `json:"links,omitempty"`
}

type Links struct {
	Next string `json:"next,omitempty"`
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// RespondWithError sends an error response (Gin)
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// RespondWithErrorHTTP sends an error response for standard net/http handlers
func RespondWithErrorHTTP(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}

// RespondWithErrorCode sends an error response with specific error code
func RespondWithErrorCode(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

// RespondWithMeta sends a response with metadata
func RespondWithMeta(c *gin.Context, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// RespondWithPagination sends a paginated response
func RespondWithPagination(c *gin.Context, data interface{}, total int64, page, limit int) {
	totalPages := (total + int64(limit) - 1) / int64(limit)
	
	RespondWithMeta(c, data, gin.H{
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"has_next":    page < int(totalPages),
		"has_prev":    page > 1,
	})
}

// RespondWithV1Success sends a LindaGrid v1 compatible success response
func RespondWithV1Success(c *gin.Context, data []interface{}, pageSize int, fingerprint string) {
	response := V1Response{
		Data:    data,
		Success: true,
		Meta: &V1Meta{
			At:          GetCurrentTimestamp(),
			PageSize:    pageSize,
			Fingerprint: fingerprint,
		},
	}

	// Add next link if fingerprint exists
	if fingerprint != "" {
		response.Meta.Links = &Links{
			Next: buildNextPageLink(c, fingerprint, pageSize),
		}
	}

	c.JSON(http.StatusOK, response)
}

// RespondWithV1Error sends a LindaGrid v1 compatible error response
func RespondWithV1Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// RespondWithCSV sends a CSV response
func RespondWithCSV(c *gin.Context, filename string, headers []string, rows [][]string) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename="+filename)
	c.Header("Cache-Control", "no-cache")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(headers); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to generate CSV")
		return
	}

	// Write rows
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			RespondWithError(c, http.StatusInternalServerError, "Failed to generate CSV")
			return
		}
	}
}

// RespondWithJSON sends a JSON response with proper content type
func RespondWithJSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// RespondWithNoContent sends a 204 No Content response
func RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// RespondWithCreated sends a 201 Created response
func RespondWithCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// RespondWithAccepted sends a 202 Accepted response
func RespondWithAccepted(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAccepted, Response{
		Success: true,
		Data:    data,
	})
}

// RespondWithBadRequest sends a 400 Bad Request response
func RespondWithBadRequest(c *gin.Context, message string) {
	RespondWithError(c, http.StatusBadRequest, message)
}

// RespondWithUnauthorized sends a 401 Unauthorized response
func RespondWithUnauthorized(c *gin.Context, message string) {
	RespondWithError(c, http.StatusUnauthorized, message)
}

// RespondWithForbidden sends a 403 Forbidden response
func RespondWithForbidden(c *gin.Context, message string) {
	RespondWithError(c, http.StatusForbidden, message)
}

// RespondWithNotFound sends a 404 Not Found response
func RespondWithNotFound(c *gin.Context, message string) {
	RespondWithError(c, http.StatusNotFound, message)
}

// RespondWithConflict sends a 409 Conflict response
func RespondWithConflict(c *gin.Context, message string) {
	RespondWithError(c, http.StatusConflict, message)
}

// RespondWithTooManyRequests sends a 429 Too Many Requests response
func RespondWithTooManyRequests(c *gin.Context, message string) {
	RespondWithError(c, http.StatusTooManyRequests, message)
}

// RespondWithInternalError sends a 500 Internal Server Error response
func RespondWithInternalError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusInternalServerError, message)
}

// RespondWithServiceUnavailable sends a 503 Service Unavailable response
func RespondWithServiceUnavailable(c *gin.Context, message string) {
	RespondWithError(c, http.StatusServiceUnavailable, message)
}

// buildNextPageLink builds the next page link for v1 responses
func buildNextPageLink(c *gin.Context, fingerprint string, pageSize int) string {
	// Build URL with fingerprint for next page
	query := c.Request.URL.Query()
	query.Set("fingerprint", fingerprint)
	query.Set("limit", strconv.Itoa(pageSize))
	
	return c.Request.URL.Path + "?" + query.Encode()
}

// GetCurrentTimestamp returns current timestamp in milliseconds
func GetCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// ParsePaginationParams parses pagination parameters from request
// func ParsePaginationParams(c *gin.Context) (limit, offset int, sort string, fingerprint string) {
// 	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))
// 	if limit <= 0 || limit > 200 {
// 		limit = 20
// 	}

// 	offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
// 	if offset < 0 {
// 		offset = 0
// 	}

// 	sort = c.DefaultQuery("sort", "-timestamp")
// 	fingerprint = c.Query("fingerprint")

// 	return limit, offset, sort, fingerprint
// }

// ParseV1PaginationParams parses v1 style pagination parameters
func ParseV1PaginationParams(c *gin.Context) (limit, start int, sort string, fingerprint string) {
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 || limit > 200 {
		limit = 20
	}

	start, _ = strconv.Atoi(c.DefaultQuery("start", "0"))
	if start < 0 {
		start = 0
	}

	sort = c.DefaultQuery("sort", "-timestamp")
	if sort == "timestamp" {
		sort = "timestamp"
	} else if sort == "-timestamp" {
		sort = "-timestamp"
	}

	fingerprint = c.Query("fingerprint")

	return limit, start, sort, fingerprint
}