package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Anacardo89/page-visitor-counter/config"
	"github.com/Anacardo89/page-visitor-counter/internal/api"
	"github.com/Anacardo89/page-visitor-counter/internal/middleware"
)

type Server struct {
	httpSrv  *http.Server
	router   http.Handler
	addr     string
	timeouts ServerTimeouts
}

type ServerTimeouts struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func NewServer(cfg *config.Server, ah api.VisitorAPI, mw *middleware.MiddlewareHandler) *Server {
	to := ServerTimeouts{
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		ShutdownTimeout: cfg.ShutdownTimeout,
	}
	s := &Server{
		router:   NewRouter(ah, mw),
		addr:     fmt.Sprintf(":%s", cfg.Port),
		timeouts: to,
	}
	return s
}

func (s *Server) Start() error {
	s.httpSrv = &http.Server{
		Addr:         s.addr,
		Handler:      s.router,
		ReadTimeout:  s.timeouts.ReadTimeout,
		WriteTimeout: s.timeouts.WriteTimeout,
	}
	slog.Info("Starting server on", "address", s.addr)
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeouts.ShutdownTimeout)
	defer cancel()
	if s.httpSrv != nil {
		return s.httpSrv.Shutdown(ctx)
	}
	return nil
}
