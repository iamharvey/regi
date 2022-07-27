package command

import (
	"bytes"
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestNewCmdImage(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	streams := io.Streams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	ctxCmd := NewCmdImage(streams)

	// Add test.
	_, err := executeCommand(ctxCmd, "list", "--withTag=true")
	assert.NoError(t, err)
}
