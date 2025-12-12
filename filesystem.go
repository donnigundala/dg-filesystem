package filesystem

import (
	"io"
	"time"
)

// Disk represents a storage disk (local or cloud).
type Disk interface {
	// Put writes the given content to a file at the given path.
	Put(path string, content []byte) error

	// PutStream writes the given stream to a file at the given path.
	PutStream(path string, content io.Reader) error

	// Get retrieves the content of a file at the given path.
	Get(path string) ([]byte, error)

	// GetStream retrieves the stream of a file at the given path.
	GetStream(path string) (io.ReadCloser, error)

	// Exists checks if a file exists at the given path.
	Exists(path string) (bool, error)

	// Delete removes a file at the given path.
	Delete(path string) error

	// Url returns the public URL for the file at the given path.
	Url(path string) string

	// SignedUrl returns a temporary URL for the file at the given path, valid for the given duration.
	// Useful for giving temporary access to private files.
	SignedUrl(path string, expiration time.Duration) (string, error)

	// MakeDirectory creates a directory at the given path.
	MakeDirectory(path string) error

	// DeleteDirectory removes a directory at the given path.
	DeleteDirectory(path string) error
}

// DriverConstructor is a function that creates a new Disk instance.
type DriverConstructor func(config map[string]interface{}) (Disk, error)
