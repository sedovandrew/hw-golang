package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReturnCode(t *testing.T) {
	if os.Getenv("MAIN_TEST_RETURN_CODE") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0], "testdata/env", "testdata/rc.sh"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestReturnCode")
	cmd.Env = append(os.Environ(), "MAIN_TEST_RETURN_CODE=1")
	err := cmd.Run()

	var exitError *exec.ExitError
	ok := errors.As(err, &exitError)
	require.True(t, ok)
	require.Equal(t, "exit status 7", exitError.Error())
}

func TestPipe(t *testing.T) {
	// Run pipe.sh
	if os.Getenv("MAIN_TEST_PIPE") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0], "testdata/env", "testdata/pipe.sh"}
		main()
		return
	}

	// Run test one more time but with TEST_MAIN=1
	cmd := exec.Command(os.Args[0], "-test.run=TestPipe")
	cmd.Env = append(os.Environ(), "MAIN_TEST_PIPE=1")

	// Get the stdin of the command being run.
	stdinCmd, err := cmd.StdinPipe()
	require.NoError(t, err)

	// Get the stdout of the command being run.
	stdoutCmd, err := cmd.StdoutPipe()
	require.NoError(t, err)

	// Forwarding an stderr directly.
	cmd.Stderr = os.Stderr

	// Start command
	err = cmd.Start()
	require.NoError(t, err)

	// Writing to stdin
	io.WriteString(stdinCmd, "He__o Wor_d!")
	stdinCmd.Close()

	// Reading from stdout
	stdoutBytes, err := io.ReadAll(stdoutCmd)
	require.NoError(t, err)

	cmd.Wait()
	require.Equal(t, "Hello World!", string(stdoutBytes))
}

func TestArgs(t *testing.T) {
	if os.Getenv("MAIN_TEST_ARGS") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0], "testdata/env", "testdata/args.sh", "a=1", "b=2"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestArgs")
	cmd.Env = append(os.Environ(), "MAIN_TEST_ARGS=1")
	var outBytes bytes.Buffer
	cmd.Stdout = &outBytes
	err := cmd.Run()

	var exitError *exec.ExitError
	ok := errors.As(err, &exitError)
	require.False(t, ok)
	require.Equal(t, "a=1 b=2", outBytes.String())
}

func TestWithoutArgs(t *testing.T) {
	if os.Getenv("MAIN_TEST_WITHOUT_ARGS") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0]}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWithoutArgs")
	cmd.Env = append(os.Environ(), "MAIN_TEST_WITHOUT_ARGS=1")
	err := cmd.Run()

	var exitError *exec.ExitError
	ok := errors.As(err, &exitError)
	require.True(t, ok)
	require.Equal(t, "exit status 111", exitError.Error())
}

func TestOnlyOneArg(t *testing.T) {
	if os.Getenv("MAIN_TEST_ONLY_ONE_ARG") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0], "testdata/env"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestOnlyOneArg")
	cmd.Env = append(os.Environ(), "MAIN_TEST_ONLY_ONE_ARG=1")
	err := cmd.Run()

	var exitError *exec.ExitError
	ok := errors.As(err, &exitError)
	require.True(t, ok)
	require.Equal(t, "exit status 111", exitError.Error())
}

func TestDirNotExist(t *testing.T) {
	if os.Getenv("MAIN_TEST_DIR_NOT_EXIST") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0], "testdata/not_exists", "ls"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDirNotExist")
	cmd.Env = append(os.Environ(), "MAIN_TEST_DIR_NOT_EXIST=1")
	err := cmd.Run()

	var exitError *exec.ExitError
	ok := errors.As(err, &exitError)
	require.True(t, ok)
	require.Equal(t, "exit status 111", exitError.Error())
}
