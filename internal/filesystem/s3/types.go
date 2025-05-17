package s3

import "fmt"

// ClientConfig holds configuration parameters for creating an S3 client
type ClientConfig struct {
	AccessKey    string
	SecretKey    string
	Region       string
	Endpoint     string
	UsePathStyle bool
}

type BucketPath struct {
	Bucket string
	Key    string
}

func (p BucketPath) String() string {
	return fmt.Sprintf("s3://%s/%s", p.Bucket, p.Key)
}
