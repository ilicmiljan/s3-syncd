package config

import "github.com/aws/aws-sdk-go-v2/service/s3"

type Bucket struct {
	ID           string     `yaml:"-"`
	AccessKey    string     `yaml:"access_key"`
	SecretKey    string     `yaml:"secret_key"`
	Region       string     `yaml:"region"`
	Endpoint     string     `yaml:"endpoint"`
	UsePathStyle bool       `yaml:"use_path_style"`
	S3Client     *s3.Client `yaml:"-"`
}

type Buckets struct {
	Buckets map[string]Bucket `yaml:"buckets"`
}

type Task struct {
	ID       string  `yaml:"-"`
	BucketID string  `yaml:"bucket"`
	Cron     string  `yaml:"cron"`
	Remote   string  `yaml:"remote"`
	Local    string  `yaml:"local"`
	Bucket   *Bucket `yaml:"-"`
}

type Tasks struct {
	Tasks map[string]Task `yaml:"tasks"`
}
