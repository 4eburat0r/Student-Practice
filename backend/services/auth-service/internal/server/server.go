package server

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	handler    *handler.Handler
}

func NewServer(cfg *config.Config, h *handler.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		handler: h,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
