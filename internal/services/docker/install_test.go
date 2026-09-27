package docker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tuanta7/uon/pkg/ubuntu"
)

func TestUbuntuCodename(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		want     string
	}{
		{"Ubuntu codename", "VERSION_CODENAME=generic\nUBUNTU_CODENAME=noble\n", "noble"},
		{"version fallback", "VERSION_CODENAME=\"jammy\"\n", "jammy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "os-release")
			if err := os.WriteFile(path, []byte(tt.contents), 0o600); err != nil {
				t.Fatal(err)
			}
			got, err := ubuntu.Codename(path)
			if err != nil {
				t.Fatalf("ubuntuCodename() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ubuntuCodename() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUbuntuCodenameRequiresCodename(t *testing.T) {
	path := filepath.Join(t.TempDir(), "os-release")
	if err := os.WriteFile(path, []byte("NAME=Ubuntu\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ubuntu.Codename(path); err == nil {
		t.Fatal("ubuntuCodename() succeeded, want missing-codename error")
	}
}

func TestRemoveConflictingPackagesOnlyRemovesInstalledPackages(t *testing.T) {
	binDir := t.TempDir()
	argsFile := filepath.Join(t.TempDir(), "apt-args")
	writeCommand(t, binDir, "dpkg-query", `
for argument do package=$argument; done
case "$package" in
  docker.io|runc) printf installed ;;
  *) exit 1 ;;
esac
`)
	writeCommand(t, binDir, "apt-get", "printf '%s\\n' \"$@\" > \"$APT_ARGS\"\n")
	t.Setenv("PATH", binDir)
	t.Setenv("APT_ARGS", argsFile)

	if err := removeConflictingPackages(t.Context()); err != nil {
		t.Fatalf("removeConflictingPackages() error = %v", err)
	}
	contents, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(contents), "remove\n-y\ndocker.io\nrunc\n"; got != want {
		t.Fatalf("apt-get arguments = %q, want %q", got, want)
	}
}

func TestRunEnablesAndStartsDocker(t *testing.T) {
	binDir := t.TempDir()
	argsFile := filepath.Join(t.TempDir(), "systemctl-args")
	writeCommand(t, binDir, "systemctl", "printf '%s\\n' \"$@\" > \"$SYSTEMCTL_ARGS\"\n")
	t.Setenv("PATH", binDir)
	t.Setenv("SYSTEMCTL_ARGS", argsFile)

	if err := Run(t.Context()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	contents, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(contents), "enable\n--now\ndocker.service\n"; got != want {
		t.Fatalf("systemctl arguments = %q, want %q", got, want)
	}
}

func writeCommand(t *testing.T, directory, name, body string) {
	t.Helper()
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+strings.TrimSpace(body)+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}
