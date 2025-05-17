package local

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"time"
)

func CreateTempFile(dir string) (*TempFile, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	path := filepath.Join(dir, fmt.Sprintf("%s.part", uuid.New().String()))

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file %s: %w", path, err)
	}

	cleanup := func() {
		_ = file.Close()

		if err := os.Remove(path); err == nil {
			log.Printf("cleaned up temp file: %s", path)
		}
	}

	return &TempFile{
		Path:    path,
		File:    file,
		Cleanup: cleanup,
	}, nil
}

func GetLastModified(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func Replace(source, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	return os.Rename(source, destination)
}
