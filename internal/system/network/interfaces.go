package network

import (
	"context"
	"fmt"
	"net"
)

type Interface struct {
	Name      string
	Addresses []string
}

func GetAllInterfaces(ctx context.Context) ([]Interface, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list network interfaces: %w", err)
	}

	result := make([]Interface, 0, len(interfaces))
	for _, iface := range interfaces {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		addresses, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("list addresses for interface %q: %w", iface.Name, err)
		}
		entry := Interface{Name: iface.Name}
		for _, address := range addresses {
			entry.Addresses = append(entry.Addresses, address.String())
		}
		result = append(result, entry)
	}
	return result, nil
}
