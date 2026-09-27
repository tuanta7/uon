package docker

import "github.com/spf13/cobra"

// NewCommand creates the Docker Engine command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "docker", Short: "Manage Docker Engine", Args: cobra.NoArgs}
	cmd.AddCommand(installCommand(), runCommand(), removeCommand())
	return cmd
}
