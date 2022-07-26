package command

import (
	"flag"
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/spf13/cobra"
	"os"
)

const (
	// CmdShortMsg is the short description for root cmd.
	msgShort = "Regi - A command line tool for communicating with your Docker registries."

	// CmdLongMsg is the longer description for root cmd.
	msgLong = `Regi is a command line tool for communicating with your Docker registries. 

It allows you to query, push, pull and delete images from your registries. 
You can use Regi to manage connections with multiple Docker registries.

For more information about Regi, please visit:
	 http://github.com/iamharvey/regi
`
)

// NewRegiCommand creates Regi root command.
func NewRegiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "regi",
		Short: msgShort,
		Long:  msgLong,
		Run:   runHelp,
	}

	// Initialize an io stream with standard io reader and writers.
	streams := io.Streams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	// Add sub commands.
	cmd.AddCommand(
		NewCmdContext(streams),
		NewCmdLogin(streams),
		NewCmdImage(streams),
	)

	// Add go flag set.
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return cmd
}

// runHelp prints help msg.
func runHelp(cmd *cobra.Command, _ []string) {
	cmd.Help()
}
