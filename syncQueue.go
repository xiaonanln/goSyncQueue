package sync_queue

import (
	"sync"

	"gopkg.in/eapache/queue.v1"
)

type syncQueue struct {
	lock    sync.Mutex
	popable *sync.Cond
	buffer  *queue.Queue
}

func NewSyncQueue() SyncQueue {
	ch := &syncQueue{
		buffer: queue.New(),
	}
	ch.popable = sync.NewCond(&ch.lock)
	return ch
}

func (q *syncQueue) Pop() interface{} {
	c := q.popable
	buffer := q.buffer

	c.L.Lock()
	for buffer.Length() == 0 {
		c.Wait()
	}

	c.L.Unlock()
	return nil
}

func (q *syncQueue) TryPop() (interface{}, bool) {
	return nil, false
}

func (q *syncQueue) Push(v interface{}) {
	q.buffer.Add(v)
}
