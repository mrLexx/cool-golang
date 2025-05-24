package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	run := exec.Command(cmd[0], cmd[1:]...)

	for _, e := range os.Environ() {
		k := strings.Split(e, "=")[0]
		if _, ok := env[k]; !ok {
			run.Env = append(run.Env, e)
		}
	}

	for k, v := range env {
		if !v.NeedRemove {
			run.Env = append(run.Env, fmt.Sprintf("%v=%v", k, v.Value))
		}
	}

	var out strings.Builder
	run.Stdout = &out
	err := run.Run()

	if exitError, ok := err.(*exec.ExitError); ok {
		exitCode := exitError.ExitCode()
		return exitCode
	}

	fmt.Printf("%v", out.String())

	return
}
