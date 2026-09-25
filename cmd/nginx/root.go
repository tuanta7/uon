package nginx

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nginx",
		Short: "Manage the NGINX web server",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newInstallCommand())
	cmd.AddCommand(newRunCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newRemoveCommand())
	cmd.AddCommand(newConfigCommand())
	return cmd
}
