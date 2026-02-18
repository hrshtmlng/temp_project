package main

import (
	"context"
	"sync"
)

func RunWorkerPool(ctx context.Context, numJobs, numWorkers int) []Result {
	jobs := make(chan int, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(ctx, w, jobs, results, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var collected []Result
	for res := range results {
		collected = append(collected, res)
	}

	return collected

}
