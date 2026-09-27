package docker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tuanta7/uon/pkg/command"
	"github.com/tuanta7/uon/pkg/ubuntu"
)

const (
	dockerKeyURL = "https://download.docker.com/linux/ubuntu/gpg"
	dockerAPTURL = "https://download.docker.com/linux/ubuntu"
	keyringPath  = "/etc/apt/keyrings/docker.asc"
	sourcesPath  = "/etc/apt/sources.list.d/docker.sources"
)

var conflictingPackages = []string{
	"docker.io", "docker-compose", "docker-compose-v2", "docker-doc",
	"docker-buildx", "podman-docker", "containerd", "runc",
}

var enginePackages = []string{
	"docker-ce", "docker-ce-cli", "containerd.io",
	"docker-buildx-plugin", "docker-compose-plugin",
}

// Install installs Docker Engine from Docker's official Ubuntu APT repository.
func Install(ctx context.Context) error {
	if err := removeConflictingPackages(ctx); err != nil {
		return err
	}
	if err := command.RunQuietly(ctx, "apt-get", "update"); err != nil {
		return fmt.Errorf("update apt package index: %w", err)
	}
	if err := command.RunQuietly(ctx, "apt-get", "install", "-y", "ca-certificates", "curl"); err != nil {
		return fmt.Errorf("install Docker repository prerequisites: %w", err)
	}
	if err := command.RunQuietly(ctx, "install", "-m", "0755", "-d", filepath.Dir(keyringPath)); err != nil {
		return fmt.Errorf("create Docker keyring directory: %w", err)
	}
	if err := command.RunQuietly(ctx, "curl", "-fsSL", dockerKeyURL, "-o", keyringPath); err != nil {
		return fmt.Errorf("download Docker signing key: %w", err)
	}
	if err := command.RunQuietly(ctx, "chmod", "a+r", keyringPath); err != nil {
		return fmt.Errorf("make Docker signing key readable: %w", err)
	}

	architecture, err := command.Run(ctx, "dpkg", "--print-architecture")
	if err != nil {
		return fmt.Errorf("determine Debian architecture: %w", err)
	}
	codename, err := ubuntu.Codename()
	if err != nil {
		return err
	}
	sources := fmt.Sprintf("Types: deb\nURIs: %s\nSuites: %s\nComponents: stable\nArchitectures: %s\nSigned-By: %s\n", dockerAPTURL, codename, architecture, keyringPath)
	if err := os.WriteFile(sourcesPath, []byte(sources), 0o644); err != nil {
		return fmt.Errorf("write Docker apt source: %w", err)
	}
	if err := os.Chmod(sourcesPath, 0o644); err != nil {
		return fmt.Errorf("set Docker apt source permissions: %w", err)
	}

	if err := command.RunQuietly(ctx, "apt-get", "update"); err != nil {
		return fmt.Errorf("update Docker package index: %w", err)
	}
	args := append([]string{"install", "-y"}, enginePackages...)
	if err := command.RunQuietly(ctx, "apt-get", args...); err != nil {
		return fmt.Errorf("install Docker Engine: %w", err)
	}
	return nil
}

func removeConflictingPackages(ctx context.Context) error {
	installed := make([]string, 0, len(conflictingPackages))
	for _, name := range conflictingPackages {
		status, err := command.Run(ctx, "dpkg-query", "--show", "--showformat=${db:Status-Status}", name)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			continue
		}
		if err != nil {
			return fmt.Errorf("check conflicting package %s: %w", name, err)
		}
		if status == "installed" {
			installed = append(installed, name)
		}
	}
	if len(installed) == 0 {
		return nil
	}
	args := append([]string{"remove", "-y"}, installed...)
	if err := command.RunQuietly(ctx, "apt-get", args...); err != nil {
		return fmt.Errorf("remove conflicting Docker packages: %w", err)
	}
	return nil
}
