package docker

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	dockerservice "github.com/tuanta7/uon/internal/services/docker"
)

func installCommand() *cobra.Command { return newInstallCommand(dockerservice.Install) }

func newInstallCommand(install func(context.Context) error) *cobra.Command {
	return &cobra.Command{
		Use: "install", Short: "Install Docker Engine from Docker's APT repository", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := install(cmd.Context()); err != nil {
				return err
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "Docker Engine is installed. Run \"uon docker run\" to start it and enable it at boot.")
			return err
		},
	}
}
