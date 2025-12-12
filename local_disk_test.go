package filesystem

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocalDisk(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	config := map[string]interface{}{
		"root": tmpDir,
		"url":  "http://localhost/storage",
	}

	disk, err := NewLocalDisk(config)
	assert.NoError(t, err)

	// Test Put & Get
	err = disk.Put("test.txt", []byte("hello world"))
	assert.NoError(t, err)

	content, err := disk.Get("test.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(content))

	// Test Exists
	exists, err := disk.Exists("test.txt")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = disk.Exists("nonexistent.txt")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test Url
	url := disk.Url("test.txt")
	assert.Equal(t, "http://localhost/storage/test.txt", url)

	// Test SignedUrl (Local implementation returns standard URL)
	signedUrl, err := disk.SignedUrl("test.txt", 1*time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost/storage/test.txt", signedUrl)

	// Test Delete
	err = disk.Delete("test.txt")
	assert.NoError(t, err)

	exists, err = disk.Exists("test.txt")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test Directory
	err = disk.MakeDirectory("nested/dir")
	assert.NoError(t, err)

	err = disk.Put("nested/dir/file.txt", []byte("data"))
	assert.NoError(t, err)

	exists, err = disk.Exists("nested/dir/file.txt")
	assert.NoError(t, err)
	assert.True(t, exists)
}
