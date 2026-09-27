package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/tuanta7/uon/pkg/command"
)

var removalPackages = append(append([]string{}, enginePackages...), "docker-ce-rootless-extras")

// Remove uninstalls Docker Engine and its APT repository. Docker data is only
// deleted when prune is true.
func Remove(ctx context.Context, prune bool) error {
	args := append([]string{"purge", "-y"}, removalPackages...)
	if err := command.RunQuietly(ctx, "apt-get", args...); err != nil {
		return fmt.Errorf("remove Docker Engine: %w", err)
	}
	for _, path := range []string{sourcesPath, keyringPath} {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s: %w", path, err)
		}
	}
	if !prune {
		return nil
	}
	for _, path := range []string{"/var/lib/docker", "/var/lib/containerd"} {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove %s: %w", path, err)
		}
	}
	return nil
}
