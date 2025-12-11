package filesystem

import (
	"fmt"
	"sync"
)

// Manager handles filesystem drivers.
type Manager struct {
	drivers map[string]DriverConstructor
	disks   map[string]Disk
	mu      sync.RWMutex
}

// NewManager creates a new filesystem manager.
func NewManager() *Manager {
	return &Manager{
		drivers: make(map[string]DriverConstructor),
		disks:   make(map[string]Disk),
	}
}

// Extend registers a custom driver.
func (m *Manager) Extend(driverName string, constructor DriverConstructor) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers[driverName] = constructor
}

// Disk returns a disk instance by name.
func (m *Manager) Disk(name string, config map[string]interface{}) (Disk, error) {
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

	m.disks[name] = disk
	return disk, nil
}
