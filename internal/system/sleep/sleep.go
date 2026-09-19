package sleep

import (
	"context"
	"fmt"

	"github.com/tuanta7/uon/pkg/command"
)

// ToggleSleep allows or prevents system-wide sleep, including desktop sessions.
// When allowSleep is true, it unmasks systemd's sleep targets to allow suspend
// and hibernation. When false, it masks them to prevent both operations.
// Changes persist across reboots and require root privileges.
// Screen blanking and locking settings are unchanged.
func ToggleSleep(ctx context.Context, allowSleep bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	action := "mask"
	if allowSleep {
		action = "unmask"
	}
	if err := command.Run(ctx, "systemctl", action,
		"sleep.target",
		"suspend.target",
		"hibernate.target",
		"hybrid-sleep.target",
		"suspend-then-hibernate.target",
	); err != nil {
		return fmt.Errorf("%s system sleep targets: %w", action, err)
	}
	return nil
}
