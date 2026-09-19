package network

import (
	"context"
	"strings"

	"github.com/tuanta7/uon/pkg/netplan"
)

func UseDHCP(ctx context.Context, iface string) error {
	iface = strings.TrimSpace(iface)
	if err := validateInterface(iface); err != nil {
		return err
	}
	return netplan.Apply(ctx, iface, netplan.Interface{
		DHCP4: true,
	})
}
