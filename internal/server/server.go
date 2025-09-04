package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sibhellyx/imageProccesor/config"
)

type Server struct {
	ctx context.Context
	cfg config.Config
	srv *http.Server
}

func NewServer(ctx context.Context, cfg config.Config) *Server {
	server := new(Server)
	server.ctx = ctx
	server.cfg = cfg
	return server
}

func (s *Server) Serve() {
	s.srv = &http.Server{
		Addr: ":" + s.cfg.Port,
	}

	log.Printf("Starting server on :%s", s.cfg.Port)
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}

func (s *Server) Shutdown() {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.srv.Shutdown(ctxShutdown)
	if err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
