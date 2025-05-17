package local

import (
	"fmt"
	"github.com/google/uuid"
	"io"
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

func Rename(source, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	return os.Rename(source, destination)
}

func CopyAndDelete(source, destination string) error {
	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", filepath.Dir(destination), err)
	}

	// Open the source file
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		_ = destinationFile.Close()
	}()

	// Copy data
	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Flush to disk
	if err := destinationFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	// Remove the source file
	if err := os.Remove(source); err != nil {
		return fmt.Errorf("failed to delete source file: %w", err)
	}

	return nil
}
