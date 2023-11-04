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
		stdout = []byte{}
	}

	version := []byte(strings.TrimSpace(string(stdout)))
	os.WriteFile("assets/version", version, 0644)
}
