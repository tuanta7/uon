package nginx

import (
	"context"
	"fmt"

	"github.com/tuanta7/uon/pkg/command"
)

func Remove(ctx context.Context, prune bool) error {
	args := []string{"remove", "-y", "nginx", "libnginx-mod-stream"}
	if prune {
		// Ubuntu keeps the default site and most shared configuration
		// in the nginx-common package
		args = []string{"purge", "-y", "nginx", "nginx-common", "libnginx-mod-stream"}
	}
	if err := command.Run(ctx, "apt-get", args...); err != nil {
		return fmt.Errorf("remove nginx: %w", err)
	}
	return nil
}
