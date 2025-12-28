package server

import (
	"net/http"

	"github.com/Anacardo89/page-visitor-counter/internal/api"
	"github.com/Anacardo89/page-visitor-counter/internal/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(ah api.VisitorAPI, mw *middleware.MiddlewareHandler) http.Handler {
	r := mux.NewRouter()
	// Health check
	r.Handle("/", http.HandlerFunc(api.HealthCheck)).Methods("GET")
	// Visitors
	r.Handle("/visitors", http.HandlerFunc(ah.AddVisitor)).Methods("POST")
	r.Handle("/visitors", http.HandlerFunc(ah.CountVisitors)).Methods("GET")
	// Catch-all 404
	r.NotFoundHandler = http.HandlerFunc(api.CatchAll)
	return mw.Wrap(r)
}
