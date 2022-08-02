package command

import (
	"fmt"
	"github.com/iamharvey/regi/internal/pkg/data"
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os/exec"
)

const (
	// msgShortLoginCmd is the short version description for login root command.
	msgShortLoginCmd = "Login to current Docker registry."
)

// cmdLoginOptions eases access to storage and console io.
type cmdLoginOptions struct {
	*data.DB
	io.Streams
}

// NewCmdLoginOptions returns a new Options for login command.
func NewCmdLoginOptions(streams io.Streams) (*cmdLoginOptions, error) {
	db, err := data.NewDB()
	if err != nil {
		return nil, err
	}
	return &cmdLoginOptions{
		DB:      db,
		Streams: streams,
	}, nil
}

// NewCmdLogin creates a login command.
func NewCmdLogin(streams io.Streams) *cobra.Command {
	o, err := NewCmdLoginOptions(streams)
	if err != nil {
		streams.ErrOut.Write([]byte(err.Error()))
	}

	// Context root command.
	cmd := &cobra.Command{
		Use:                   "login",
		DisableFlagsInUseLine: true,
		Short:                 msgShortLoginCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.loginCmdRun())
		},
	}

	return cmd
}

// loginCmdRun lists all the registries.
func (o *cmdLoginOptions) loginCmdRun() error {
	reg, err := o.CurrentContext()
	if err != nil {
		return err
	}

	if reg == nil {
		return errors.New("context is not set, please set current context with 'regi ctx set <name>' first")
	}

	o.Out.Write([]byte(fmt.Sprintf("\nConnecting Docker registry with current context [%s](%s)",
		reg.Name, reg.Server)))

	c := exec.Command(
		"docker",
		"login",
		"-u", reg.User,
		"-p", reg.Password,
		reg.Server)

	out, err := c.CombinedOutput()
	if err != nil {
		return err
	}

	o.Streams.Out.Write([]byte(fmt.Sprintf("\n%s", out)))

	return nil
}
