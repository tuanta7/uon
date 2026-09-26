package nginx

import (
	"fmt"

	"github.com/spf13/cobra"
	nginxservice "github.com/tuanta7/uon/internal/services/nginx"
)

func runCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Start NGINX and enable it at boot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			installed, err := nginxservice.IsInstalled(cmd.Context())
			if err != nil {
				return err
			}
			if !installed {
				return fmt.Errorf("nginx is not installed; run \"uon nginx install\" first")
			}
			if err := nginxservice.Run(cmd.Context()); err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "NGINX is running and enabled at boot.")
			return err
		},
	}
}
