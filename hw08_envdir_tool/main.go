package main

import (
	"fmt"
	"os"
)

const (
	// The return code is the same as for the envdir utility.
	errorReturnCode = 111
	// cmd + 2 args.
	minArgCount = 3
)

func main() {
	if len(os.Args) < minArgCount {
		fmt.Println("Not enough arguments.")
		os.Exit(errorReturnCode)
	}
	dir := os.Args[1]
	cmd := os.Args[2:]

	environmentMap, err := ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory %q.\n", dir)
		os.Exit(errorReturnCode)
	}

	returnCode := RunCmd(cmd, environmentMap)
	os.Exit(returnCode)
}
