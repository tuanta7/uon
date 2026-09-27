package command

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func RunQuietly(ctx context.Context, name string, args ...string) error {
	_, err := Run(ctx, name, args...)
	return err
}

func Run(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		message := strings.TrimSpace(string(output))
		if message != "" {
			return "", fmt.Errorf("%w: %s", err, message)
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
