package api

import (
	"encoding/json"
	"net/http"
)

// Healthcheck - /
type HealthCheckResp struct {
	Status string `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := HealthCheckResp{
		Status: "OK",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}
