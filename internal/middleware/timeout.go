package middleware

import (
	"net/http"
)

func (m *MiddlewareHandler) Timeout(next http.Handler) http.Handler {
	return http.TimeoutHandler(
		next,
		m.writeTimeout,
		"request timed out",
	)
}
