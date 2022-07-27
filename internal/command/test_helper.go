package command

import (
	"bytes"
	"github.com/spf13/cobra"
)

// executeCommand executes command.
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

// executeCommandC does the actual job of executing command.
func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
