package dtos

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

func WriteErrorResponse(w http.ResponseWriter, message, error string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := ErrorResponse{
		Message: message,
		Error:   error,
		Code:    statusCode,
	}

	json.NewEncoder(w).Encode(errorResp)
}
