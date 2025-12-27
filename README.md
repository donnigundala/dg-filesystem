# dg-filesystem

`dg-filesystem` is a powerful, abstracted file storage library for the DG Framework. It provides a simple, unified API for many filesystems (Local, S3) allowing you to swap storage drivers without changing your code.

## 📦 Installation

```bash
go get github.com/donnigundala/dg-filesystem
```

## 🚀 Usage

```go
package main

import (
    "github.com/donnigundala/dg-core/foundation"
    "github.com/donnigundala/dg-filesystem"
    _ "github.com/donnigundala/dg-filesystem/drivers/s3" // Register S3 driver
)

func main() {
    app := foundation.New(".")
    
    // Register provider (uses 'filesystem' key in config)
    app.Register(dgfilesystem.NewFilesystemServiceProvider(nil))
    
    app.Start()
    
    // Usage
    disk := dgfilesystem.MustResolve(app)
    disk.Put("hello.txt", []byte("Hello World"))
}
```

### Integration via InfrastructureSuite
In your `bootstrap/app.go`, you typically use the declarative suite pattern:

```go
import _ "github.com/donnigundala/dg-filesystem/drivers/s3"

func InfrastructureSuite(workerMode bool) []foundation.ServiceProvider {
	return []foundation.ServiceProvider{
		dgfilesystem.NewFilesystemServiceProvider(nil),
		// ... other providers
	}
}
```

> [!TIP]
> To use the S3 driver, you MUST include a blank import of `github.com/donnigundala/dg-filesystem/drivers/s3` in your bootstrap file to trigger its self-registration.

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

## Configuration

The plugin uses the `filesystem` key in your configuration file.

### Configuration Mapping (YAML vs ENV)

| YAML Key | Environment Variable | Default | Description |
| :--- | :--- | :--- | :--- |
| `filesystem.default` | `FILESYSTEM_DEFAULT` | `local` | Default disk name |
| `filesystem.disks.<name>.driver` | - | - | `local`, `s3` |
| `filesystem.disks.<name>.root` | - | - | Local root path |
| `filesystem.disks.<name>.bucket` | - | - | S3 bucket name |
| `filesystem.disks.<name>.region` | - | - | S3 region |
| `filesystem.disks.<name>.key` | - | - | AWS Access Key |
| `filesystem.disks.<name>.secret` | - | - | AWS Secret Key|

### Example YAML

```yaml
filesystem:
  default: local
  disks:
    local:
      driver: local
      root: "./storage/app"
    s3:
      driver: s3
      bucket: "my-app"
      region: "us-east-1"
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

## 📊 Observability

`dg-filesystem` is instrumented with OpenTelemetry metrics. If `dg-observability` is registered and enabled, the following metrics are automatically emitted for all disk operations:

*   `filesystem.disk.operation.count`: Counter (labels: `disk.name`, `operation`, `status`)
*   `filesystem.disk.operation.duration`: Histogram (labels: `disk.name`, `operation`, `status`) - measured in milliseconds.
*   `filesystem.disk.operation.bytes`: Counter (labels: `disk.name`, `operation`, `type`) - tracks bytes read or written.

To enable observability, ensure the `dg-observability` plugin is registered and configured:

```yaml
observability:
  enabled: true
  service_name: "my-app"
```
