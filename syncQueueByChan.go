package sync_queue

type syncQueueByChannel struct {
	channel chan interface{}
}

func NewSyncQueueByChan() SyncQueue {
	ch := &syncQueueByChannel{
		channel: make(chan interface{}, 1000000),
	}
	return ch
}

func (q *syncQueueByChannel) Pop() interface{} {
	return <-q.channel
}

func (q *syncQueueByChannel) TryPop() (interface{}, bool) {
	select {
	case v := <-q.channel:
		return v, true
	default:
		return nil, false
	}
}

func (q *syncQueueByChannel) Push(v interface{}) {
	q.channel <- v
}

func (q *syncQueueByChannel) Len() int {
	return -1
}

func (q *syncQueueByChannel) Close() {
	close(q.channel)
}
