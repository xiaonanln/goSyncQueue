//The MIT License (MIT)
//
//Copyright (c) 2013 Evan Huus
//
//Permission is hereby granted, free of charge, to any person obtaining a copy of
//this software and associated documentation files (the "Software"), to deal in
//the Software without restriction, including without limitation the rights to
//use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
//the Software, and to permit persons to whom the Software is furnished to do so,
//subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
//FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
//COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
//IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
//CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

func (ch *syncQueueWithGoroutine) Close() {
	close(ch.input)
}
