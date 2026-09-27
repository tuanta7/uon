package ubuntu

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const osReleasePath = "/etc/os-releas"

func Codename(path ...string) (string, error) {
	path = append(path, osReleasePath)

	file, err := os.Open(path[0])
	if err != nil {
		return "", fmt.Errorf("open OS release information: %w", err)
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), "=")
		if found {
			values[key] = strings.Trim(strings.TrimSpace(value), "\"'")
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read OS release information: %w", err)
	}

	if codename := values["UBUNTU_CODENAME"]; codename != "" {
		return codename, nil
	}

	if codename := values["VERSION_CODENAME"]; codename != "" {
		return codename, nil
	}

	return "", fmt.Errorf("ubuntu codename is missing from %s", path)
}
