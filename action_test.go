package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetOutput(t *testing.T) {
	t.Parallel()

	tempFileName, err := uuid.NewRandom()
	require.NoError(t, err)
	tempFile, err := os.CreateTemp("", tempFileName.String())
	require.NoError(t, err)
	defer func(name string) {
		err := os.Remove(name)
		require.NoError(t, err)
	}(tempFile.Name())

	err = os.Setenv("GITHUB_OUTPUT", tempFile.Name())
	require.NoError(t, err)

	tests := []struct {
		name  string
		value string
	}{
		{"TEST_VAR1", "value1"},
		{"TEST_VAR2", "value2"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s=%s", tc.name, tc.value), func(t *testing.T) {
			err := setOutput(tc.name, tc.value)
			require.NoError(t, err)

			fileContent, err := os.ReadFile(tempFile.Name())
			require.NoError(t, err)

			want := fmt.Sprintf("%s=%s", tc.name, tc.value)
			assert.True(t, strings.Contains(string(fileContent), want))
		})
	}
}
