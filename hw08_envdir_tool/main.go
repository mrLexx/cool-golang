package main

import (
	"fmt"
	"os"
)

var (
	dir string
	cmd []string
)

func main() {
	if len(os.Args[1:]) < 2 {
		fmt.Printf("Usage of %v:\n", os.Args[0])
		fmt.Printf("  %v dir command arg1 arg2 ... argN:\n", os.Args[0])
		os.Exit(1)
	}

	dir = os.Args[1]
	cmd = os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(2)
	}

	code := RunCmd(cmd, env)
	os.Exit(code)
}
