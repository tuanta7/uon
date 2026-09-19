package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tuanta7/uon/internal/system/sleep"
)

func newSleepCommand() *cobra.Command {
	return &cobra.Command{
		Use:       "sleep <on|off>",
		Short:     "Allow or prevent system sleep",
		Long:      "Turn system sleep on to allow suspend and hibernation, or off to prevent them.",
		ValidArgs: []string{"on", "off"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			allowSleep := args[0] == "on"
			if err := sleep.ToggleSleep(cmd.Context(), allowSleep); err != nil {
				return err
			}

			message := "System sleep is off."
			if allowSleep {
				message = "System sleep is on."
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), message)
			return err
		},
	}
}
