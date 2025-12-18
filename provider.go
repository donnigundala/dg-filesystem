package filesystem

import (
	"github.com/donnigundala/dg-core/container"
	"github.com/donnigundala/dg-core/contracts/foundation"
)

// FilesystemServiceProvider is the service provider for the filesystem module.
type FilesystemServiceProvider struct {
	// Config is auto-injected by dg-core if using config tags
	Config Config `config:"filesystem"`

	Manager *Manager
}

// NewFilesystemServiceProvider creates a new filesystem service provider.
func NewFilesystemServiceProvider() *FilesystemServiceProvider {
	return &FilesystemServiceProvider{}
}

// Name returns the provider name.
func (p *FilesystemServiceProvider) Name() string {
	return "filesystem"
}

// Version returns the provider version.
func (p *FilesystemServiceProvider) Version() string {
	return "1.0.0"
}

// Dependencies returns the provider dependencies.
func (p *FilesystemServiceProvider) Dependencies() []string {
	return []string{}
}

// Register registers the filesystem service.
func (p *FilesystemServiceProvider) Register(app foundation.Application) error {
	p.Manager = NewManager()

	// Register built-in drivers
	p.Manager.Extend("local", NewLocalDisk)

	app.Singleton("filesystem", func(c container.Container) (interface{}, error) {
		return p.Manager, nil
	})

	return nil
}

// Boot boots the filesystem service.
func (p *FilesystemServiceProvider) Boot(app foundation.Application) error {
	for name, diskConfig := range p.Config.Disks {
		// Create disk instance
		disk, err := p.Manager.Disk(name, diskConfig)
		if err != nil {
			return err
		}

		// Register disk in container as filesystem.name
		app.Singleton("filesystem."+name, func(c container.Container) (interface{}, error) {
			return disk, nil
		})

		// If this is the default disk, register it as "disk"
		if name == p.Config.Default {
			app.Singleton("disk", func(c container.Container) (interface{}, error) {
				return disk, nil
			})
		}
	}

	return nil
}
