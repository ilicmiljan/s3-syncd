package validation

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/robfig/cron/v3"
	"go-sync-s3.miljanilic.com/internal/config"
)

func (e *TaskValidationError) Error() string {
	var lines []string
	for _, err := range e.Errors {
		lines = append(lines, fmt.Sprintf("task[%s]: %s (%s)", err.TaskID, err.Message, err.Field))
	}
	return "task validation failed:\n" + strings.Join(lines, "\n")
}

func ValidateTasks(tasks map[string]config.Task, buckets map[string]config.Bucket) error {
	var errs []TaskFieldError

	for id, t := range tasks {
		if t.ID == "" {
			errs = append(errs, TaskFieldError{id, "id", "missing required field"})
		}

		if t.BucketID == "" {
			errs = append(errs, TaskFieldError{id, "bucket", "missing required field"})
		} else if _, ok := buckets[t.BucketID]; !ok {
			errs = append(errs, TaskFieldError{id, "bucket", fmt.Sprintf("unknown bucket_id '%s'", t.BucketID)})
		}

		if t.Cron == "" {
			errs = append(errs, TaskFieldError{id, "cron", "missing required field"})
		} else if _, err := cron.ParseStandard(t.Cron); err != nil {
			errs = append(errs, TaskFieldError{id, "cron", fmt.Sprintf("invalid cron expression '%s': %v", t.Cron, err)})
		}

		if t.Remote == "" {
			errs = append(errs, TaskFieldError{id, "remote", "missing required field"})
		}

		if t.Local == "" {
			errs = append(errs, TaskFieldError{id, "local", "missing required field"})
		} else if filepath.Base(t.Local) == "." {
			errs = append(errs, TaskFieldError{id, "local", fmt.Sprintf("invalid local path '%s'", t.Local)})
		}

		if !t.Mode.IsValid() {
			errs = append(errs, TaskFieldError{id, "mode", fmt.Sprintf("invalid mode '%s' (expected: copy, rename)", t.Mode)})
		}
	}

	if len(errs) > 0 {
		return &TaskValidationError{Errors: errs}
	}
	return nil
}
