package docker

import (
	"context"
	"fmt"

	"github.com/tuanta7/uon/pkg/command"
)

// Run starts Docker Engine immediately and enables it at boot.
func Run(ctx context.Context) error {
	if err := command.RunQuietly(ctx, "systemctl", "enable", "--now", "docker.service"); err != nil {
		return fmt.Errorf("run Docker service: %w", err)
	}
	return nil
}
