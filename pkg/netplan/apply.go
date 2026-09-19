package netplan

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/tuanta7/uon/pkg/command"
	"go.yaml.in/yaml/v3"
)

func Apply(ctx context.Context, iface string, settings Interface) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	path := os.Getenv("UON_NETPLAN_FILE")
	if path == "" {
		path = "/etc/netplan/99-uon.yaml"
	}
	config := Config{Network: Network{
		Version:   2,
		Ethernets: map[string]Interface{},
	}}
	if existing, err := os.ReadFile(path); err == nil && len(existing) > 0 {
		if err := yaml.Unmarshal(existing, &config); err != nil {
			return fmt.Errorf("read existing Netplan configuration: %w", err)
		}
		if config.Network.Ethernets == nil {
			config.Network.Ethernets = map[string]Interface{}
		}
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read existing Netplan configuration: %w", err)
	}
	config.Network.Version = 2
	config.Network.Ethernets[iface] = settings
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal Netplan configuration: %w", err)
	}
	if err := writeAtomic(path, data); err != nil {
		return fmt.Errorf("write Netplan configuration: %w", err)
	}
	if err := command.Run(ctx, "netplan", "apply"); err != nil {
		return fmt.Errorf("apply Netplan configuration: %w", err)
	}
	return nil
}
