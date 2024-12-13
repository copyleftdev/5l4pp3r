package model

import "time"

// ConfigFile represents a configuration file's metadata and compressed data.
type ConfigFile struct {
	ID           int64     // Primary key
	SystemID     int64     // Foreign key referencing SystemInfo.ID
	FilePath     string    // Full path to the file
	Size         int64     // Original file size in bytes
	Permissions  string    // File permission string
	ModifiedTime time.Time // Last modification time
	Data         []byte    // Compressed file data
}
