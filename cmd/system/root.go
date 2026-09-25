package system

import (
	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/cmd/system/network"
	"github.com/tuanta7/uon/cmd/system/sleep"
	"github.com/tuanta7/uon/cmd/system/ssh"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "Manage system settings",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(network.Command())
	cmd.AddCommand(sleep.Command())
	cmd.AddCommand(ssh.Command())
	return cmd
}
