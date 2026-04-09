package concurrency

import (
	"context"
	"time"

	"github.com/naghinezhad/BookingResourceSystem/internal/service"
)

type Job struct {
	Ctx        context.Context
	ResourceID string
	Start      time.Time
	End        time.Time
	ResultChan chan error
}

type WorkerPool struct {
	jobs chan Job
}

func NewWorkerPool(workerCount int, queueSize int, service *service.ReservationService) *WorkerPool {
	pool := &WorkerPool{
		jobs: make(chan Job, queueSize),
	}

	for i := range workerCount {
		go pool.worker(i, service)
	}

	return pool
}

func (p *WorkerPool) worker(_ int, s *service.ReservationService) {
	for job := range p.jobs {
		err := s.Reserve(job.Ctx, job.ResourceID, job.Start, job.End)
		job.ResultChan <- err
	}
}

func (p *WorkerPool) Submit(job Job) bool {
	select {
	case p.jobs <- job:
		return true
	default:
		return false
	}
}
