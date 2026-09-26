package dev

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/pkg/command"
)

func upCommand() *cobra.Command {
	return newUpCommand(command.Run)
}

func newUpCommand(run func(context.Context, string, ...string) error) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Start development infrastructure",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(cmd.Context(), "docker", "compose", "--file", composeFile, "up", "--detach", "--wait"); err != nil {
				return fmt.Errorf("start development infrastructure: %w", err)
			}

			_, err := fmt.Fprint(cmd.OutOrStdout(), `Development infrastructure is ready.

Expose MinIO and Redis through the host NGINX service:
  uon nginx install
  uon nginx run
  uon nginx load --stream static/uon.dev.nginx.conf

Then use:
  MinIO API:     http://<server-LAN-IP>:9000
  MinIO console: http://<server-LAN-IP>:9001
  MinIO login:   minio / password
  Redis:         <server-LAN-IP>:6379
`)
			return err
		},
	}
}
