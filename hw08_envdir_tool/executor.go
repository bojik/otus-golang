package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 0
	}
	for name, val := range env {
		if val.NeedRemove {
			os.Unsetenv(name)
			continue
		}
		os.Setenv(name, val.Value)
	}
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Env = os.Environ()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return err.(*exec.ExitError).ExitCode() //nolint:errorlint
		}
	}
	return 0
}
