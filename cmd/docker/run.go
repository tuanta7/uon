package docker

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	dockerservice "github.com/tuanta7/uon/internal/services/docker"
)

func runCommand() *cobra.Command { return newRunCommand(dockerservice.Run) }

func newRunCommand(run func(context.Context) error) *cobra.Command {
	return &cobra.Command{
		Use: "run", Short: "Start Docker Engine and enable it at boot", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(cmd.Context()); err != nil {
				return err
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "Docker Engine is running and enabled at boot.")
			return err
		},
	}
}
