// internal/models/types.go
package models

import "encoding/json"

// JSON is a type alias for json.RawMessage
type JSON json.RawMessage

// PaginationRequest represents common pagination parameters
type PaginationRequest struct {
	Limit       int    `json:"limit" form:"limit" default:"20"`
	Start       int    `json:"start" form:"start" default:"0"`
	Fingerprint string `json:"fingerprint" form:"fingerprint"`
	Sort        string `json:"sort" form:"sort" default:"-timestamp"`
}