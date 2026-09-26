package nginx

import (
	"fmt"

	"github.com/spf13/cobra"
	nginxservice "github.com/tuanta7/uon/internal/services/nginx"
)

func removeCommand() *cobra.Command {
	var prune bool
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove NGINX while preserving its configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := nginxservice.Remove(cmd.Context(), prune); err != nil {
				return err
			}
			if prune {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), "NGINX is removed. Configuration files were pruned.")
				return err
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "NGINX is removed. Configuration files were preserved.")
			return err
		},
	}
	cmd.Flags().BoolVar(&prune, "prune", false, "Remove NGINX configuration files as well")
	return cmd
}
