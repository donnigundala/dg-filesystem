package filesystem

import (
	"github.com/donnigundala/dg-core/container"
	"github.com/donnigundala/dg-core/contracts/foundation"
	"github.com/spf13/viper"
)

// FilesystemServiceProvider is the service provider for the filesystem module.
type FilesystemServiceProvider struct {
	Manager *Manager
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
	config := viper.GetStringMap("filesystem.disks")

	for name, diskConfig := range config {
		cfg, ok := diskConfig.(map[string]interface{})
		if !ok {
			continue
		}

		// Create disk instance
		disk, err := p.Manager.Disk(name, cfg)
		if err != nil {
			return err
		}

		// Register disk in container as filesystem.name
		app.Singleton("filesystem."+name, func(c container.Container) (interface{}, error) {
			return disk, nil
		})

		// If this is the default disk, register it as "disk"
		if name == viper.GetString("filesystem.default") {
			app.Singleton("disk", func(c container.Container) (interface{}, error) {
				return disk, nil
			})
		}
	}

	return nil
}
