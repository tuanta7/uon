package nginx

import (
	"context"
	"fmt"

	"github.com/tuanta7/uon/pkg/command"
)

// Run starts nginx immediately and enables it at boot.
func Run(ctx context.Context) error {
	if err := command.Run(ctx, "systemctl", "enable", "--now", "nginx.service"); err != nil {
		return fmt.Errorf("run nginx service: %w", err)
	}
	return nil
}
