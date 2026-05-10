// Package worker contains the background video processing pipeline:
// queue → download → FFmpeg HLS transcode → upload → asset registration.
package worker

import (
	"errors"
	"sync"
)

// Job describes a single video to be processed.
type Job struct {
	VideoID     string
	StoragePath string // path inside the raw-videos bucket
	OwnerID     string
}

// Queue is the interface used by services to enqueue work. The default
// implementation is in-memory; swap with a Redis-backed queue for HA.
type Queue interface {
	Enqueue(j Job) error
	Dequeue() (Job, bool)
	Len() int
	Close()
}

type MemoryQueue struct {
	mu     sync.Mutex
	cond   *sync.Cond
	items  []Job
	closed bool
}

func NewMemoryQueue() *MemoryQueue {
	q := &MemoryQueue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

var ErrQueueClosed = errors.New("queue closed")

func (q *MemoryQueue) Enqueue(j Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return ErrQueueClosed
	}
	q.items = append(q.items, j)
	q.cond.Signal()
	return nil
}

// Dequeue blocks until an item is available. Returns false when the queue is closed and drained.
func (q *MemoryQueue) Dequeue() (Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 && !q.closed {
		q.cond.Wait()
	}
	if len(q.items) == 0 {
		return Job{}, false
	}
	j := q.items[0]
	q.items = q.items[1:]
	return j, true
}

func (q *MemoryQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

func (q *MemoryQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
}
