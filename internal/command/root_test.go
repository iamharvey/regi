package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCmdRoot(t *testing.T) {
	rootCmd := NewRegiCommand()
	assert.NotNil(t, rootCmd)

	/*  There are actually 5 commands:
	- completion		Generate the autocompletion script for the specified shell
	- context		Manage connection settings of multiple Docker registries.
	- help			Help about any command
	- image			Pull, push, delete and list images over Docker registry
	- login			Login to current Docker registry.
	*/
	assert.Equal(t, 5, len(rootCmd.Commands()))
	assert.Equal(t, msgShort, rootCmd.Short)
	assert.Equal(t, msgShort, rootCmd.Long)

	// Root.
	_, err := executeCommand(rootCmd)
	assert.NoError(t, err)
}
