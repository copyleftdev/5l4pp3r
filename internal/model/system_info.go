package model

import "time"

// SystemInfo represents basic system-level information at snapshot time.
type SystemInfo struct {
	ID        int64     // Primary key in the database
	Hostname  string    // Hostname of the system
	CreatedAt time.Time // Timestamp of when the snapshot was taken
}
