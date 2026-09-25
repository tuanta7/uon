package nginx

import (
	"context"
	"fmt"
)

// Status describes the package and systemd service state.
type Status struct {
	Installed bool
	Active    bool
	Enabled   bool
}

// GetStatus returns the nginx package, runtime, and boot-time states.
func GetStatus(ctx context.Context) (Status, error) {
	installed, err := IsInstalled(ctx)
	if err != nil {
		return Status{}, err
	}
	status := Status{Installed: installed}
	if !installed {
		return status, nil
	}

	status.Active, err = commandSucceedsOrHasStateExit(ctx, "systemctl", 3, "is-active", "nginx.service")
	if err != nil {
		return Status{}, fmt.Errorf("check whether nginx is active: %w", err)
	}
	status.Enabled, err = commandSucceedsOrHasStateExit(ctx, "systemctl", 1, "is-enabled", "nginx.service")
	if err != nil {
		return Status{}, fmt.Errorf("check whether nginx is enabled: %w", err)
	}
	return status, nil
}
