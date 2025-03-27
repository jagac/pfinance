package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/jagac/pfinance/pkg/cache"
)

// MiddlewareFunc is worker middleware
type MiddlewareFunc func(Job) Job

// Job is what's going to be run on background
type Job func(c context.Context) (any, error)

// Task represents the Name and Job to be run on background
type Task struct {
	OriginContext context.Context
	Name          string
	Job           Job
	TTL           time.Duration
}

type TaskResult struct {
	Value any   `json:"value"`
	Error error `json:"error"`
}

// Worker is a process that runs tasks
type Worker interface {
	Run(id string)
	Enqueue(task Task)
	Use(middleware MiddlewareFunc)
	Length() int64
	Shutdown(ctx context.Context) error
	GetResult(taskName string) (TaskResult, bool)
}

// BackgroundWorker is a worker that runs tasks on background
type BackgroundWorker struct {
	context.Context
	logger      *slog.Logger
	queue       chan Task
	len         int64
	middleware  MiddlewareFunc
	resultCache *cache.Cache[string, TaskResult]
	sync.RWMutex
}

var maxQueueSize = 100

func New(logger *slog.Logger, cache *cache.Cache[string, TaskResult]) *BackgroundWorker {
	ctx := context.Background()

	return &BackgroundWorker{
		Context:     ctx,
		logger:      logger,
		queue:       make(chan Task, maxQueueSize),
		resultCache: cache,
		middleware: func(next Job) Job {
			return next
		},
	}
}

// Run initializes the worker loop
func (w *BackgroundWorker) Run(workerID string) {
	w.logger.Info(workerID, "started", "listening...")
	for task := range w.queue {
		w.logger.Info(workerID, "Running task", task.Name)
		result, err := w.middleware(task.Job)(task.OriginContext)
		if err != nil {
			w.logger.Error(err.Error())
		}
		// Use the TTL from the task when storing the result in cache
		ttl := task.TTL
		if ttl == 0 {
			ttl = 10 * time.Minute
		}

		if err := w.resultCache.Set(task.Name, TaskResult{Value: result, Error: err}, ttl); err != nil {
			w.logger.Error(err.Error())
		}

		w.Lock()
		w.len--
		w.Unlock()
		w.logger.Info(workerID, "Task completed", task.Name)
	}
}

// Shutdown current worker
func (w *BackgroundWorker) Shutdown(ctx context.Context) error {
	if w.Length() > 0 {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			count := w.Length()
			if count == 0 {
				return nil
			}

			select {
			case <-ctx.Done():
				return errors.New("timeout waiting for worker queue")
			case <-ticker.C:
			}
		}
	}
	return nil
}

// Enqueue a task on current worker with a TTL
func (w *BackgroundWorker) Enqueue(task Task) {
	w.logger.Info("New task enqueued", "Task", task.Name)
	w.Lock()
	w.len++
	w.Unlock()
	w.queue <- task
}

// Length from current queue length
func (w *BackgroundWorker) Length() int64 {
	w.RLock()
	defer w.RUnlock()
	return w.len
}

// Use this to inject worker dependencies
func (w *BackgroundWorker) Use(middleware MiddlewareFunc) {
	w.middleware = middleware
}

func (w *BackgroundWorker) GetResult(taskName string) (TaskResult, bool) {
	return w.resultCache.Get(taskName)
}
