package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CtxKey string

const (
	CtxKeyReqID CtxKey = "request_id"
)

func (m *MiddlewareHandler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), CtxKeyReqID, reqID)
		start := time.Now()
		slog.Info("request received",
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"time_received", start,
			"client_ip", r.RemoteAddr,
		)

		rw := newLogRW(w)
		next.ServeHTTP(rw, r.WithContext(ctx))

		duration := time.Since(start)
		slog.Info("request completed",
			"request_id", reqID,
			"status", rw.Status(),
			"duration_ms", duration.Milliseconds(),
			"size", rw.Size(),
		)
	})
}
