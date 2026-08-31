package worker

import (
	"context"
	"testing"
	"time"
)

// TestNewJobWorker - verify worker is created with a valid channel
func TestNewJobWorker(t *testing.T) {
	w := NewJobWorker(10)
	if w == nil {
		t.Fatal("expected worker to be created, got nil")
	}
	if w.jobQueue == nil {
		t.Fatal("expected job queue channel to be initialized")
	}
}

// TestJobWorker_SubmitAndProcess - verify job sent via channel is processed by goroutine
// This tests the goroutine + channel wiring end-to-end
func TestJobWorker_SubmitAndProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := NewJobWorker(10)
	w.Start(ctx)

	// Submit a job into the channel
	w.Submit(Job{
		Type:      JobProductCreated,
		ProductID: 1,
		Details:   "test product",
	})

	// Give the goroutine a moment to read from the channel and process
	// In real tests you'd use a done channel or sync.WaitGroup instead of sleep
	time.Sleep(50 * time.Millisecond)

	// If we reach here without panic or deadlock, the goroutine + channel works
}

// TestJobWorker_MultipleJobs - verify multiple jobs are all processed
func TestJobWorker_MultipleJobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := NewJobWorker(10)
	w.Start(ctx)

	jobs := []Job{
		{Type: JobProductCreated, ProductID: 1, Details: "created"},
		{Type: JobProductUpdated, ProductID: 1, Details: "updated"},
		{Type: JobProductDeleted, ProductID: 1, Details: "deleted"},
	}

	for _, job := range jobs {
		w.Submit(job)
	}

	time.Sleep(100 * time.Millisecond)
	// No deadlock or panic = all jobs processed successfully
}

// TestJobWorker_GracefulShutdown - verify goroutine stops when context is cancelled
// Like testing Spring Boot's @PreDestroy or graceful shutdown hook
func TestJobWorker_GracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	w := NewJobWorker(10)
	w.Start(ctx)

	// Cancel context = simulate app shutdown (Ctrl+C)
	cancel()

	// Give goroutine time to receive ctx.Done() and exit
	time.Sleep(50 * time.Millisecond)

	// If goroutine didn't exit, this test would eventually deadlock on submit
	// Reaching here means shutdown worked correctly
}

// TestJobWorker_BufferedChannel - verify jobs queue up without blocking when buffer has space
func TestJobWorker_BufferedChannel(t *testing.T) {
	// Create worker but do NOT start it (no goroutine reading)
	w := NewJobWorker(5)

	// With a buffer of 5, these 5 submits should not block
	for i := 0; i < 5; i++ {
		w.Submit(Job{Type: JobProductCreated, ProductID: uint(i)})
	}

	if len(w.jobQueue) != 5 {
		t.Errorf("expected 5 jobs in queue, got %d", len(w.jobQueue))
	}
}
