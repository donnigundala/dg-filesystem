# dg-filesystem

`dg-filesystem` is a powerful, abstracted file storage library for the DG Framework. It provides a simple, unified API for many filesystems (Local, S3) allowing you to swap storage drivers without changing your code.

## 📦 Installation

```bash
go get github.com/donnigundala/dg-filesystem
```

## 🚀 Usage

### 1. Configuration

Add the filesystem configuration to your `config/filesystem.yaml` (or via Viper):

```yaml
default: "local"

disks:
  local:
    driver: "local"
    root: "./storage/app"
    url: "http://localhost:8080/storage"

  s3:
    driver: "s3"
    region: "us-east-1"
    bucket: "my-bucket"
    url: "https://my-bucket.s3.amazonaws.com"
    # Env vars AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are used automatically
```

### 2. Registration

If using `skeleton`, register the provider in `bootstrap/providers.go`:

```go
import (
    "github.com/donnigundala/dg-filesystem"
    "github.com/donnigundala/dg-filesystem/drivers/s3"
)

// ...

func Providers() []interface{} {
    return []interface{}{
        // ...
        &filesystem.FilesystemServiceProvider{},
    }
}
```

To enable the S3 driver (it is not enabled by default to keep dependencies light), you must register it manually or in a custom provider:

```go
// In AppServiceProvider.Boot or similar
fsManager := app.Make("filesystem").(*filesystem.Manager)
fsManager.Extend("s3", s3.NewS3Disk)
```

**Wait**, my current `FilesystemServiceProvider` creates a NEW manager every time `Register` is called? No, it's a singleton.
But `FilesystemServiceProvider.Register` calls `p.Manager = NewManager()`.
If the user wants to extending it, they need access to the manager.
They can get it via `app.Make("filesystem")`.

### 3. Basic Usage

```go
package main

import (
    "github.com/donnigundala/dg-filesystem"
)

type UserController struct {
    // Inject default disk
    Disk filesystem.Disk
}

func (c *UserController) UploadProfile(ctx *gin.Context) {
    file, _ := ctx.FormFile("avatar")
    f, _ := file.Open()
    defer f.Close()

    // Store file
    err := c.Disk.PutStream("avatars/"+file.Filename, f)
    
    // Get URL
    url := c.Disk.Url("avatars/"+file.Filename)
    
    // Get Temporary URL (Signed) - Valid for 15 minutes
    // Great for private S3 files
    signedUrl, _ := c.Disk.SignedUrl("avatars/"+file.Filename, 15*time.Minute)

    ctx.JSON(200, gin.H{
        "url": url,
        "signed_url": signedUrl,
    })
}
```

### 4. Direct Access

```go
// Get default disk
disk := filesystem.MustResolve(container)

// Get specific disk
manager := container.MustMake("filesystem").(*filesystem.Manager)
s3Disk, _ := manager.Disk("s3", s3Config)
```

## 🔌 Drivers

### Local
Supported out of the box.

### S3
Requires `github.com/donnigundala/dg-filesystem/drivers/s3`.

```go
import "github.com/donnigundala/dg-filesystem/drivers/s3"

manager.Extend("s3", s3.NewS3Disk)
```

## 📝 Interface

```go
type Disk interface {
    Put(path string, content []byte) error
    PutStream(path string, content io.Reader) error
    Get(path string) ([]byte, error)
    GetStream(path string) (io.ReadCloser, error)
    Exists(path string) (bool, error)
    Delete(path string) error
    Url(path string) string
    SignedUrl(path string, expiration time.Duration) (string, error)
    MakeDirectory(path string) error
    DeleteDirectory(path string) error
}
```
