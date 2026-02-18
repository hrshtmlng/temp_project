package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var ErrTemporary = errors.New("temporary failure")

type JobError struct {
	JobID int
	Err   error
}

func (e JobError) Error() string {
	return fmt.Sprintf("job %d: %v", e.JobID, e.Err)
}

type Result struct {
	JobID int
	Value int
	Err   error
}

func worker(ctx context.Context, id int, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			// simulate work
			time.Sleep(time.Millisecond * time.Duration(200+rand.Intn(300)))

			// random failure
			if rand.Intn(4) == 0 {
				err := JobError{
					JobID: job,
					Err:   fmt.Errorf("%w: network glitch", ErrTemporary),
				}
				results <- Result{JobID: job, Err: err}
				continue
			}
			results <- Result{JobID: job, Value: job * 2}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	results := RunWorkerPool(ctx, 20, 4)

	for _, res := range results {
		if res.Err != nil {
			log.Printf("ERROR job_id=%d err=%v\n", res.JobID, res.Err)
			continue
		}
		log.Printf("INFO job completed job_id=%d result=%d\n", res.JobID, res.Value)
	}

	log.Println("All workers finished")
}
