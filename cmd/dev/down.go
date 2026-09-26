package dev

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/pkg/command"
)

func downCommand() *cobra.Command {
	return newDownCommand(command.Run)
}

func newDownCommand(run func(context.Context, string, ...string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop development infrastructure and delete its data",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(cmd.Context(), "docker", "compose", "--file", composeFile, "down", "--volumes", "--remove-orphans"); err != nil {
				return fmt.Errorf("stop development infrastructure: %w", err)
			}

			_, err := fmt.Fprintln(cmd.OutOrStdout(), "Development infrastructure and data were removed.")
			return err
		},
	}
}
