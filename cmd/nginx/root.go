package nginx

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nginx",
		Short: "Manage the NGINX web server",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(installCommand())
	cmd.AddCommand(runCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(removeCommand())
	cmd.AddCommand(configCommand())
	cmd.AddCommand(loadCommand())
	return cmd
}
