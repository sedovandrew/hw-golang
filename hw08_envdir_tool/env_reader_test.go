package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("parse env directory", func(t *testing.T) {
		environmentMap, err := ReadDir("testdata/env/")
		require.NoError(t, err)
		require.IsType(t, Environment{}, environmentMap)
		// BAR
		require.Equal(t, "bar", environmentMap["BAR"].Value)
		require.False(t, environmentMap["BAR"].NeedRemove)
		// EMPTY
		require.Equal(t, "", environmentMap["EMPTY"].Value)
		require.False(t, environmentMap["EMPTY"].NeedRemove)
		// FOO
		require.Equal(t, "   foo\nwith new line", environmentMap["FOO"].Value)
		require.False(t, environmentMap["FOO"].NeedRemove)
		// HELLO
		require.Equal(t, "\"hello\"", environmentMap["HELLO"].Value)
		require.False(t, environmentMap["HELLO"].NeedRemove)
		// UNSET
		require.Equal(t, "", environmentMap["UNSET"].Value)
		require.True(t, environmentMap["UNSET"].NeedRemove)
		// TAB
		require.Equal(t, "Tabs at the end of the line", environmentMap["TAB"].Value)
		require.False(t, environmentMap["TAB"].NeedRemove)
		// A=A
		require.NotContains(t, environmentMap, "A=A")
	})
}
