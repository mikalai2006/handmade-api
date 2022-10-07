package server

import (
	"context"
	"net/http"

	"github.com/mikalai2006/handmade/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, hadler http.Handler) *Server  {
	return &Server{
		httpServer: &http.Server{
			Addr: ":" + cfg.HTTP.Port,
			Handler: hadler,
			ReadTimeout: cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}