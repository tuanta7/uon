package nginx

import (
	"context"
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
		"load":    false,
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

func TestLoadCommandRequiresOneFilePath(t *testing.T) {
	for _, args := range [][]string{{"load"}, {"load", "one.conf", "two.conf"}} {
		cmd := NewCommand()
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		cmd.SetArgs(args)
		if err := cmd.Execute(); err == nil {
			t.Fatalf("Execute() with args %v succeeded, want argument error", args)
		}
	}
}

func TestLoadCommandSelectsStreamConfiguration(t *testing.T) {
	var gotPath string
	var gotStream bool
	cmd := newLoadCommand(func(_ context.Context, path string, stream bool) error {
		gotPath = path
		gotStream = stream
		return nil
	})
	cmd.SetOut(io.Discard)
	cmd.SetArgs([]string{"--stream", "redis.conf"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotPath != "redis.conf" || !gotStream {
		t.Fatalf("load arguments = (%q, %v), want (%q, true)", gotPath, gotStream, "redis.conf")
	}
}

func TestInstallCommandPrintsConfigNextStep(t *testing.T) {
	binDir := t.TempDir()
	argsFile := filepath.Join(t.TempDir(), "apt-args")
	writeCommand(t, binDir, "apt-get", "printf '%s\\n' \"$@\" > \"$APT_ARGS\"\n")
	t.Setenv("PATH", binDir)
	t.Setenv("APT_ARGS", argsFile)

	cmd := installCommand()
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := output.String(); !strings.Contains(got, "NGINX is installed") || !strings.Contains(got, "uon nginx config") {
		t.Fatalf("output = %q, want success and config guidance", got)
	}
	assertFileContents(t, argsFile, "install\n-y\nnginx\nlibnginx-mod-stream\n")
}

func TestRunCommandRequiresInstalledPackage(t *testing.T) {
	binDir := t.TempDir()
	writeCommand(t, binDir, "dpkg-query", "exit 1\n")
	t.Setenv("PATH", binDir)

	cmd := runCommand()
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

	cmd := runCommand()
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
	cmd := configCommand()
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

func TestRemoveCommandPreservesConfigurationByDefault(t *testing.T) {
	binDir := t.TempDir()
	argsFile := filepath.Join(t.TempDir(), "apt-args")
	writeCommand(t, binDir, "apt-get", "printf '%s\\n' \"$@\" > \"$APT_ARGS\"\n")
	t.Setenv("PATH", binDir)
	t.Setenv("APT_ARGS", argsFile)

	cmd := removeCommand()
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	assertFileContents(t, argsFile, "remove\n-y\nnginx\nlibnginx-mod-stream\n")
	if !strings.Contains(output.String(), "Configuration files were preserved") {
		t.Fatalf("output = %q, want preserved configuration message", output.String())
	}
}

func TestRemoveCommandPrunesConfiguration(t *testing.T) {
	binDir := t.TempDir()
	argsFile := filepath.Join(t.TempDir(), "apt-args")
	writeCommand(t, binDir, "apt-get", "printf '%s\\n' \"$@\" > \"$APT_ARGS\"\n")
	t.Setenv("PATH", binDir)
	t.Setenv("APT_ARGS", argsFile)

	cmd := removeCommand()
	cmd.SetArgs([]string{"--prune"})
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	assertFileContents(t, argsFile, "purge\n-y\nnginx\nnginx-common\nlibnginx-mod-stream\n")
	if !strings.Contains(output.String(), "Configuration files were pruned") {
		t.Fatalf("output = %q, want pruned configuration message", output.String())
	}
}

func TestCommandPropagatesActionError(t *testing.T) {
	binDir := t.TempDir()
	writeCommand(t, binDir, "apt-get", "echo 'apt failed' >&2\nexit 1\n")
	t.Setenv("PATH", binDir)

	cmd := installCommand()
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

func assertFileContents(t *testing.T, path, want string) {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(contents); got != want {
		t.Fatalf("file contents = %q, want %q", got, want)
	}
}
