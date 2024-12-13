package model

// NetworkInterface represents a single network interface's details.
type NetworkInterface struct {
	ID            int64  // Primary key in the database
	SystemID      int64  // Foreign key referencing SystemInfo.ID
	InterfaceName string // e.g., "eth0", "wlan0"
	IPAddress     string // e.g., "192.168.1.10"
	MACAddress    string // e.g., "00:1A:2B:3C:4D:5E"
}
