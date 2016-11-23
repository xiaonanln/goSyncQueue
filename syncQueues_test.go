package sync_queue

import (
	"math/rand"
	"testing"
)

const (
	SEQ_TEST_N   = 100000
	FUZZY_TEST_N = 100000
)

func TestSyncQueueWithGoroutine_Seq(t *testing.T) {
	q := NewSyncQueueWithGoroutine()
	seqTestSyncQueue(t, q)
}

func TestSyncQueueWithGoroutine_Fuzzy(t *testing.T) {
	q := NewSyncQueueWithGoroutine()
	fuzzyTestSyncQueue(t, q)
}

func TestSyncQueue_Seq(t *testing.T) {
	q := NewSyncQueue()
	seqTestSyncQueue(t, q)
}

func TestSyncQueue_Fuzzy(t *testing.T) {
	q := NewSyncQueue()
	fuzzyTestSyncQueue(t, q)
}

func TestSyncQueueByChan_Seq(t *testing.T) {
	q := NewSyncQueueByChan()
	seqTestSyncQueue(t, q)
}

func TestSyncQueueByChan_Fuzzy(t *testing.T) {
	q := NewSyncQueueByChan()
	fuzzyTestSyncQueue(t, q)
}

func seqTestSyncQueue(t *testing.T, q SyncQueue) {
	vals := []interface{}{}
	for i := 0; i < SEQ_TEST_N; i++ {
		vals = append(vals, rand.Int())
	}

	for i, val := range vals {
		q.Push(val)
		if q.Len() != -1 && q.Len() != i+1 {
			t.Fatalf("queue length should be %v, but is %v", i+1, q.Len())
		}
	}

	for i := 0; i < SEQ_TEST_N; i++ {
		val := q.Pop()
		if val != vals[i] {
			t.Fatalf("pop val should be %v, but is %v", vals[i], val)
		}
		if q.Len() != -1 && q.Len() != SEQ_TEST_N-i-1 {
			t.Fatalf("queue length should be %v, but is %v", SEQ_TEST_N-i-1, q.Len())
		}
	}
	q.Close()
}

func fuzzyTestSyncQueue(t *testing.T, q SyncQueue) {
	vals := []interface{}{}

	for i := 0; i < FUZZY_TEST_N; i++ {
		if q.Len() > 0 && rand.Float64() < 0.4 {
			v := q.Pop()
			if v != vals[0] {
				t.Fatalf("pop val should be %v, but is %v", vals[i], v)
			}
			vals = vals[1:]
		} else {
			v := rand.Int()
			vals = append(vals, v)
			q.Push(v)
		}

		if q.Len() != -1 && q.Len() != len(vals) {
			t.Fatalf("queue length should be %v, but is %v", len(vals), q.Len())
		}
	}

	for _, val := range vals {
		pv := q.Pop()
		if val != pv {
			t.Fatalf("pop val should be %v, but is %v", val, pv)
		}
	}
	q.Close()
}
