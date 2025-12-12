package filesystem

import (
	"fmt"

	"github.com/donnigundala/dg-core/container"
)

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
