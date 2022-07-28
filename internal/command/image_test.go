package command

import (
	"bytes"
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestCmdImage(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	streams := io.Streams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	ctxCmd := NewCmdImage(streams)

	// Docker pull hello-world image for later push test.
	// Execute command.
	c := exec.Command(
		"docker",
		"pull",
		"hello-world:latest",
	)
	_, err := c.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	// Push.
	_, err = executeCommand(ctxCmd, "push", "hello-world", "latest")
	assert.NoError(t, err)

	// Pull.
	_, err = executeCommand(ctxCmd, "pull", "hello-world", "latest")
	assert.NoError(t, err)

	// List.
	_, err = executeCommand(ctxCmd, "list", "withTag=true")
	assert.NoError(t, err)

	// Delete.
	_, err = executeCommand(ctxCmd, "delete", "hello-world", "latest")
	assert.NoError(t, err)

}
