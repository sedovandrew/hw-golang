package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("return code", func(t *testing.T) {
		command := []string{"testdata/rc.sh"}
		returnCode := RunCmd(command, Environment{})
		require.Equal(t, 7, returnCode)
	})
}
