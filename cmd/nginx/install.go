package nginx

import (
	"fmt"

	"github.com/spf13/cobra"
	nginxservice "github.com/tuanta7/uon/internal/services/nginx"
)

func newInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install NGINX with APT",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := nginxservice.Install(cmd.Context()); err != nil {
				return err
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "NGINX is installed.\nRun \"uon nginx config\" for configuration locations and setup steps.")
			return err
		},
	}
}
