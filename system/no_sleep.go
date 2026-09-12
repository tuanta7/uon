package system

import (
	"fmt"
	"os/exec"
)

func DisableSleep() error {
	commands := [][]string{
		{"pmset", "-a", "sleep", "0"},        // Disable system sleep
		{"pmset", "-a", "displaysleep", "0"}, // Disable display sleep
		{"pmset", "-a", "disksleep", "0"},    // Disable disk sleep
		{"pmset", "-a", "womp", "1"},         // Wake on Magic Packet (LAN)
		{"pmset", "-a", "autorestart", "1"},  // Automatically restart after power failure
	}

	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to execute command %v: %v, output: %s", command, err, output)
		}
		fmt.Println("Command executed successfully:", string(output))
	}
	return nil
}

func RestoreSleep(minutes int) error {
	commands := [][]string{
		{"pmset", "-a", "sleep", fmt.Sprintf("%d", minutes)},
		{"pmset", "-a", "displaysleep", fmt.Sprintf("%d", minutes)},
		{"pmset", "-a", "disksleep", fmt.Sprintf("%d", minutes)},
	}

	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to execute command %v: %v, output: %s", command, err, output)
		}
		fmt.Println("Command executed successfully:", string(output))
	}
	return nil
}
