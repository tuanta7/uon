package docker

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	dockerservice "github.com/tuanta7/uon/internal/services/docker"
)

func removeCommand() *cobra.Command { return newRemoveCommand(dockerservice.Remove) }

func newRemoveCommand(remove func(context.Context, bool) error) *cobra.Command {
	var prune bool
	cmd := &cobra.Command{
		Use: "remove", Short: "Remove Docker Engine while preserving its data", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := remove(cmd.Context(), prune); err != nil {
				return err
			}
			message := "Docker Engine is removed. Images, containers, volumes, and configuration were preserved."
			if prune {
				message = "Docker Engine and its stored images, containers, and volumes were removed."
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), message)
			return err
		},
	}
	cmd.Flags().BoolVar(&prune, "prune", false, "Remove images, containers, and volumes in Docker's data directories")
	return cmd
}
