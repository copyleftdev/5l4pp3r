package gatherer

import (
	"net"
	"strings"

	"github.com/copyleftdev/5l4pp3r/internal/model"
)

// GatherNetworkInfo enumerates network interfaces and captures their IP and MAC addresses.
func GatherNetworkInfo() ([]*model.NetworkInterface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var results []*model.NetworkInterface
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ipAddress string
		for _, addr := range addrs {
			ip := addr.String()
			if strings.Contains(ip, ".") {
				ipAddress = strings.Split(ip, "/")[0]
				break
			}
		}

		mac := iface.HardwareAddr.String()

		results = append(results, &model.NetworkInterface{
			InterfaceName: iface.Name,
			IPAddress:     ipAddress,
			MACAddress:    mac,
		})
	}

	return results, nil
}
