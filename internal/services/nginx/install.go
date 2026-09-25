package nginx

import (
	"context"
	"fmt"

	"github.com/tuanta7/uon/pkg/command"
)

// Install installs the Ubuntu nginx package.
func Install(ctx context.Context) error {
	if err := command.Run(ctx, "apt-get", "install", "-y", "nginx"); err != nil {
		return fmt.Errorf("install nginx: %w", err)
	}
	return nil
}

// IsInstalled reports whether the nginx Debian package is fully installed.
func IsInstalled(ctx context.Context) (bool, error) {
	status, err := command.Output(ctx, "dpkg-query", "--show", "--showformat=${db:Status-Status}", "nginx")
	if hasExitCode(err, 1) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("check nginx installation: %w", err)
	}

	return status == "installed", nil
}
