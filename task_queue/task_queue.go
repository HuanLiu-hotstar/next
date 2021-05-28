package taskqueue

import "time"

type Func func()
type TaskQueue struct {
	config TaskQueueConfig
	queue  chan Func
}

type TaskQueueConfig struct {
	MaxTaskCount int
	WorkerCount  int
}

func (t *TaskQueue) Config() TaskQueueConfig {
	return t.config
}
func (t *TaskQueueConfig) Validate() TaskQueueConfig {
	if t.MaxTaskCount < 1 {
		t.MaxTaskCount = 1
	}
	if t.WorkerCount < 1 {
		t.WorkerCount = 1
	}
	return *t
}

func NewTaskQueue(config TaskQueueConfig) *TaskQueue {
	t := &TaskQueue{
		config: config.Validate(),
		queue:  make(chan Func, config.MaxTaskCount),
	}
	t.start()
	return t
}

func (t *TaskQueue) Enqueue(f Func) {
	t.queue <- f
}

func (t *TaskQueue) EnqueueTimeout(timeout time.Duration, f Func) bool {
	select {
	case t.queue <- f:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (t *TaskQueue) TryEnqueue(f Func) bool {
	select {
	case t.queue <- f:
		return true
	default:
		return false
	}
}

func (t *TaskQueue) start() {
	for i := 0; i < t.config.WorkerCount; i++ {
		go func() {
			for f := range t.queue {
				f()
			}
		}()
	}
}
