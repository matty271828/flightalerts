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
	queue       chan WriteRequest
	workerCount int
}

// NewQueuer is used to create a buffered channel in order
// to provide a queue for write requests to google sheets.
func NewQueuer(workerCount int) *Queuer {
	q := &Queuer{
		queue:       make(chan WriteRequest, 100), // Adjust buffer size as needed
		workerCount: workerCount,
	}

	for i := 0; i < workerCount; i++ {
		go q.worker()
	}

	return q
}

func (q *Queuer) worker() {
	rateLimiter := time.NewTicker(time.Minute / 60) // 60 writes per minute
	defer rateLimiter.Stop()

	for req := range q.queue {
		<-rateLimiter.C
		err := req.Function()
		if err != nil {
			if req.RetryCount < 3 { // Retry up to 3 times
				req.RetryCount++
				q.queue <- req
			} else {
				log.Printf("Error executing request after retries: %v", err)
			}
		}
	}
}
