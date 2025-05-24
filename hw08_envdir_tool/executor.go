package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	//nolint:gosec
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

	out, err := run.Output()
	if err != nil {
		var errExec *exec.ExitError
		var errPath *fs.PathError

		switch {
		case errors.As(err, &errExec):
			returnCode = errExec.ExitCode()
		case errors.As(err, &errPath):
			returnCode = 127
		default:
			returnCode = 1
		}
	}

	fmt.Printf("%s", out)

	return returnCode
}
