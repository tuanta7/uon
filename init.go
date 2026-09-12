package main

import (
	"fmt"
	"minion/system"
	"os"

	"os/exec"
)

func init() {
	if os.Geteuid() != 0 {
		fmt.Println("This program must be run as an administrator.")
		os.Exit(1)
	}
	checkBrew()
	alwaysOn()
}

func checkBrew() {
	_, err := exec.LookPath("brew")
	if err != nil {
		fmt.Println("Homebrew is not installed")
		return
	}

	out, err := exec.Command("brew", "--version").Output()
	if err != nil {
		fmt.Println("Failed to get Homebrew version")
	} else {
		fmt.Println(string(out))
	}
}

func alwaysOn() {
	err := system.DisableSleep()
	if err != nil {
		fmt.Println("Failed to disable sleep:", err)
	}
}
