package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage: %s /path/to/env/dir command arg1 arg2\n", args[0])
		return
	}
	dir := args[1]
	cmds := args[2:]
	envs, err := ReadDir(dir)
	if err != nil {
		panic(err)
	}
	os.Exit(RunCmd(cmds, envs))
}
