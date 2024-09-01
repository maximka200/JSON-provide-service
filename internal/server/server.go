package server

import (
	"context"
	"fmt"
	"jps/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg config.Config, handler http.Handler) error {
	op := "server.Run"
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 mb
		ReadTimeout:    cfg.Timeout,
		WriteTimeout:   cfg.Timeout,
	}
	return fmt.Errorf("%w:%s", s.httpServer.ListenAndServe(), op)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
