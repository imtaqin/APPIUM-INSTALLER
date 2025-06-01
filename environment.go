package main

import (
	"fmt"
	"os/exec"
)

func setEnvVarPermanent(name, value string) error {
	return runCommandWait("setx", name, value)
}

func addToPathPermanent(newPath string) error {
	return runCommandWait("setx", "PATH", fmt.Sprintf("%%PATH%%;%s", newPath))
}

func refreshEnvironmentCmd() error {
	cmd := exec.Command("cmd", "/C", "refreshenv")
	return cmd.Run()
}