package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/sibhellyx/imageProccesor/internal/errors"
	"github.com/sibhellyx/imageProccesor/internal/repository"
	"github.com/sibhellyx/imageProccesor/internal/workerpool/pool"
)

type Service struct {
	repository *repository.Repository
	pool       *pool.Pool
	taskQueue  chan string
	mu         *sync.Mutex
}

func NewService(repo *repository.Repository, pool *pool.Pool, queueCapacity int) *Service {
	s := &Service{
		repository: repo,
		pool:       pool,
		taskQueue:  make(chan string, queueCapacity),
		mu:         &sync.Mutex{},
	}
	s.pool.Create()

	go s.procces()

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

func (s *Service) StartProcces() string {
	for {
		select {
		case imagePath, ok := <-s.taskQueue:
			if !ok {
				return ""
			}
			return s.pool.Handle(imagePath)
		default:
			return ""
		}
	}
}

func (s *Service) procces() {
	for imagePath := range s.taskQueue {
		go func(path string) {
			s := s.pool.Handle(path)
			fmt.Println("Procces result ", s)
		}(imagePath)
	}
}
