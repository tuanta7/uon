package dev

import "github.com/spf13/cobra"

const composeFile = "static/docker-compose.dev.yml"

// NewCommand creates the development infrastructure command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Manage local development infrastructure",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(upCommand())
	cmd.AddCommand(downCommand())
	return cmd
}
