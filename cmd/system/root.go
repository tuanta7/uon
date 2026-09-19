package system

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/internal/system/network"
)

var Cmd = &cobra.Command{
	Use:   "system",
	Short: "Manage system settings",
	Long:  "Manage system settings. Without a subcommand, list all network interfaces and their IP addresses with CIDR prefixes.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		interfaces, err := network.GetAllInterfaces(cmd.Context())
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "INTERFACE\tIP ADDRESSES"); err != nil {
			return err
		}
		for _, iface := range interfaces {
			addresses := strings.Join(iface.Addresses, ", ")
			if addresses == "" {
				addresses = "-"
			}
			if _, err := fmt.Fprintf(w, "%s\t%s\n", iface.Name, addresses); err != nil {
				return err
			}
		}
		return w.Flush()
	},
}

func init() {
	Cmd.AddCommand(ipCmd)
	Cmd.AddCommand(newSleepCommand())
}
