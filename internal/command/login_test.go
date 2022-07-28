package command

import (
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCmdLogin(t *testing.T) {
	streams := io.Streams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	ctxCmd := NewCmdLogin(streams)
	assert.NotNil(t, ctxCmd)

	// Login.
	_, err := executeCommand(ctxCmd)
	assert.NoError(t, err)

}
