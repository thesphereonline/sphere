package server

import (
	"context"
	"fmt"
	"net/http"

	"log"

	"github.com/go-chi/chi/v5"
	"github.com/thesphereonline/sphere/config"
	"github.com/thesphereonline/sphere/internal/infra/db"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, pg *db.Postgres, logger *log.Logger) *Server {
	r := chi.NewRouter()

	// Sample route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: r,
	}

	return &Server{httpServer: s}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
