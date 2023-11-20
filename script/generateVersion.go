package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	stdout, err := cmd.Output()
	if err != nil {
		return
	}

	version := []byte(strings.TrimSpace(string(stdout)))
	err = os.WriteFile("assets/version", version, 0600)
	if err != nil {
		return
	}
}
