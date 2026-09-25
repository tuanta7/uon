package nginx

import (
	"fmt"

	"github.com/spf13/cobra"
	nginxservice "github.com/tuanta7/uon/internal/services/nginx"
)

func newRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove NGINX while preserving its configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := nginxservice.Remove(cmd.Context()); err != nil {
				return err
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "NGINX is removed. Configuration files were preserved.")
			return err
		},
	}
}
