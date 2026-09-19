package network

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func validateAddress(addr string) (string, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", errors.New("IP address is required")
	}

	if _, _, err := net.ParseCIDR(addr); err != nil {
		return "", fmt.Errorf("invalid address %q: expected IP address with CIDR prefix: %w", addr, err)
	}
	return addr, nil
}

func validateInterface(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("network interface is required")
	}
	if _, err := net.InterfaceByName(name); err != nil {
		return fmt.Errorf("find network interface %q: %w", name, err)
	}
	return nil
}
