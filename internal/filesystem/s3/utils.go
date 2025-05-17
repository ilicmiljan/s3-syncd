package s3

import (
	"fmt"
	"strings"
)

func ParseS3Path(path string) (BucketPath, error) {
	trimmed := strings.TrimPrefix(path, "s3a://")
	trimmed = strings.TrimPrefix(trimmed, "s3://")

	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return BucketPath{}, fmt.Errorf("invalid S3 path: %s", path)
	}

	return BucketPath{
		Bucket: parts[0],
		Key:    parts[1],
	}, nil
}
