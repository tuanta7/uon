package system

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/internal/system/network"
)

var ipCmd = newIPCommand()

func newIPCommand() *cobra.Command {
	var iface, address, gateway string
	var dhcp bool
	cmd := &cobra.Command{
		Use:   "ip",
		Short: "Configure a network interface with a static IP or DHCP",
		Example: `  uon system ip --interface eth0 --static 192.168.1.10/24 --gateway 192.168.1.1/24
  uon system ip --interface eth0 --dhcp`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(iface) == "" {
				return fmt.Errorf("--interface is required")
			}
			if dhcp {
				return network.UseDHCP(cmd.Context(), iface)
			}
			if !cmd.Flags().Changed("static") {
				return fmt.Errorf("provide --static or --dhcp")
			}
			return network.UseStaticIP(cmd.Context(), network.StaticIPConfig{
				Interface: iface,
				Address:   address,
				Gateway:   gateway,
			})
		},
	}
	cmd.Flags().StringVarP(&iface, "interface", "i", "", "Network interface to configure")
	cmd.Flags().StringVar(&address, "static", "", "Static IP address with CIDR prefix (for example 192.168.1.10/24)")
	cmd.Flags().StringVar(&gateway, "gateway", "", "Optional gateway with CIDR prefix (for example 192.168.1.1/24); requires --static")
	cmd.Flags().BoolVar(&dhcp, "dhcp", false, "Use DHCP")
	cmd.MarkFlagsMutuallyExclusive("static", "dhcp")
	cmd.MarkFlagsMutuallyExclusive("gateway", "dhcp")
	return cmd
}
