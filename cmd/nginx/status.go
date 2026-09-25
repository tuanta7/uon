package nginx

import (
	"fmt"

	"github.com/spf13/cobra"
	nginxservice "github.com/tuanta7/uon/internal/services/nginx"
)

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show NGINX package and service status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := nginxservice.GetStatus(cmd.Context())
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Installed: %s\nActive: %s\nEnabled: %s\n",
				yesNo(status.Installed), yesNo(status.Active), yesNo(status.Enabled))
			return err
		},
	}
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
