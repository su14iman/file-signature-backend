package models

import (
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type DocumentInfo struct {
	FilePath      string
	FileType      string
	FileSignature string
	CreatedAt     time.Time
}
