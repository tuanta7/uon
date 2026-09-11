package main

import (
	"fmt"
	"os/exec"
)

func init() {
	checkBrew()
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
