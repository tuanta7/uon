package dev

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestCommandContainsLifecycleActions(t *testing.T) {
	cmd := NewCommand()
	found := map[string]bool{"up": false, "down": false}
	for _, child := range cmd.Commands() {
		if _, ok := found[child.Name()]; ok {
			found[child.Name()] = true
		}
	}
	for name, present := range found {
		if !present {
			t.Errorf("dev command is missing %q subcommand", name)
		}
	}
}

func TestLifecycleCommandsRejectArguments(t *testing.T) {
	for _, name := range []string{"up", "down"} {
		t.Run(name, func(t *testing.T) {
			cmd := NewCommand()
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.SetArgs([]string{name, "extra"})
			if err := cmd.Execute(); err == nil ||
				(!strings.Contains(err.Error(), "accepts 0 arg") && !strings.Contains(err.Error(), "unknown command")) {
				t.Fatalf("Execute() error = %v, want argument error", err)
			}
		})
	}
}

func TestUpRunsComposeAndPrintsEndpoints(t *testing.T) {
	var gotName string
	var gotArgs []string
	cmd := newUpCommand(func(_ context.Context, name string, args ...string) error {
		gotName = name
		gotArgs = args
		return nil
	})
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotName != "docker" {
		t.Fatalf("command = %q, want docker", gotName)
	}
	wantArgs := []string{"compose", "--file", composeFile, "up", "--detach", "--wait"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("arguments = %q, want %q", gotArgs, wantArgs)
	}
	for _, endpoint := range []string{"<server-LAN-IP>:9000", "<server-LAN-IP>:9001", "<server-LAN-IP>:6379"} {
		if !strings.Contains(output.String(), endpoint) {
			t.Errorf("output = %q, want endpoint %q", output.String(), endpoint)
		}
	}
	for _, nextStep := range []string{"uon nginx install", "uon nginx run", "uon nginx load --stream static/uon.dev.nginx.conf"} {
		if !strings.Contains(output.String(), nextStep) {
			t.Errorf("output = %q, want next step %q", output.String(), nextStep)
		}
	}
}

func TestDownRunsComposeWithVolumeRemoval(t *testing.T) {
	var gotArgs []string
	cmd := newDownCommand(func(_ context.Context, name string, args ...string) error {
		if name != "docker" {
			t.Fatalf("command = %q, want docker", name)
		}
		gotArgs = args
		return nil
	})
	cmd.SetOut(io.Discard)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	wantArgs := []string{"compose", "--file", composeFile, "down", "--volumes", "--remove-orphans"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("arguments = %q, want %q", gotArgs, wantArgs)
	}
}

func TestLifecycleCommandsWrapComposeErrors(t *testing.T) {
	composeError := errors.New("compose failed")
	tests := []struct {
		name string
		cmd  interface{ Execute() error }
		want string
	}{
		{"up", newUpCommand(func(context.Context, string, ...string) error { return composeError }), "start development infrastructure"},
		{"down", newDownCommand(func(context.Context, string, ...string) error { return composeError }), "stop development infrastructure"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Execute()
			if !errors.Is(err, composeError) || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Execute() error = %v, want wrapped compose error containing %q", err, tt.want)
			}
		})
	}
}
