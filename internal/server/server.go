package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sibhellyx/imageProccesor/api"
	"github.com/sibhellyx/imageProccesor/config"
	"github.com/sibhellyx/imageProccesor/internal/handlers"
	"github.com/sibhellyx/imageProccesor/internal/repository"
	"github.com/sibhellyx/imageProccesor/internal/service"
	"github.com/sibhellyx/imageProccesor/internal/workerpool/pool"
	"github.com/sibhellyx/imageProccesor/pkg/image"
)

type Server struct {
	ctx     context.Context
	cfg     config.Config
	srv     *http.Server
	service *service.Service
	pool    *pool.Pool
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
	repo := repository.NewRepository()
	s.pool = pool.NewPool(image.ProccesImage, s.cfg.Workers)
	s.service = service.NewService(repo, s.pool, s.cfg.QueueCapacity)
	handler := handlers.NewHandler(s.service)
	routes := api.CreateRoutes(handler)

	s.srv = &http.Server{
		Addr:    ":" + s.cfg.Port,
		Handler: routes,
	}

	log.Printf("Starting server on :%s", s.cfg.Port)
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}

func (s *Server) Shutdown() {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.pool.Shutdown()

	err := s.srv.Shutdown(ctxShutdown)
	if err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

}
