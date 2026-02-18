package main

import (
	"context"
	"testing"
	"time"
)

func TestRunWorkerPool_Completes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results := RunWorkerPool(ctx, 5, 3)

	if len(results) == 0 {
		t.Fatal("expected some results, got none")
	}
}

func TestRunWorkerPool_ResultCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	numJobs := 10
	results := RunWorkerPool(ctx, numJobs, 3)

	if len(results) != numJobs {
		t.Fatalf("expected %d results, got %d", numJobs, len(results))
	}
}

func TestJobError_Type(t *testing.T) {
	err := JobError{JobID: 1, Err: ErrTemporary}

	if err.JobID != 1 {
		t.Fatal("wrong job id")
	}

	if err.Err != ErrTemporary {
		t.Fatal("Wrong error")
	}
}

func TestRunWorkerPool_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	results := RunWorkerPool(ctx, 50, 5)

	if len(results) == 50 {
		t.Fatal("expected some jobs to be cancelled")
	}
}

func TestRunWorkerPool_ZeroWorkers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	results := RunWorkerPool(ctx, 5, 0)

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestJobError_ErrorString(t *testing.T) {
	err := JobError{JobID: 3, Err: ErrTemporary}

	msg := err.Error()

	if msg == "" {
		t.Fatal("expected error message")
	}

}
