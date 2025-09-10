package pool

import (
	"context"
	"fmt"
	"sync"

	"github.com/sibhellyx/imageProccesor/internal/models"
	"github.com/sibhellyx/imageProccesor/internal/workerpool/worker"
)

type Pool struct {
	pool    chan *worker.Worker
	handler func(int, *models.ImageTask) error
	workers []*worker.Worker

	ctx      context.Context
	cancel   context.CancelFunc
	shutdown bool

	mu sync.Mutex
	wg sync.WaitGroup
}

func NewPool(handler func(int, *models.ImageTask) error, countWorkers int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool{
		handler:  handler,
		pool:     make(chan *worker.Worker, countWorkers),
		workers:  make([]*worker.Worker, 0, countWorkers),
		ctx:      ctx,
		cancel:   cancel,
		shutdown: false,
	}

	for i := 0; i < countWorkers; i++ {
		p.workers = append(p.workers, &worker.Worker{
			Id: i,
		})
	}
	return p
}

func (p *Pool) Create() {
	for _, w := range p.workers {
		p.pool <- w
	}
}

func (p *Pool) Handle(imagePath *models.ImageTask) <-chan error {
	p.mu.Lock()
	if p.shutdown {
		p.mu.Unlock()
		ch := make(chan error)
		close(ch)
		return ch
	}
	p.mu.Unlock()

	resultChan := make(chan error, 1)

	select {
	case w := <-p.pool:
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			defer func() { p.pool <- w }()

			err := p.handler(w.Id, imagePath)
			w.JobsCompleted++
			resultChan <- err

		}()
	case <-p.ctx.Done():
		close(resultChan)
	}

	return resultChan
}

func (p *Pool) Wait() {
	for i := 0; i < len(p.workers); i++ {
		<-p.pool
	}
}

func (p *Pool) Shutdown() {
	p.mu.Lock()
	p.shutdown = true
	p.mu.Unlock()
	p.cancel()
	p.wg.Wait()
	p.Wait()
	close(p.pool)
	// fmt.Println("Worker pool was graceful shutdown")
}

func (p *Pool) Stats() {
	for _, w := range p.workers {
		fmt.Println("Worker ", w.Id, "processed ", w.JobsCompleted)
	}
}
