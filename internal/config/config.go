package config

import (
	"go-sync-s3.miljanilic.com/internal/filesystem/s3"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadBuckets(path string) (map[string]Bucket, error) {
	var cfg Buckets

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)

	for id, task := range cfg.Buckets {
		task.ID = id
		cfg.Buckets[id] = task
	}

	return cfg.Buckets, err
}

func (b *Bucket) WithS3Client() {
	if b.S3Client != nil {
		return
	}

	s3cfg := s3.ClientConfig{
		AccessKey:    b.AccessKey,
		SecretKey:    b.SecretKey,
		Region:       b.Region,
		Endpoint:     b.Endpoint,
		UsePathStyle: b.UsePathStyle,
	}

	b.S3Client = s3.CreateClient(s3cfg)
}

func LoadTasks(path string) (map[string]Task, error) {
	var cfg Tasks

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)

	for id, task := range cfg.Tasks {
		task.ID = id
		cfg.Tasks[id] = task
	}

	return cfg.Tasks, err
}
