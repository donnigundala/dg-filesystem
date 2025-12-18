package filesystem

// Config holds the configuration for the filesystem service.
type Config struct {
	// Default is the name of the default disk.
	Default string `mapstructure:"default"`
	// Disks contains configuration for all disks.
	Disks map[string]map[string]interface{} `mapstructure:"disks"`
}
