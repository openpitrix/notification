package etcdutil

import "github.com/coreos/etcd/contrib/recipes"

type Queue struct {
	*recipe.Queue
}


func (q *Queue) Enqueue(val string) error {
	return q.Queue.Enqueue(val)
}

// Dequeue returns Enqueue()'d elements in FIFO order. If the
// queue is empty, Dequeue blocks until elements are available.
func (q *Queue) Dequeue() (string, error) {
	return q.Queue.Dequeue()
}
