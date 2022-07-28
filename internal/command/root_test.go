package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCmdRoot(t *testing.T) {
	rootCmd := NewRegiCommand()
	assert.NotNil(t, rootCmd)

	/*  There are actually 3 commands:
	- context		Manage connection settings of multiple Docker registries.
	- image			Pull, push, delete and list images over Docker registry
	- login			Login to current Docker registry.
	*/
	assert.Equal(t, 3, len(rootCmd.Commands()))
	assert.Equal(t, msgShort, rootCmd.Short)
	assert.Equal(t, msgLong, rootCmd.Long)

	// Root.
	_, err := executeCommand(rootCmd)
	assert.NoError(t, err)
}
