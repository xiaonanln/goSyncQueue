package sync_queue

type SyncQueue interface {
	Pop() interface{}
	TryPop() (interface{}, bool)
	Push(v interface{})
	Close()
	Len() int
}
