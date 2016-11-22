package sync_queue

import "gopkg.in/eapache/queue.v1"

// InfiniteChannel implements the Channel interface with an infinite buffer between the input and the output.
type syncQueueWithGoroutine struct {
	input, output chan interface{}
	length        chan int
	buffer        *queue.Queue
}

func NewSyncQueueWithGoroutine() SyncQueue {
	ch := &syncQueueWithGoroutine{
		input:  make(chan interface{}),
		output: make(chan interface{}),
		length: make(chan int),
		buffer: queue.New(),
	}
	go ch.infiniteBuffer()
	return ch
}

func (ch *syncQueueWithGoroutine) Pop() interface{} {
	return <-ch.output
}

func (ch *syncQueueWithGoroutine) TryPop() (interface{}, bool) {
	select {
	case v := <-ch.output:
		return v, true
	default:
		return nil, false
	}
}

func (ch *syncQueueWithGoroutine) Push(v interface{}) {
	ch.input <- v
}

func (ch *syncQueueWithGoroutine) Len() int {
	return <-ch.length
}

func (ch *syncQueueWithGoroutine) infiniteBuffer() {
	var input, output chan interface{}
	var next interface{}
	input = ch.input

	for input != nil || output != nil {
		select {
		case elem, open := <-input:
			if open {
				ch.buffer.Add(elem)
			} else {
				input = nil
			}
		case output <- next:
			ch.buffer.Remove()
		case ch.length <- ch.buffer.Length():
		}

		if ch.buffer.Length() > 0 {
			output = ch.output
			next = ch.buffer.Peek()
		} else {
			output = nil
			next = nil
		}
	}

	close(ch.output)
	close(ch.length)
}
