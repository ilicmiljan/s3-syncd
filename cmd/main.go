package main

import (
	"context"
	"errors"
	"fmt"
	"go-sync-s3.miljanilic.com/internal/config"
	"go-sync-s3.miljanilic.com/internal/scheduler"
	"go-sync-s3.miljanilic.com/internal/validation"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buckets, err := config.LoadBuckets("config/buckets.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	tasks, err := config.LoadTasks("config/tasks.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := validation.ValidateBuckets(buckets); err != nil {
		var validationErrors *validation.BucketValidationError

		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors.Errors {
				fmt.Printf("error in bucket '%s', field '%s': %s\n", e.BucketID, e.Field, e.Message)
			}
		}
		os.Exit(1)
	}

	if err := validation.ValidateTasks(tasks, buckets); err != nil {
		var validationErrors *validation.TaskValidationError
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors.Errors {
				fmt.Printf("error in task '%s', field '%s': %s\n", e.TaskID, e.Field, e.Message)
			}
		}
		os.Exit(1)
	}

	cron, err := scheduler.ScheduleTasks(ctx, tasks, buckets)
	if err != nil {
		log.Fatalf("failed to schedule tasks: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("received signal: %s. shutting down...", sig)
		cancel()
	}()

	<-ctx.Done() // Wait for context to be cancelled
	cron.Stop()  // Stop scheduler after cancellation
	log.Println("shutdown complete.")
}
