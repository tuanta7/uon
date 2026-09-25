package sleep

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/internal/system/sleep"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:       "sleep <enable|disable>",
		Short:     "Allow or prevent system sleep",
		Long:      "Enable system sleep to allow suspend and hibernation, or disable it to prevent them.",
		ValidArgs: []string{"enable", "disable"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			allowSleep := args[0] == "enable"
			if err := sleep.ToggleSleep(cmd.Context(), allowSleep); err != nil {
				return err
			}

			message := "System sleep is disabled."
			if allowSleep {
				message = "System sleep is enabled."
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), message)
			return err
		},
	}
}
