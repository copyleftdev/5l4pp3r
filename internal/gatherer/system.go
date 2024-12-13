package gatherer

import (
	"os"
	"time"

	"github.com/copyleftdev/5l4pp3r/internal/model"
)

// GatherSystemInfo collects basic system-level information, such as hostname and current timestamp.
func GatherSystemInfo() (*model.SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &model.SystemInfo{
		Hostname:  hostname,
		CreatedAt: time.Now().UTC(),
	}, nil
}
