package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

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
			fmt.Printf("Worker %d stopping: %v\n", id, ctx.Err())
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Printf("Worker %d processing job %d\n", id, job)

			// simulate work
			time.Sleep(time.Millisecond * time.Duration(300+rand.Intn(400)))
			results <- Result{JobID: job, Value: job * 2}
		}
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	numJobs := 20
	numWorker := 4

	jobs := make(chan int, numJobs)
	results := make(chan Result, numWorker)

	// cancel everything after 1 sec
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// start workers
	for w := 1; w <= numWorker; w++ {
		wg.Add(1)
		go worker(ctx, w, jobs, results, &wg)
	}

	// send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// close results after workers exit
	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Printf("Result: job %d → %d\n", res.JobID, res.Value)
	}
	fmt.Println("All done")
}
