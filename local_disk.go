package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalDisk implements the Disk interface for local storage.
type LocalDisk struct {
	root string
	url  string
}

// NewLocalDisk creates a new local disk instance.
func NewLocalDisk(config map[string]interface{}) (Disk, error) {
	root, ok := config["root"].(string)
	if !ok {
		return nil, fmt.Errorf("local driver requires 'root' config")
	}

	url, _ := config["url"].(string)

	// Ensure root directory exists
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("failed to create root directory: %w", err)
	}

	return &LocalDisk{
		root: root,
		url:  url,
	}, nil
}

func (d *LocalDisk) fullPath(path string) string {
	return filepath.Join(d.root, path)
}

func (d *LocalDisk) Put(path string, content []byte) error {
	fullPath := d.fullPath(path)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, content, 0644)
}

func (d *LocalDisk) PutStream(path string, content io.Reader) error {
	fullPath := d.fullPath(path)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, content)
	return err
}

func (d *LocalDisk) Get(path string) ([]byte, error) {
	return os.ReadFile(d.fullPath(path))
}

func (d *LocalDisk) GetStream(path string) (io.ReadCloser, error) {
	return os.Open(d.fullPath(path))
}

func (d *LocalDisk) Exists(path string) (bool, error) {
	_, err := os.Stat(d.fullPath(path))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (d *LocalDisk) Delete(path string) error {
	return os.Remove(d.fullPath(path))
}

func (d *LocalDisk) Url(path string) string {
	if d.url == "" {
		return path
	}
	// Clean path to ensure forward slashes
	cleanPath := filepath.ToSlash(path)
	// Remove leading slash if present
	cleanPath = strings.TrimPrefix(cleanPath, "/")

	return fmt.Sprintf("%s/%s", strings.TrimRight(d.url, "/"), cleanPath)
}

func (d *LocalDisk) SignedUrl(path string, expiration time.Duration) (string, error) {
	// Local driver does not support true signed URLs.
	// We return the standard URL as a fallback.
	return d.Url(path), nil
}

func (d *LocalDisk) MakeDirectory(path string) error {
	return os.MkdirAll(d.fullPath(path), 0755)
}

func (d *LocalDisk) DeleteDirectory(path string) error {
	return os.RemoveAll(d.fullPath(path))
}
