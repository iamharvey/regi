package command

import (
	"bytes"
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestCmdLogin(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	streams := io.Streams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	ctxCmd := NewCmdLogin(streams)

	// Login.
	_, err := executeCommand(ctxCmd)
	assert.NoError(t, err)

}
