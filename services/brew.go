package services

import (
	"fmt"
	"os/exec"
)

const (
	defaultVersion = "6.0.22" // updated 12 Sep 2026
	downloadURL    = "https://github.com/Homebrew/brew/releases/download/%s/Homebrew.pkg"
	pkgPath        = "Homebrew.pkg"
)

func InstallBrewGitHub(version string) error {
	cmd := exec.Command("curl", "-L", "-O", fmt.Sprintf(downloadURL, version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute command %v: %v, output: %s", cmd.Args, err, output)
	}
	fmt.Println("Command executed successfully:", string(output))

	cmd = exec.Command("installer", "-pkg", pkgPath, "-target", "/")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute command %v: %v, output: %s", cmd.Args, err, output)
	}

	fmt.Println("Command executed successfully:", string(output))
	return nil
}
