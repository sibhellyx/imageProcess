package service

import (
	"github.com/sibhellyx/imageProccesor/internal/repository"
	"github.com/sibhellyx/imageProccesor/internal/workerpool/pool"
)

type Service struct {
	repository *repository.Repository
	pool       *pool.Pool
}

func NewService(repo *repository.Repository, pool *pool.Pool) *Service {
	return &Service{
		repository: repo,
		pool:       pool,
	}
}
