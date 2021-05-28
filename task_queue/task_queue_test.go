package taskqueue

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestTaskQueue(t *testing.T) {
	c := NewTaskQueue(TaskQueueConfig{
		MaxTaskCount: 5,
		WorkerCount:  1,
	})
	c.Enqueue(func() { fmt.Println("A") })
	c.Enqueue(func() { fmt.Println("B") })
	time.Sleep(2 * time.Second)
	assert.Equal(t, c.Config().MaxTaskCount, 5)
}
