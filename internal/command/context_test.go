package command

import (
	"bytes"
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestCmdContext(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	streams := io.Streams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	ctxCmd := NewCmdContext(streams)

	// Add test.
	_, err := executeCommand(ctxCmd, "add", "-n=context-123", "-s=http://localhost:5000")
	assert.NoError(t, err)

	// Set test.
	_, err = executeCommand(ctxCmd, "set", "context-123")
	assert.NoError(t, err)

	// List test.
	_, err = executeCommand(ctxCmd, "list")
	assert.NoError(t, err)

	// Get test.
	_, err = executeCommand(ctxCmd, "get", "context-123")
	assert.NoError(t, err)

}
