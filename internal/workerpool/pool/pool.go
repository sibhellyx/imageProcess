package pool

import (
	"context"
	"fmt"
	"sync"

	"github.com/sibhellyx/imageProccesor/internal/workerpool/worker"
)

type Pool struct {
	pool     chan *worker.Worker
	handler  func(int, string) string
	workers  []*worker.Worker
	ctx      context.Context
	cancel   context.CancelFunc
	mu       *sync.Mutex
	shutdown bool
	wg       sync.WaitGroup
}

func NewPool(handler func(int, string) string, countWorkers int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool{
		handler:  handler,
		pool:     make(chan *worker.Worker, countWorkers),
		workers:  make([]*worker.Worker, 0, countWorkers),
		ctx:      ctx,
		cancel:   cancel,
		mu:       &sync.Mutex{},
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

func (p *Pool) Handle(imagePath string) <-chan string {
	p.mu.Lock()
	if p.shutdown {
		p.mu.Unlock()
		ch := make(chan string)
		close(ch)
		return ch
	}
	p.mu.Unlock()

	resultChan := make(chan string, 1)

	select {
	case w := <-p.pool:
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			defer func() { p.pool <- w }()

			result := p.handler(w.Id, imagePath)
			w.JobsCompleted++
			resultChan <- result

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
	fmt.Println("Worker pool was graceful shutdown")
}

func (p *Pool) Stats() {
	for _, w := range p.workers {
		fmt.Println("Worker ", w.Id, "processed ", w.JobsCompleted)
	}
}
