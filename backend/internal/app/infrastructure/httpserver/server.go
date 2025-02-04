package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Mixturka/DockerLens/backend/internal/app/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	router     chi.Router
}

func NewServer(r chi.Router) *Server {
	s := &Server{router: r}
	return s
}

func (s *Server) Start(log *slog.Logger, config config.Config) error {
	log.Info("Http server listening", "address", config.ListenAddr)
	s.httpServer = &http.Server{Addr: config.ListenAddr, Handler: s.router}
	err := s.httpServer.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Error("Http server stopped unexpected")
	}

	return nil
}

func (s *Server) Shutdown(log *slog.Logger) {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to shutdown http server correctly %s", err.Error()))
		} else {
			s.httpServer = nil
		}
	}
}
