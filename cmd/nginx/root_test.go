package nginx

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandContainsLifecycleAndConfigActions(t *testing.T) {
	cmd := NewCommand()
	want := map[string]bool{
		"install": false,
		"run":     false,
		"status":  false,
		"remove":  false,
		"config":  false,
	}
	for _, child := range cmd.Commands() {
		if _, ok := want[child.Name()]; ok {
			want[child.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("nginx command is missing %q subcommand", name)
		}
	}
}

func TestNGINXSubcommandsRejectArguments(t *testing.T) {
	for _, name := range []string{"install", "run", "status", "remove", "config"} {
		t.Run(name, func(t *testing.T) {
			cmd := NewCommand()
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.SetArgs([]string{name, "extra"})
			if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "unknown command") && !strings.Contains(err.Error(), "accepts 0 arg") {
				t.Fatalf("Execute() error = %v, want argument error", err)
			}
		})
	}
}

func TestInstallCommandPrintsConfigNextStep(t *testing.T) {
	binDir := t.TempDir()
	writeCommand(t, binDir, "apt-get", "exit 0\n")
	t.Setenv("PATH", binDir)

	cmd := newInstallCommand()
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := output.String(); !strings.Contains(got, "NGINX is installed") || !strings.Contains(got, "uon nginx config") {
		t.Fatalf("output = %q, want success and config guidance", got)
	}
}

func TestRunCommandRequiresInstalledPackage(t *testing.T) {
	binDir := t.TempDir()
	writeCommand(t, binDir, "dpkg-query", "exit 1\n")
	t.Setenv("PATH", binDir)

	cmd := newRunCommand()
	cmd.SetOut(io.Discard)

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "uon nginx install") {
		t.Fatalf("Execute() error = %v, want install-first error", err)
	}
}

func TestRunCommandStartsInstalledNGINX(t *testing.T) {
	binDir := t.TempDir()
	writeCommand(t, binDir, "dpkg-query", "printf installed\n")
	writeCommand(t, binDir, "systemctl", "exit 0\n")
	t.Setenv("PATH", binDir)

	cmd := newRunCommand()
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(output.String(), "running and enabled at boot") {
		t.Fatalf("output = %q, want run success", output.String())
	}
}

func TestStatusCommandFormatsAllStates(t *testing.T) {
	binDir := t.TempDir()
	writeCommand(t, binDir, "dpkg-query", "printf installed\n")
	writeCommand(t, binDir, "systemctl", "case \"$1\" in\nis-active) exit 3 ;;\nis-enabled) exit 0 ;;\nesac\nexit 2\n")
	t.Setenv("PATH", binDir)

	cmd := newStatusCommand()
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if want := "Installed: yes\nActive: no\nEnabled: yes\n"; output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestConfigCommandShowsUbuntuWorkflow(t *testing.T) {
	cmd := newConfigCommand()
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{
		"/etc/nginx/nginx.conf",
		"/etc/nginx/sites-available/",
		"/etc/nginx/sites-enabled/",
		"nginx -t",
		"systemctl reload nginx",
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("output = %q, want it to contain %q", output.String(), want)
		}
	}
}

func TestCommandPropagatesActionError(t *testing.T) {
	binDir := t.TempDir()
	writeCommand(t, binDir, "apt-get", "echo 'apt failed' >&2\nexit 1\n")
	t.Setenv("PATH", binDir)

	cmd := newInstallCommand()
	cmd.SetOut(io.Discard)

	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "apt failed") {
		t.Fatalf("Execute() error = %v, want apt failure", err)
	}
}

func writeCommand(t *testing.T, directory, name, body string) {
	t.Helper()
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
}
