package filesystem

import (
	"fmt"

	"github.com/donnigundala/dg-core/container"
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
	return []string{"config"}
}

// Register registers the filesystem service.
func (p *FilesystemServiceProvider) Register(c container.Container) error {
	p.Manager = NewManager()

	// Register built-in drivers
	p.Manager.Extend("local", NewLocalDisk)

	c.Singleton("filesystem", func(c container.Container) (interface{}, error) {
		return p.Manager, nil
	})

	return nil
}

// Boot boots the filesystem service.
func (p *FilesystemServiceProvider) Boot(c container.Container) error {
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
		c.Singleton("filesystem."+name, func(c container.Container) (interface{}, error) {
			return disk, nil
		})

		// If this is the default disk, register it as "disk"
		if name == viper.GetString("filesystem.default") {
			c.Singleton("disk", func(c container.Container) (interface{}, error) {
				return disk, nil
			})
		}
	}

	return nil
}

// Injectable allows automatic dependency injection of the default Disk.
type Injectable struct {
	Disk Disk
}

// Provide implements the Injectable interface for DI.
func (i *Injectable) Provide(c container.Container) error {
	var err error
	i.Disk, err = Resolve(c)
	return err
}

// Resolve returns the default disk from the container.
func Resolve(c container.Container) (Disk, error) {
	instance, err := c.Make("disk")
	if err != nil {
		return nil, err
	}
	disk, ok := instance.(Disk)
	if !ok {
		return nil, fmt.Errorf("resolved object is not a Disk")
	}
	return disk, nil
}

// MustResolve returns the default disk or panics.
func MustResolve(c container.Container) Disk {
	disk, err := Resolve(c)
	if err != nil {
		panic(err)
	}
	return disk
}
