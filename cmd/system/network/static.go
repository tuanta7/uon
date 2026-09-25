package network

import (
	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/internal/system/network"
)

func useStaticCommand() *cobra.Command {
	var address, gateway string
	cmd := &cobra.Command{
		Use:   "static <interface>",
		Short: "Configure a static IP address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network.UseStaticIP(cmd.Context(), network.StaticIPConfig{
				Interface: args[0],
				Address:   address,
				Gateway:   gateway,
			})
		},
	}
	cmd.Flags().StringVar(&address, "address", "", "Static IP address with CIDR prefix (for example 192.168.1.10/24)")
	cmd.Flags().StringVar(&gateway, "gateway", "", "Optional gateway with CIDR prefix (for example 192.168.1.1/24)")
	_ = cmd.MarkFlagRequired("address")
	return cmd
}
