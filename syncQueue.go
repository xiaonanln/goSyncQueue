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

	q.lock.Lock()
	for buffer.Length() == 0 {
		c.Wait()
	}

	v := buffer.Peek()
	buffer.Remove()

	q.lock.Unlock()
	return v
}

func (q *syncQueue) TryPop() (interface{}, bool) {
	buffer := q.buffer

	q.lock.Lock()

	if buffer.Length() > 0 {
		v := buffer.Peek()
		buffer.Remove()
		q.lock.Unlock()
		return v, true
	} else {
		q.lock.Unlock()
		return nil, false
	}
}

func (q *syncQueue) Push(v interface{}) {
	q.lock.Lock()
	q.buffer.Add(v)
	q.lock.Unlock()
}

func (q *syncQueue) Len() (l int) {
	q.lock.Lock()
	l = q.buffer.Length()
	q.lock.Unlock()
	return
}

func (q *syncQueue) Close() {
}
