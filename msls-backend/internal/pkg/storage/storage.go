// Package storage provides file storage abstractions for the application.
package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// Storage defines the interface for file storage operations.
type Storage interface {
	// Upload stores a file at the given path.
	Upload(ctx context.Context, path string, data []byte, contentType string) error

	// Download retrieves a file from the given path.
	Download(ctx context.Context, path string) ([]byte, error)

	// Delete removes a file at the given path.
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists at the given path.
	Exists(ctx context.Context, path string) (bool, error)

	// GetPresignedURL generates a presigned URL for downloading a file.
	// For local storage, this returns a direct file path or URL.
	GetPresignedURL(ctx context.Context, path string, fileName string) (string, error)
}

// LocalStorage implements Storage interface using the local filesystem.
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new local storage instance.
func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	// Ensure base directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

// Upload stores a file at the given path.
func (s *LocalStorage) Upload(ctx context.Context, path string, data []byte, contentType string) error {
	fullPath := filepath.Join(s.basePath, path)

	// Ensure parent directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// Download retrieves a file from the given path.
func (s *LocalStorage) Download(ctx context.Context, path string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("read file: %w", err)
	}

	return data, nil
}

// Delete removes a file at the given path.
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // File already doesn't exist
		}
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}

// Exists checks if a file exists at the given path.
func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.basePath, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat file: %w", err)
	}

	return true, nil
}

// GetPresignedURL returns a URL for downloading the file.
// For local storage, this returns a public URL path.
func (s *LocalStorage) GetPresignedURL(ctx context.Context, path string, fileName string) (string, error) {
	// For local storage, return a URL that can be served by the static file handler
	downloadURL := fmt.Sprintf("%s/%s", s.baseURL, path)

	// URL encode the path
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return "", fmt.Errorf("parse URL: %w", err)
	}

	// Add filename as a query parameter for download hint
	q := parsedURL.Query()
	q.Set("filename", fileName)
	q.Set("expires", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))
	parsedURL.RawQuery = q.Encode()

	return parsedURL.String(), nil
}

// GetReader returns an io.ReadCloser for streaming large files.
func (s *LocalStorage) GetReader(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("open file: %w", err)
	}

	return file, nil
}
