package pool

import (
	"fmt"

	"github.com/sibhellyx/imageProccesor/internal/workerpool/worker"
)

type Pool struct {
	pool    chan *worker.Worker
	handler func(int, string) string
	workers []*worker.Worker
}

func NewPool(handler func(int, string) string, countWorkers int) *Pool {
	p := &Pool{
		handler: handler,
		pool:    make(chan *worker.Worker, countWorkers),
		workers: make([]*worker.Worker, 0, countWorkers),
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
	resultChan := make(chan string, 1)
	w := <-p.pool
	go func() {
		resultChan <- p.handler(w.Id, imagePath)
		w.JobsCompleted++
		p.pool <- w
	}()
	return resultChan
}

func (p *Pool) Wait() {
	for range len(p.workers) {
		<-p.pool
	}
}

func (p *Pool) Stats() {
	for _, w := range p.workers {
		fmt.Println("Worker ", w.Id, "proccesed ", w.JobsCompleted)
	}
}
