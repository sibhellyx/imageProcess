package service

import (
	"context"
	"fmt"

	"github.com/sibhellyx/imageProccesor/internal/errors"
	"github.com/sibhellyx/imageProccesor/internal/repository"
	"github.com/sibhellyx/imageProccesor/internal/workerpool/pool"
)

type Service struct {
	repository *repository.Repository
	pool       *pool.Pool
	taskQueue  chan string
	limiter    chan struct{}
}

func NewService(repo *repository.Repository, pool *pool.Pool, queueCapacity int) *Service {
	s := &Service{
		repository: repo,
		pool:       pool,
		taskQueue:  make(chan string, queueCapacity),
		limiter:    make(chan struct{}, queueCapacity),
	}
	s.pool.Create()

	go s.proccess()

	return s
}

func (s *Service) Create(ctx context.Context, path string) error {
	select {
	case s.taskQueue <- path:
		//запись в репо
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return errors.ErrQueueTasksFull
	}
}

func (s *Service) proccess() {
	for imagePath := range s.taskQueue {
		s.limiter <- struct{}{}
		go func(path string) {
			defer func() { <-s.limiter }()
			result := s.pool.Handle(path)
			fmt.Println("Procces result ", <-result)
		}(imagePath)

	}
}
