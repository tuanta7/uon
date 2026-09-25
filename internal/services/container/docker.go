package container

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

var packagesToRemove = []string{
	"docker.io",
	"docker-compose",
	"docker-compose-v2",
	"docker-doc",
	"docker-buildx",
	"podman-docker",
	"containerd",
	"runc",
}

type DockerManager struct {
	version string
}

func NewDockerManager() DockerManager {
	cmd := exec.Command("docker", "--version")
	output, err := cmd.Output()
	if err != nil {
		return DockerManager{}
	}

	return DockerManager{
		version: strings.TrimSpace(string(output)),
	}
}

func (d *DockerManager) Install(ctx context.Context) error {
	dpkgCmd := fmt.Sprintf("dpkg --get-selections %s | cut -f1", strings.Join(packagesToRemove, " "))
	cmd := exec.CommandContext(ctx, "apt", "remove", fmt.Sprintf("$(%s)", dpkgCmd))
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (d *DockerManager) Update(ctx context.Context, version string) error {
	return nil
}

// Uninstall implements the steps defined at
// https://docs.docker.com/engine/install/ubuntu/#uninstall-docker-engine
func (d *DockerManager) Uninstall(ctx context.Context) error {
	if d.version == "" {
		return nil
	}

	cmd := exec.CommandContext(ctx, "apt", "purge",
		"docker-ce",
		"docker-ce-cli",
		"containerd.io",
		"docker-buildx-plugin",
		"docker-compose-plugin",
		"docker-ce-rootless-extras",
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	removeDockerCmd := exec.CommandContext(ctx, "rm", "-rf", "/var/lib/docker")
	if err := removeDockerCmd.Run(); err != nil {
		return err
	}

	removeContainerdCmd := exec.CommandContext(ctx, "rm", "-rf", "/var/lib/containerd")
	if err := removeContainerdCmd.Run(); err != nil {
		return err
	}

	return nil
}
