package validation

import (
	"fmt"
	"strings"

	"go-sync-s3.miljanilic.com/internal/config"
)

func (e *BucketValidationError) Error() string {
	var lines []string
	for _, err := range e.Errors {
		lines = append(lines, fmt.Sprintf("bucket[%s]: %s (%s)", err.BucketID, err.Message, err.Field))
	}
	return "bucket validation failed:\n" + strings.Join(lines, "\n")
}

func ValidateBuckets(buckets map[string]config.Bucket) error {
	var errs []BucketFieldError

	for name, b := range buckets {
		if b.AccessKey == "" {
			errs = append(errs, BucketFieldError{name, "access_key", "missing required field"})
		}
		if b.SecretKey == "" {
			errs = append(errs, BucketFieldError{name, "secret_key", "missing required field"})
		}
		if b.Region == "" {
			errs = append(errs, BucketFieldError{name, "region", "missing required field"})
		}
		if b.Endpoint == "" {
			errs = append(errs, BucketFieldError{name, "endpoint", "missing required field"})
		}
	}

	if len(errs) > 0 {
		return &BucketValidationError{Errors: errs}
	}
	return nil
}
