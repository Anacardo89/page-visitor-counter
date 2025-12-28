package middleware

import (
	"net/http"
	"time"
)

type MiddlewareHandler struct {
	writeTimeout time.Duration
}

func NewMiddlewareHandler(to time.Duration) *MiddlewareHandler {
	return &MiddlewareHandler{
		writeTimeout: to - time.Second,
	}
}

func (m *MiddlewareHandler) Wrap(next http.Handler) http.Handler {
	return m.Log(m.Timeout(next))
}
