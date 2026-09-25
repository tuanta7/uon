package network

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage network interfaces",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(listCommand())
	cmd.AddCommand(useStaticCommand())
	cmd.AddCommand(useDHCPCommand())
	return cmd
}
