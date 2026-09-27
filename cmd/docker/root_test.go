package docker

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestCommandContainsLifecycleActions(t *testing.T) {
	cmd := NewCommand()
	found := map[string]bool{"install": false, "run": false, "remove": false}
	for _, child := range cmd.Commands() {
		if _, ok := found[child.Name()]; ok {
			found[child.Name()] = true
		}
	}
	for name, present := range found {
		if !present {
			t.Errorf("docker command is missing %q subcommand", name)
		}
	}
}

func TestLifecycleCommandsRejectArguments(t *testing.T) {
	for _, name := range []string{"install", "run", "remove"} {
		t.Run(name, func(t *testing.T) {
			cmd := NewCommand()
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.SetArgs([]string{name, "extra"})
			if err := cmd.Execute(); err == nil {
				t.Fatal("Execute() succeeded, want argument error")
			}
		})
	}
}

func TestInstallAndRunCommandsInvokeActions(t *testing.T) {
	for _, tt := range []struct {
		name string
		cmd  interface {
			Execute() error
			SetOut(io.Writer)
		}
		want string
	}{
		{"install", newInstallCommand(func(context.Context) error { return nil }), "Docker Engine is installed"},
		{"run", newRunCommand(func(context.Context) error { return nil }), "running and enabled at boot"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var output strings.Builder
			tt.cmd.SetOut(&output)
			if err := tt.cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if !strings.Contains(output.String(), tt.want) {
				t.Fatalf("output = %q, want %q", output.String(), tt.want)
			}
		})
	}
}

func TestRemoveCommandPassesPruneFlag(t *testing.T) {
	var gotPrune bool
	cmd := newRemoveCommand(func(_ context.Context, prune bool) error {
		gotPrune = prune
		return nil
	})
	cmd.SetArgs([]string{"--prune"})
	var output strings.Builder
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !gotPrune {
		t.Fatal("remove action received prune=false, want true")
	}
	if !strings.Contains(output.String(), "stored images, containers, and volumes were removed") {
		t.Fatalf("output = %q, want pruned-data message", output.String())
	}
}

func TestLifecycleCommandsPropagateErrors(t *testing.T) {
	want := errors.New("action failed")
	commands := []interface{ Execute() error }{
		newInstallCommand(func(context.Context) error { return want }),
		newRunCommand(func(context.Context) error { return want }),
		newRemoveCommand(func(context.Context, bool) error { return want }),
	}
	for _, cmd := range commands {
		if err := cmd.Execute(); !errors.Is(err, want) {
			t.Fatalf("Execute() error = %v, want %v", err, want)
		}
	}
}
