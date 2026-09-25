package remote

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToggleSSH(t *testing.T) {
	tests := []struct {
		name  string
		allow bool
		want  string
	}{
		{"allow", true, "enable --now ssh.service"},
		{"prevent", false, "disable --now ssh.service"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			binDir := t.TempDir()
			argsFile := filepath.Join(t.TempDir(), "args")
			script := "#!/bin/sh\nprintf '%s' \"$*\" > \"$SSH_ARGS_FILE\"\n"
			if err := os.WriteFile(filepath.Join(binDir, "systemctl"), []byte(script), 0o755); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", binDir)
			t.Setenv("SSH_ARGS_FILE", argsFile)

			if err := ToggleSSH(context.Background(), tt.allow); err != nil {
				t.Fatalf("ToggleSSH() error = %v", err)
			}
			got, err := os.ReadFile(argsFile)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("systemctl arguments = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsOpenSSHServerInstalled(t *testing.T) {
	tests := []struct {
		name       string
		scriptBody string
		want       bool
	}{
		{"installed", "printf installed\n", true},
		{"missing", "exit 1\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			binDir := t.TempDir()
			script := "#!/bin/sh\n" + tt.scriptBody
			if err := os.WriteFile(filepath.Join(binDir, "dpkg-query"), []byte(script), 0o755); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", binDir)

			got, err := IsOpenSSHServerInstalled(context.Background())
			if err != nil {
				t.Fatalf("IsOpenSSHServerInstalled() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("IsOpenSSHServerInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstallOpenSSHServer(t *testing.T) {
	binDir := t.TempDir()
	argsFile := filepath.Join(t.TempDir(), "args")
	script := "#!/bin/sh\nprintf '%s' \"$*\" > \"$SSH_ARGS_FILE\"\n"
	if err := os.WriteFile(filepath.Join(binDir, "apt-get"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir)
	t.Setenv("SSH_ARGS_FILE", argsFile)

	if err := InstallOpenSSHServer(context.Background()); err != nil {
		t.Fatalf("InstallOpenSSHServer() error = %v", err)
	}
	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	if want := "install -y openssh-server"; string(got) != want {
		t.Fatalf("apt-get arguments = %q, want %q", got, want)
	}
}

func TestToggleSSHWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := ToggleSSH(ctx, true)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ToggleSSH() error = %v, want context.Canceled", err)
	}
}

func TestToggleSSHIncludesSystemctlError(t *testing.T) {
	binDir := t.TempDir()
	script := "#!/bin/sh\necho 'service not found' >&2\nexit 1\n"
	if err := os.WriteFile(filepath.Join(binDir, "systemctl"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir)

	err := ToggleSSH(context.Background(), true)
	if err == nil || !strings.Contains(err.Error(), "enable SSH service") || !strings.Contains(err.Error(), "service not found") {
		t.Fatalf("ToggleSSH() error = %v, want action and systemctl output", err)
	}
}
