package system

import (
	"fmt"
	"os/exec"
)

func AllowSSH() error {
	cmd := exec.Command("systemsetup", "-setremotelogin", "on")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute command %v: %v, output: %s", cmd.Args, err, output)
	}
	fmt.Println("Command executed successfully:", string(output))
	return nil
}
