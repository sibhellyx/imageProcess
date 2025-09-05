package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sibhellyx/imageProccesor/internal/errors"
	"github.com/sibhellyx/imageProccesor/internal/repository"
	"github.com/sibhellyx/imageProccesor/internal/workerpool/pool"
)

type Service struct {
	repository *repository.Repository
	pool       *pool.Pool
	taskQueue  chan string
	limiter    chan struct{}

	ctx      context.Context
	cancel   context.CancelFunc
	shutdown bool

	mu sync.Mutex
	wg sync.WaitGroup
}

func NewService(repo *repository.Repository, pool *pool.Pool, queueCapacity int) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		repository: repo,
		pool:       pool,
		taskQueue:  make(chan string, queueCapacity),
		limiter:    make(chan struct{}, queueCapacity),

		ctx:      ctx,
		cancel:   cancel,
		shutdown: false,
	}
	s.pool.Create()

	s.wg.Add(1)
	go s.proccess()

	return s
}

func (s *Service) Create(ctx context.Context, path string) error {

	s.mu.Lock()
	if s.shutdown {
		s.mu.Unlock()
		return errors.ErrServerShuttingDown
	}
	s.mu.Unlock()

	select {
	case s.taskQueue <- path:
		//запись в репо
		return nil
	case <-ctx.Done(): //контекст отслеживающий отмену запроса
		return ctx.Err()
	case <-s.ctx.Done(): //shutdown
		return errors.ErrServerShuttingDown
	default:
		return errors.ErrQueueTasksFull
	}
}

func (s *Service) proccess() {
	defer s.wg.Done()
	for {
		select {
		case imagePath, ok := <-s.taskQueue:
			if !ok {
				fmt.Println("task queue channel was closed")
				return
			}
			select {
			case s.limiter <- struct{}{}:
				s.wg.Add(1)
				time.Sleep(1 * time.Second)
				go func(path string) {
					defer func() {
						<-s.limiter
						s.wg.Done()
					}()
					result := s.pool.Handle(path)
					fmt.Println("Procces result ", <-result)
				}(imagePath)
			case <-s.ctx.Done():
				fmt.Println("Skipping task due to shutdown:", imagePath)
			}
		case <-s.ctx.Done():
			return
		}
	}

}

func (s *Service) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.shutdown = true
	s.mu.Unlock()

	s.cancel()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	close(s.taskQueue)
	s.wg.Wait()

	s.pool.Shutdown()
	select {
	case <-done:
		fmt.Println("Service was gracefull shutdown")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// for imagePath := range s.taskQueue {
// 	s.limiter <- struct{}{}
// 	go func(path string) {
// 		defer func() { <-s.limiter }()
// 		result := s.pool.Handle(path)
// 		fmt.Println("Procces result ", <-result)
// 	}(imagePath)

// }
