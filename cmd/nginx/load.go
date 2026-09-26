package nginx

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	nginxservice "github.com/tuanta7/uon/internal/services/nginx"
)

func loadCommand() *cobra.Command {
	return newLoadCommand(nginxservice.Load)
}

func newLoadCommand(load func(context.Context, string, bool) error) *cobra.Command {
	var stream bool
	cmd := &cobra.Command{
		Use:   "load <file>",
		Short: "Load an NGINX configuration and reload NGINX",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := load(cmd.Context(), args[0], stream); err != nil {
				return err
			}
			configurationType := "site"
			if stream {
				configurationType = "stream"
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Loaded %s as an NGINX %s configuration and reloaded NGINX.\n", filepath.Base(filepath.Clean(args[0])), configurationType)
			return err
		},
	}
	cmd.Flags().BoolVar(&stream, "stream", false, "load a main-context TCP/UDP stream configuration")
	return cmd
}
