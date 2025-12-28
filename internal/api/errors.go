package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResp struct {
	Error string `json:"error"`
}

// 404 Handler
func CatchAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := ErrorResp{
		Error: "endpoint not found",
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(body)
}

// Error helper
func failHTTP(w http.ResponseWriter, reqID any, status int, outMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResp{Error: outMsg}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode error response body", "request_id", reqID, "error", err)
	}
}
