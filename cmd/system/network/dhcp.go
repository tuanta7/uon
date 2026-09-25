package network

import (
	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/internal/system/network"
)

func useDHCPCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dhcp <interface>",
		Short: "Configure an interface to use DHCP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network.UseDHCP(cmd.Context(), args[0])
		},
	}
}
