package dgfilesystem

import (
	"fmt"
	"sync"

	"github.com/donnigundala/dg-core/contracts/filesystem"
)

var (
	globalDrivers   = make(map[string]DriverConstructor)
	globalDriversMu sync.RWMutex
)

// RegisterDriver registers a driver constructor globally.
func RegisterDriver(name string, constructor DriverConstructor) {
	globalDriversMu.Lock()
	defer globalDriversMu.Unlock()
	globalDrivers[name] = constructor
}

// Manager handles filesystem drivers.
type Manager struct {
	drivers map[string]DriverConstructor
	disks   map[string]filesystem.Disk
	mu      sync.RWMutex
}

// NewManager creates a new filesystem manager.
func NewManager() *Manager {
	m := &Manager{
		drivers: make(map[string]DriverConstructor),
		disks:   make(map[string]filesystem.Disk),
	}

	// Load globally registered drivers
	globalDriversMu.RLock()
	for name, constructor := range globalDrivers {
		m.drivers[name] = constructor
	}
	globalDriversMu.RUnlock()

	return m
}

// Extend registers a custom driver.
func (m *Manager) Extend(driverName string, constructor DriverConstructor) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers[driverName] = constructor
}

// Disk returns a disk instance by name.
func (m *Manager) Disk(name string, config map[string]interface{}) (filesystem.Disk, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return existing instance if available (singleton per disk name)
	if disk, ok := m.disks[name]; ok {
		return disk, nil
	}

	driverName, ok := config["driver"].(string)
	if !ok {
		return nil, fmt.Errorf("disk '%s' configuration missing 'driver'", name)
	}

	constructor, ok := m.drivers[driverName]
	if !ok {
		return nil, fmt.Errorf("driver '%s' not supported", driverName)
	}

	disk, err := constructor(config)
	if err != nil {
		return nil, err
	}

	// Wrap with observability decorator
	observedDisk := NewObservedDisk(disk, name)

	m.disks[name] = observedDisk
	return observedDisk, nil
}
