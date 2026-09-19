package network

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/tuanta7/uon/pkg/netplan"
)

type StaticIPConfig struct {
	Interface string
	Address   string // IP address with a CIDR prefix, for example 192.168.1.10/24.
	Gateway   string // Optional gateway with a CIDR prefix, for example 192.168.1.1/24.
}

func UseStaticIP(ctx context.Context, c StaticIPConfig) error {
	iface := strings.TrimSpace(c.Interface)
	if err := validateInterface(iface); err != nil {
		return err
	}

	address, err := validateAddress(c.Address)
	if err != nil {
		return fmt.Errorf("invalid static IP configuration: %w", err)
	}

	gateway := strings.TrimSpace(c.Gateway)
	addressIP, _, _ := net.ParseCIDR(address)
	if gateway != "" {
		gateway, err = validateAddress(gateway)
		if err != nil {
			return fmt.Errorf("invalid gateway: %w", err)
		}
		gatewayIP, _, _ := net.ParseCIDR(gateway)
		if (addressIP.To4() != nil) != (gatewayIP.To4() != nil) {
			return fmt.Errorf("gateway %q does not match address family", c.Gateway)
		}
		gateway = gatewayIP.String()
	}

	return netplan.Apply(ctx, iface, netplan.Interface{
		DHCP4:     false,
		DHCP6:     false,
		Addresses: []string{address},
		Routes:    netplan.DefaultRoutes(gateway),
	})
}
