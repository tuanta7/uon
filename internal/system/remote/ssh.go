package remote

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/tuanta7/uon/pkg/command"
)

// IsOpenSSHServerInstalled reports whether the openssh-server Debian package is
// fully installed.
func IsOpenSSHServerInstalled(ctx context.Context) (bool, error) {
	status, err := command.Output(ctx, "dpkg-query", "--show", "--showformat=${db:Status-Status}", "openssh-server")
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("check openssh-server installation: %w", err)
	}
	return status == "installed", nil
}

// InstallOpenSSHServer installs the Ubuntu OpenSSH server package.
func InstallOpenSSHServer(ctx context.Context) error {
	if err := command.Run(ctx, "apt-get", "install", "-y", "openssh-server"); err != nil {
		return fmt.Errorf("install openssh-server: %w", err)
	}
	return nil
}

// ToggleSSH controls whether the system accepts SSH connections.
// Enabling SSH starts the service immediately and configures it to start on boot.
// Disabling SSH stops the service and prevents it from starting on boot.
func ToggleSSH(ctx context.Context, allow bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	action := "disable"
	if allow {
		action = "enable"
	}
	if err := command.Run(ctx, "systemctl", action, "--now", "ssh.service"); err != nil {
		return fmt.Errorf("%s SSH service: %w", action, err)
	}
	return nil
}
