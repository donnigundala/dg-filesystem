package dgfilesystem

import (
	"github.com/donnigundala/dg-core/contracts/filesystem"
)

// DriverConstructor is a function that creates a new Disk instance.
type DriverConstructor func(config map[string]interface{}) (filesystem.Disk, error)
