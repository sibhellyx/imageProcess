package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sibhellyx/imageProccesor/internal/errors"
	"github.com/sibhellyx/imageProccesor/internal/models"
	"github.com/sibhellyx/imageProccesor/internal/repository"
	"github.com/sibhellyx/imageProccesor/pkg/actions"
)

type WorkerPool interface {
	Create()
	Handle(imagePath *models.ImageTask) <-chan error
	Wait()
	Shutdown()
	Stats()
}

type Service struct {
	repository *repository.Repository
	pool       WorkerPool
	taskQueue  chan *models.ImageTask
	limiter    chan struct{}

	ctx      context.Context
	cancel   context.CancelFunc
	shutdown bool

	mu sync.Mutex
	wg sync.WaitGroup
}

func NewService(repo *repository.Repository, pool WorkerPool, queueCapacity int) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		repository: repo,
		pool:       pool,
		taskQueue:  make(chan *models.ImageTask, queueCapacity),
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

func (s *Service) AddImageTask(req models.ImageRequestAction) (*models.ImageTask, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	imageTask := &models.ImageTask{
		Name:         "",
		DownloadPath: "",
		Path:         req.Path,
		Status:       models.StatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Actions:      req.Actions,
	}

	return imageTask, nil
}

func (s *Service) Download(req models.ImageRequestDownload) (string, error) {
	err := req.Validate()
	if err != nil {
		return "", err
	}

	path, err := actions.DownloadImageWithResty(req.Url, req.Name)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (s *Service) Proccess(ctx context.Context, imageTask *models.ImageTask) error {
	s.mu.Lock()
	if s.shutdown {
		s.mu.Unlock()
		return errors.ErrServerShuttingDown
	}
	s.mu.Unlock()
	select {
	case s.taskQueue <- imageTask:
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
		case imageTask, ok := <-s.taskQueue:
			if !ok {
				fmt.Println("task queue channel was closed")
				return
			}
			select {
			case s.limiter <- struct{}{}:
				s.wg.Add(1)
				imageTask.Status = models.StatusProcessing
				go func(imageTask *models.ImageTask) {
					defer func() {
						<-s.limiter
						s.wg.Done()
					}()
					result := s.pool.Handle(imageTask)
					if <-result != nil {
						imageTask.Status = models.StatusFailed
					} else {
						imageTask.Status = models.StatusCompleted
					}
				}(imageTask)
			case <-s.ctx.Done():

				fmt.Println("Skipping task due to shutdown:", imageTask.Path)
				imageTask.Status = models.StatusCanceled
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

	s.pool.Shutdown()
	select {
	case <-done:
		fmt.Println("Service was gracefull shutdown")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
