package sync_queue

import (
	"math/rand"
	"testing"
)

const (
	SEQ_TEST_N = 1000000
)

func TestSyncQueueWithGoroutine(t *testing.T) {
	q := NewSyncQueueWithGoroutine()
	testSyncQueue(t, q)
}

func TestSyncQueue(t *testing.T) {
	q := NewSyncQueue()
	testSyncQueue(t, q)
}

func testSyncQueue(t *testing.T, q SyncQueue) {
	vals := []interface{}{}
	for i := 0; i < SEQ_TEST_N; i++ {
		vals = append(vals, rand.Int())
	}

	for _, val := range vals {
		q.Push(val)
	}
	for i := 0; i < SEQ_TEST_N; i++ {
		val := q.Pop()
		if val != vals[i] {
			t.FailNow()
		}
	}
}
