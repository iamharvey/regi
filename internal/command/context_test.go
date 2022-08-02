package command

import (
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCmdContext(t *testing.T) {
	streams := io.Streams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	ctxCmd := NewCmdContext(streams)
	assert.NotNil(t, ctxCmd)

	// Add test.
	_, err := executeCommand(ctxCmd, "add", "-n=context-123", "-s=http://localhost:5000", "-u=regi", "-p=regi")
	assert.NoError(t, err)

	// SetCurrentContext test.
	_, err = executeCommand(ctxCmd, "set", "context-123")
	assert.NoError(t, err)

	// ListContexts test.
	_, err = executeCommand(ctxCmd, "list")
	assert.NoError(t, err)

	// GetContext test.
	_, err = executeCommand(ctxCmd, "get", "context-123")
	assert.NoError(t, err)

}
