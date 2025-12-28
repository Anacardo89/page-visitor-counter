package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Anacardo89/page-visitor-counter/internal/middleware"
	"github.com/Anacardo89/page-visitor-counter/internal/repo"
	"github.com/Anacardo89/page-visitor-counter/pkg/validator"
)

// Handler
type visitorAPIHandler struct {
	repo repo.VisitorRepo
}

func NewVisitorAPI(repo repo.VisitorRepo) VisitorAPI {
	return &visitorAPIHandler{
		repo: repo,
	}
}

// POST /visitors
type AddVisitorReq struct {
	Url       string `json:"url" validate:"required"`
	VisitorID string `json:"visitor_id" validate:"required"`
}

func (h *visitorAPIHandler) AddVisitor(w http.ResponseWriter, r *http.Request) {
	// Setup
	reqID := r.Context().Value(middleware.CtxKeyReqID)

	// Execution
	// Input validation
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", "request_id", reqID, "error", err)
		failHTTP(w, reqID, http.StatusBadRequest, "invalid request body")
		return
	}
	var reqBody AddVisitorReq
	if err := validator.ParseAndValidate(raw, &reqBody); err != nil {
		var ve validator.ValidationError
		if errors.As(err, &ve) {
			slog.Error("missing required fields", "request_id", reqID, "error", err)
			failHTTP(w, reqID, http.StatusBadRequest, err.Error())
		} else {
			slog.Error("failed to parse JSON from body", "request_id", reqID, "error", err)
			failHTTP(w, reqID, http.StatusBadRequest, "invalid request body")
		}
		return
	}
	// Repo operation
	h.repo.AddVisitor(reqBody.Url, reqBody.VisitorID)
	// Response
	w.WriteHeader(http.StatusNoContent)
}

// GET /visitors
type CountVisitorsResp struct {
	VisitorCount int `json:"visitor_count"`
}

func (h *visitorAPIHandler) CountVisitors(w http.ResponseWriter, r *http.Request) {
	// Setup
	reqID := r.Context().Value(middleware.CtxKeyReqID)
	w.Header().Set("Content-Type", "application/json")

	// Execution
	// Input validation
	url := r.URL.Query().Get("url")
	if url == "" {
		slog.Error("no url query param provided", "request_id", reqID)
		failHTTP(w, reqID, http.StatusBadRequest, "invalid request")
		return
	}
	// Repo operation
	visitorCount := h.repo.CountVisitors(url)
	// Response
	respBody := CountVisitorsResp{
		VisitorCount: visitorCount,
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(respBody); err != nil {
		slog.Error("failed to encode response body", "request_id", reqID, "error", err)
		failHTTP(w, reqID, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf.Bytes()); err != nil {
		slog.Error("failed to send response to client", "request_id", reqID, "error", err)
	}
}
