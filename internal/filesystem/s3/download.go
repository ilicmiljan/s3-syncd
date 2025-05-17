package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"time"
)

func GetLastModified(ctx context.Context, client *s3.Client, source BucketPath) (time.Time, error) {
	out, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &source.Bucket,
		Key:    &source.Key,
	})
	if err != nil {
		return time.Time{}, err
	}

	return *out.LastModified, nil
}

func DownloadS3Object(ctx context.Context, client *s3.Client, source BucketPath, destination *os.File) error {
	downloader := manager.NewDownloader(client)

	_, err := downloader.Download(ctx, destination, &s3.GetObjectInput{
		Bucket: aws.String(source.Bucket),
		Key:    aws.String(source.Key),
	})

	return err
}
