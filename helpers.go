package dgfilesystem

import (
	"fmt"

	"github.com/donnigundala/dg-core/container"
	"github.com/donnigundala/dg-core/contracts/filesystem"
)

// Resolve returns the default disk from the container.
func Resolve(c container.Container) (filesystem.Disk, error) {
	instance, err := c.Make(Binding)
	if err != nil {
		return nil, err
	}
	disk, ok := instance.(filesystem.Disk)
	if !ok {
		return nil, fmt.Errorf("resolved object is not a Disk")
	}
	return disk, nil
}

// MustResolve returns the default disk or panics.
func MustResolve(c container.Container) filesystem.Disk {
	disk, err := Resolve(c)
	if err != nil {
		panic(err)
	}
	return disk
}
