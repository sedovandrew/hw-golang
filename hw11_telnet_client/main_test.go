package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArgs(t *testing.T) {
	if os.Getenv("OTUS_TELNET_MAIN_TEST_ARGS") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0]}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestArgs")
	cmd.Env = append(os.Environ(), "OTUS_TELNET_MAIN_TEST_ARGS=1")
	err := cmd.Run()

	var exitError *exec.ExitError
	ok := errors.As(err, &exitError)
	require.True(t, ok)
	require.Equal(t, "exit status 1", exitError.Error())
}

func TestWrongPort(t *testing.T) {
	if os.Getenv("OTUS_TELNET_MAIN_TEST_WRONG_PORT") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0], "localhost", "wrong_network_port"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWrongPort")
	cmd.Env = append(os.Environ(), "OTUS_TELNET_MAIN_TEST_WRONG_PORT=1")
	err := cmd.Run()

	var exitError *exec.ExitError
	ok := errors.As(err, &exitError)
	require.True(t, ok)
	require.Equal(t, "exit status 1", exitError.Error())
}

func TestWrongPortRange(t *testing.T) {
	if os.Getenv("OTUS_TELNET_MAIN_TEST_WRONG_PORT_RANGE") == "1" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()
		os.Args = []string{os.Args[0], "localhost", "65537"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWrongPortRange")
	cmd.Env = append(os.Environ(), "OTUS_TELNET_MAIN_TEST_WRONG_PORT_RANGE=1")
	err := cmd.Run()

	var exitError *exec.ExitError
	ok := errors.As(err, &exitError)
	require.True(t, ok)
	require.Equal(t, "exit status 1", exitError.Error())
}
