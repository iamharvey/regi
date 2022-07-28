package command

import (
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestCmdImage(t *testing.T) {
	streams := io.Streams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	imgCmd := NewCmdImage(streams)
	assert.NotNil(t, imgCmd)

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
	_, err = executeCommand(imgCmd, "push", "hello-world", "latest")
	assert.NoError(t, err)

	// Pull.
	_, err = executeCommand(imgCmd, "pull", "hello-world", "latest")
	assert.NoError(t, err)

	// List.
	_, err = executeCommand(imgCmd, "list", "withTag=true")
	assert.NoError(t, err)

	// Delete.
	_, err = executeCommand(imgCmd, "delete", "hello-world", "latest")
	assert.NoError(t, err)

}
