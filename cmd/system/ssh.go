package system

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

func newSSHCommand() *cobra.Command {
	return newSSHCommandWithActions(remote.IsOpenSSHServerInstalled, remote.InstallOpenSSHServer, remote.ToggleSSH)
}

func newSSHCommandWithActions(
	isInstalled func(context.Context) (bool, error),
	install func(context.Context) error,
	toggle func(context.Context, bool) error,
) *cobra.Command {
	return &cobra.Command{
		Use:       "ssh <on|off>",
		Short:     "Allow or prevent SSH access",
		Long:      "Turn SSH on to start it and enable it at boot, or off to stop it and disable it at boot.",
		ValidArgs: []string{"on", "off"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			allowSSH := args[0] == "on"
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
						_, err := fmt.Fprintln(cmd.OutOrStdout(), "OpenSSH server was not installed; SSH access remains off.")
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

			message := "SSH access is off."
			if allowSSH {
				message = "SSH access is on."
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
