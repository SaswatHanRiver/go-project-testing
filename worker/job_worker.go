package worker

import (
	"context"
	"log/slog"
)

// JobType - what kind of event happened
type JobType string

const (
	JobProductCreated JobType = "PRODUCT_CREATED"
	JobProductUpdated JobType = "PRODUCT_UPDATED"
	JobProductDeleted JobType = "PRODUCT_DELETED"
)

// Job - the message sent through the channel
// Think of this like a message in a Java BlockingQueue or Spring ApplicationEvent
type Job struct {
	Type      JobType
	ProductID uint
	Details   string
}

// JobWorker holds the channel and runs in the background
// In Spring Boot you'd use @Async or a ThreadPoolTaskExecutor
// In Go: one goroutine + one channel is all you need
type JobWorker struct {
	jobQueue chan Job // channel = thread-safe message queue
}

// NewJobWorker creates a worker with a buffered channel
// bufferSize = how many jobs can wait before the sender blocks
// Like setting the queue capacity in a Java ThreadPoolExecutor
func NewJobWorker(bufferSize int) *JobWorker {
	return &JobWorker{
		jobQueue: make(chan Job, bufferSize),
	}
}

// Start launches the background goroutine
// "go func()" = fire and forget goroutine (like new Thread().start() in Java)
// context.Context handles graceful shutdown - when ctx is cancelled, goroutine exits
func (w *JobWorker) Start(ctx context.Context) {
	go func() {
		slog.Info("Background job worker started")
		for {
			select {
			case job := <-w.jobQueue:
				// This runs every time a job arrives in the channel
				w.processJob(job)
			case <-ctx.Done():
				// ctx.Done() fires when the app is shutting down
				// Like @PreDestroy in Spring Boot
				slog.Info("Job worker shutting down gracefully")
				return
			}
		}
	}()
}

// Submit sends a job into the channel (non-blocking if buffer has space)
// Called from controllers after each CRUD operation
func (w *JobWorker) Submit(job Job) {
	w.jobQueue <- job
}

// processJob is the actual work done per event
// Right now it just logs — in a real app: send email, update cache, call another service
func (w *JobWorker) processJob(job Job) {
	slog.Info("Job processed",
		"type", job.Type,
		"productID", job.ProductID,
		"details", job.Details,
	)
}
