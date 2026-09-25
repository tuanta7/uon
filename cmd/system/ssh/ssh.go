package ssh

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/internal/system/remote"
)

func Command() *cobra.Command {
	return newSSHCommandWithActions(remote.IsOpenSSHServerInstalled, remote.InstallOpenSSHServer, remote.ToggleSSH)
}

func newSSHCommandWithActions(
	isInstalled func(context.Context) (bool, error),
	install func(context.Context) error,
	toggle func(context.Context, bool) error,
) *cobra.Command {
	return &cobra.Command{
		Use:       "ssh <enable|disable>",
		Short:     "Allow or prevent SSH access",
		Long:      "Enable SSH to start it now and at boot, or disable it to stop it now and at boot.",
		ValidArgs: []string{"enable", "disable"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			allowSSH := args[0] == "enable"
			if allowSSH {
				installed, err := isInstalled(cmd.Context())
				if err != nil {
					return err
				}
				if !installed {
					confirmed, err := confirmOpenSSHInstall(cmd)
					if err != nil {
						return err
					}
					if !confirmed {
						_, err := fmt.Fprintln(cmd.OutOrStdout(), "OpenSSH server was not installed; SSH access remains disabled.")
						return err
					}
					if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Installing openssh-server..."); err != nil {
						return err
					}
					if err := install(cmd.Context()); err != nil {
						return err
					}
				}
			}

			if err := toggle(cmd.Context(), allowSSH); err != nil {
				return err
			}

			message := "SSH access is disabled."
			if allowSSH {
				message = "SSH access is enabled."
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), message)
			return err
		},
	}
}

func confirmOpenSSHInstall(cmd *cobra.Command) (bool, error) {
	if _, err := fmt.Fprint(cmd.OutOrStdout(), "openssh-server is not installed. Install it now? [y/N]: "); err != nil {
		return false, err
	}
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("read confirmation: %w", err)
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}
