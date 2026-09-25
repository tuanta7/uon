package ssh

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestSSHCommandRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"missing state", nil, "accepts 1 arg"},
		{"invalid state", []string{"on"}, "invalid argument"},
		{"too many arguments", []string{"enable", "disable"}, "accepts 1 arg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newSSHCommand()
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Execute() error = %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func TestSSHEnableInstallsMissingServerAfterConfirmation(t *testing.T) {
	var actions []string
	cmd := newSSHCommandWithActions(
		func(context.Context) (bool, error) { return false, nil },
		func(context.Context) error {
			actions = append(actions, "install")
			return nil
		},
		func(_ context.Context, allow bool) error {
			if !allow {
				t.Fatal("toggle called with allow=false")
			}
			actions = append(actions, "toggle")
			return nil
		},
	)
	cmd.SetIn(strings.NewReader("yes\n"))
	var output strings.Builder
	cmd.SetOut(&output)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"enable"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := strings.Join(actions, ","); got != "install,toggle" {
		t.Fatalf("actions = %q, want %q", got, "install,toggle")
	}
	if got := output.String(); !strings.Contains(got, "Install it now?") || !strings.Contains(got, "SSH access is enabled.") {
		t.Fatalf("output = %q, want installation prompt and success message", got)
	}
}

func TestSSHEnableDoesNotInstallWithoutConfirmation(t *testing.T) {
	installCalled := false
	toggleCalled := false
	cmd := newSSHCommandWithActions(
		func(context.Context) (bool, error) { return false, nil },
		func(context.Context) error {
			installCalled = true
			return nil
		},
		func(context.Context, bool) error {
			toggleCalled = true
			return nil
		},
	)
	cmd.SetIn(strings.NewReader("n\n"))
	var output strings.Builder
	cmd.SetOut(&output)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"enable"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if installCalled || toggleCalled {
		t.Fatalf("install called = %v, toggle called = %v; want neither", installCalled, toggleCalled)
	}
	if !strings.Contains(output.String(), "SSH access remains disabled") {
		t.Fatalf("output = %q, want declined message", output.String())
	}
}

func TestSSHDisableDoesNotCheckOrInstallPackage(t *testing.T) {
	cmd := newSSHCommandWithActions(
		func(context.Context) (bool, error) {
			t.Fatal("package check called for ssh disable")
			return false, nil
		},
		func(context.Context) error {
			t.Fatal("installer called for ssh disable")
			return nil
		},
		func(_ context.Context, allow bool) error {
			if allow {
				t.Fatal("toggle called with allow=true")
			}
			return nil
		},
	)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"disable"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}
