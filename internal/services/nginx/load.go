package nginx

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tuanta7/uon/pkg/command"
)

const (
	sitesAvailableDirectory = "/etc/nginx/sites-available"
	sitesEnabledDirectory   = "/etc/nginx/sites-enabled"
	modulesEnabledDirectory = "/etc/nginx/modules-enabled"
)

// Load installs an HTTP site configuration, or a main-context configuration
// when stream is true. It then validates and reloads NGINX. Existing
// configurations with the same file name are replaced.
func Load(ctx context.Context, sourcePath string, stream bool) error {
	fileName := filepath.Base(filepath.Clean(sourcePath))
	if fileName == "." || fileName == string(filepath.Separator) {
		return fmt.Errorf("load nginx configuration: invalid file path %q", sourcePath)
	}

	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open nginx configuration %q: %w", sourcePath, err)
	}
	defer source.Close()

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("inspect nginx configuration %q: %w", sourcePath, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("load nginx configuration: %q is not a regular file", sourcePath)
	}

	availablePath := filepath.Join(sitesAvailableDirectory, fileName)
	enabledPath := filepath.Join(sitesEnabledDirectory, fileName)
	if stream {
		availablePath = filepath.Join(modulesEnabledDirectory, "90-"+fileName)
		enabledPath = ""
	}
	if err := replaceFile(source, availablePath); err != nil {
		return fmt.Errorf("copy nginx configuration to %q: %w", availablePath, err)
	}

	if enabledPath != "" {
		if err := replaceSymlink(availablePath, enabledPath); err != nil {
			return fmt.Errorf("enable nginx configuration at %q: %w", enabledPath, err)
		}
	}

	if err := command.Run(ctx, "nginx", "-t"); err != nil {
		return fmt.Errorf("validate nginx configuration: %w", err)
	}
	if err := command.Run(ctx, "systemctl", "reload", "nginx.service"); err != nil {
		return fmt.Errorf("reload nginx service: %w", err)
	}
	return nil
}

func replaceFile(source *os.File, destination string) (returnErr error) {
	directory := filepath.Dir(destination)
	temporary, err := os.CreateTemp(directory, "."+filepath.Base(destination)+".*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() {
		if returnErr != nil {
			_ = os.Remove(temporaryPath)
		}
	}()

	if _, err := io.Copy(temporary, source); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, destination)
}

func replaceSymlink(target, destination string) error {
	directory := filepath.Dir(destination)
	temporary, err := os.CreateTemp(directory, "."+filepath.Base(destination)+".*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	if err := temporary.Close(); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	if err := os.Remove(temporaryPath); err != nil {
		return err
	}
	if err := os.Symlink(target, temporaryPath); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, destination); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	return nil
}
