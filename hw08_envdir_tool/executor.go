package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) int {
	for envKey, envValue := range env {
		if envValue.NeedRemove {
			os.Unsetenv(envKey)
		} else {
			os.Setenv(envKey, envValue.Value)
		}
	}

	c := exec.Command(cmd[0], cmd[1:]...) //#nosec G204
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	c.Run()

	returnCode := c.ProcessState.ExitCode()
	return returnCode
}
