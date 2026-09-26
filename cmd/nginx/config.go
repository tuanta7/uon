package nginx

import (
	"fmt"

	"github.com/spf13/cobra"
)

func configCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show where to define and enable NGINX configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprint(cmd.OutOrStdout(), `Main configuration: /etc/nginx/nginx.conf
Site definitions:   /etc/nginx/sites-available/
Enabled sites:      /etc/nginx/sites-enabled/

To add a site:
  1. Create /etc/nginx/sites-available/<site>.conf
  2. Run: ln -s /etc/nginx/sites-available/<site>.conf /etc/nginx/sites-enabled/<site>.conf
  3. Run: nginx -t
  4. Run: systemctl reload nginx

Use "uon nginx load <file>" for an HTTP site configuration.
Use "uon nginx load --stream <file>" for a TCP/UDP stream configuration.
`)
			return err
		},
	}
}
