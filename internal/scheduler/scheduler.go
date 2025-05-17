package scheduler

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"go-sync-s3.miljanilic.com/internal/config"
	"go-sync-s3.miljanilic.com/internal/worker"
	"log"
)

func ScheduleTasks(ctx context.Context, tasks map[string]config.Task, buckets map[string]config.Bucket) (*cron.Cron, error) {
	c := cron.New()

	for id, task := range tasks {
		id, task := id, task // Capture range variable!

		if b, ok := buckets[task.BucketID]; ok {
			task.Bucket = &b
		} else {
			log.Printf("task[%s]: unknown bucket ID '%s'", id, task.BucketID)
			continue
		}

		_, err := c.AddFunc(task.Cron, func() {
			if ctx.Err() != nil {
				log.Printf("task[%s]: skipped due to shutdown", id)
				return
			}

			if err := worker.Execute(ctx, task); err != nil {
				log.Printf("task[%s]: error: %v", id, err)
			}
		})
		if err != nil {
			return nil, fmt.Errorf("task[%s]: failed to schedule: %w", id, err)
		}

		log.Printf("task[%s]: scheduled task with cron '%s'", id, task.Cron)
	}

	c.Start()
	return c, nil
}
