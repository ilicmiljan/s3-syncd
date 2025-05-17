package worker

import (
	"context"
	"fmt"
	"go-sync-s3.miljanilic.com/internal/config"
	"go-sync-s3.miljanilic.com/internal/filesystem/local"
	"go-sync-s3.miljanilic.com/internal/filesystem/s3"
	"log"
)

func Execute(ctx context.Context, t config.Task) error {
	logPrefix := fmt.Sprintf("task[%s]:", t.ID)

	if ctx.Err() != nil {
		log.Printf("%s context cancelled before execution", logPrefix)
		return ctx.Err()
	}

	log.Printf("%s starting execution...", logPrefix)

	bucket := t.Bucket
	if bucket == nil {
		return fmt.Errorf("%s no bucket assigned", logPrefix)
	}

	bucket.WithS3Client()

	tempFile, err := local.CreateTempFile("tmp/")
	if err != nil {
		return fmt.Errorf("%s failed to create temp file: %w", logPrefix, err)
	}
	defer tempFile.Cleanup()

	source, err := s3.ParseS3Path(t.Remote)
	if err != nil {
		return fmt.Errorf("%s invalid remote path: %w", logPrefix, err)
	}
	log.Printf("%s resolved S3 path: %s", logPrefix, source)

	if !ShouldDownload(ctx, t, bucket, source, logPrefix) {
		log.Printf("%s skipping download: remote file %s has not been modified since local file %s", logPrefix, t.Remote, t.Local)
		return nil
	}

	log.Printf("%s downloading from S3...", logPrefix)
	err = s3.DownloadS3Object(ctx, bucket.S3Client, source, tempFile.File)
	if err != nil {
		return fmt.Errorf("%s download failed: %w", logPrefix, err)
	}
	log.Printf("%s download complete", logPrefix)

	if err := tempFile.File.Close(); err != nil {
		return fmt.Errorf("%s failed to close file: %w", logPrefix, err)
	}

	switch t.Mode {
	case config.ModeCopy:
		log.Printf("%s using copy mode", logPrefix)
		if err := local.CopyAndDelete(tempFile.Path, t.Local); err != nil {
			return fmt.Errorf("%s failed to copy and delete: %w", logPrefix, err)
		}
	case config.ModeRename, "":
		log.Printf("%s using rename mode", logPrefix)
		if err := local.Rename(tempFile.Path, t.Local); err != nil {
			return fmt.Errorf("%s failed to rename temp file: %w", logPrefix, err)
		}
	default:
		return fmt.Errorf("%s unknown mode: %s", logPrefix, t.Mode)
	}

	log.Printf("%s moved %s to %s", logPrefix, tempFile.Path, t.Local)
	return nil
}

func ShouldDownload(ctx context.Context, t config.Task, bucket *config.Bucket, source s3.BucketPath, logPrefix string) bool {
	lastModifiedSource, err := s3.GetLastModified(ctx, bucket.S3Client, source)
	if err != nil {
		return true
	}

	log.Printf("%s file %s last modified: %s", logPrefix, t.Remote, lastModifiedSource)

	lastModifiedDestination, err := local.GetLastModified(t.Local)
	if err != nil {
		return true
	}

	log.Printf("%s file %s last modified: %s", logPrefix, t.Local, lastModifiedDestination)

	return lastModifiedSource.After(lastModifiedDestination)
}
