package nginx

import (
	"context"
	"fmt"

	"github.com/tuanta7/uon/pkg/command"
)

// Remove removes the package while preserving its configuration files.
func Remove(ctx context.Context) error {
	if err := command.Run(ctx, "apt-get", "remove", "-y", "nginx"); err != nil {
		return fmt.Errorf("remove nginx: %w", err)
	}
	return nil
}
