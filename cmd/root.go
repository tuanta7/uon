package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/cmd/nginx"
	"github.com/tuanta7/uon/cmd/system"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uon",
		Short: "Set up an Ubuntu server for self-hosting",
		Long: `UON helps you set up an Ubuntu server on a laptop or PC for a lightweight
self-hosted environment. Manage system settings, configure network interfaces,
and install tools for your server.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if os.Getuid() != 0 {
				return fmt.Errorf("must be run as root, current: %d", os.Getuid())
			}
			return nil
		},
	}
	cmd.AddCommand(nginx.NewCommand())
	cmd.AddCommand(system.NewCommand())
	return cmd
}

func Execute() {
	err := NewCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
