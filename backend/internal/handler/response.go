package handler

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
    HttpCode int    `json:"httpCode"`
    Message  string `json:"message"`
    Data     any    `json:"data"`
    Error    any    `json:"error,omitempty"`

    Page  *int   `json:"page,omitempty"`
    Limit *int   `json:"limit,omitempty"`
    Total *int64 `json:"total,omitempty"`
}

type APIError struct {
    Code string `json:"code"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    _ = json.NewEncoder(w).Encode(data)
}

func writeResponse(
    w http.ResponseWriter,
    status int,
    message string,
    data any,
) {
    writeJSON(w, status, APIResponse{
        HttpCode: status,
        Message:  message,
        Data:     data,
    })
}

func writeListResponse(
    w http.ResponseWriter,
    message string,
    data any,
    page int,
    limit int,
    total int64,
) {
    writeJSON(w, http.StatusOK, APIResponse{
        HttpCode: http.StatusOK,
        Message:  message,
        Data:     data,
        Page:     &page,
        Limit:    &limit,
        Total:    &total,
    })
}

func writeError(
    w http.ResponseWriter,
    status int,
    message string,
    errorCode string,
) {
    writeJSON(w, status, APIResponse{
        HttpCode: status,
        Message:  message,
        Error: APIError{
            Code: errorCode,
        },
    })
}