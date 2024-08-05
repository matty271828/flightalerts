package google

import (
	"log"
	"time"
)

type WriteRequest struct {
	Function   func() error
	RetryCount int
}

type Queuer struct {
	queue     chan WriteRequest
	errorChan chan error
}

// NewQueuer is used to create a buffered channel in order
// to provide a queue for write requests to google sheets.
func NewQueuer(bufferSize int) *Queuer {
	q := &Queuer{
		queue:     make(chan WriteRequest, bufferSize), // Adjust buffer size as needed
		errorChan: make(chan error),
	}

	go q.worker()

	return q
}

func (q *Queuer) worker() {
	for req := range q.queue {
		err := req.Function()
		if err != nil {
			if req.RetryCount < 3 {
				req.RetryCount++
				q.queue <- req
			} else {
				log.Printf("Error executing request after retries: %v", err)
				q.errorChan <- err
			}
		}

		// Introduce a delay to throttle requests
		time.Sleep(time.Second)
	}
}

// QueueWork is used to add a function that returns an error to the queue.
func (q *Queuer) QueueWork(work func() error) {
	writeRequest := WriteRequest{
		Function: work,
	}
	q.queue <- writeRequest
}

// ErrorChan returns the error channel for the queuer.
func (q *Queuer) ErrorChan() <-chan error {
	return q.errorChan
}

// IsQueueFull checks if the queue is full.
func (q *Queuer) IsQueueFull() bool {
	return len(q.queue) == cap(q.queue)
}
