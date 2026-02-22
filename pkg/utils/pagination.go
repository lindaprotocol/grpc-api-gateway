package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Limit       int
	Offset      int
	Page        int
	Sort        string
	Fingerprint string
	Total       int64 // Add this field
}

// ParsePaginationParams parses pagination parameters from request
func ParsePaginationParams(c *gin.Context) *PaginationParams {
	params := &PaginationParams{
		Limit:       20,
		Offset:      0,
		Page:        1,
		Sort:        "-timestamp",
		Fingerprint: "",
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > 200 {
				params.Limit = 200
			} else {
				params.Limit = limit
			}
		}
	}

	// Parse offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			params.Offset = offset
			params.Page = offset/params.Limit + 1
		}
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
			params.Offset = (page - 1) * params.Limit
		}
	}

	// Parse start (v1 style)
	if startStr := c.Query("start"); startStr != "" {
		if start, err := strconv.Atoi(startStr); err == nil && start >= 0 {
			params.Offset = start
			params.Page = start/params.Limit + 1
		}
	}

	// Parse sort
	if sortStr := c.Query("sort"); sortStr != "" {
		params.Sort = sortStr
	}

	// Parse fingerprint
	if fpStr := c.Query("fingerprint"); fpStr != "" {
		params.Fingerprint = fpStr
	}

	return params
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total       int64  `json:"total"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	TotalPages  int    `json:"total_pages"`
	HasNext     bool   `json:"has_next"`
	HasPrev     bool   `json:"has_prev"`
	NextCursor  string `json:"next_cursor,omitempty"`
	PrevCursor  string `json:"prev_cursor,omitempty"`
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(total int64, page, limit int) *PaginationMeta {
	totalPages := (int(total) + limit - 1) / limit

	return &PaginationMeta{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// CursorPagination represents cursor-based pagination
type CursorPagination struct {
	Limit       int
	After       string
	Before      string
	Sort        string
}

// ParseCursorPagination parses cursor pagination parameters
func ParseCursorPagination(c *gin.Context) *CursorPagination {
	params := &CursorPagination{
		Limit:  20,
		After:  c.Query("after"),
		Before: c.Query("before"),
		Sort:   c.DefaultQuery("sort", "-timestamp"),
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > 200 {
				params.Limit = 200
			} else {
				params.Limit = limit
			}
		}
	}

	return params
}

// BuildNextPageLink builds the next page link
func BuildNextPageLink(baseURL string, params *PaginationParams) string {
	if params.Offset+params.Limit >= int(params.Total) {
		return ""
	}

	nextOffset := params.Offset + params.Limit
	return baseURL + "?offset=" + strconv.Itoa(nextOffset) + "&limit=" + strconv.Itoa(params.Limit)
}

// BuildPrevPageLink builds the previous page link
func BuildPrevPageLink(baseURL string, params *PaginationParams) string {
	if params.Offset == 0 {
		return ""
	}

	prevOffset := params.Offset - params.Limit
	if prevOffset < 0 {
		prevOffset = 0
	}
	return baseURL + "?offset=" + strconv.Itoa(prevOffset) + "&limit=" + strconv.Itoa(params.Limit)
}

// BuildCursorNextLink builds the next cursor link
func BuildCursorNextLink(baseURL string, cursor string, limit int) string {
	if cursor == "" {
		return ""
	}
	return baseURL + "?after=" + cursor + "&limit=" + strconv.Itoa(limit)
}

// BuildCursorPrevLink builds the previous cursor link
func BuildCursorPrevLink(baseURL string, cursor string, limit int) string {
	if cursor == "" {
		return ""
	}
	return baseURL + "?before=" + cursor + "&limit=" + strconv.Itoa(limit)
}

// CalculateOffset calculates offset from page and limit
func CalculateOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}

// CalculatePage calculates page from offset and limit
func CalculatePage(offset, limit int) int {
	if limit <= 0 {
		return 1
	}
	return offset/limit + 1
}

// ValidatePagination validates pagination parameters
func ValidatePagination(limit, maxLimit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

// SortDirection returns the sort direction from a sort string
func SortDirection(sort string) (field string, desc bool) {
	if len(sort) == 0 {
		return "timestamp", true
	}
	if sort[0] == '-' {
		return sort[1:], true
	}
	return sort, false
}