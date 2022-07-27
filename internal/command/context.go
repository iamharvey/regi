package command

import (
	"fmt"
	"github.com/iamharvey/regi/internal/pkg/data"
	"github.com/iamharvey/regi/internal/pkg/io"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	// msgShortCtxCmd is the short description for context root command.
	msgShortCtxCmd = "Manage connection settings of multiple Docker registries."

	// msgExamplesCtxCmd is the example description for context root command.
	msgExamplesCtxCmd = `
  # Get all the contexts.
  regi context list

  # Set current context.
  regi context set <context-name>

  # Add new context
  regi context add -n=context1 -s=192.168.0.168:5000

  # Get context info
  regi context get context1
`

	// msgShortCtxListCmd is the short version description for `context list` command.
	msgShortCtxListCmd = "List all the contexts."

	// msgShortCtxSetCmd is the short version description for `context set` command.
	msgShortCtxSetCmd = "Set current context with context name."

	// msgShortCtxAddCmd is the short version description for `context add` command.
	msgShortCtxAddCmd = "Add a new context."

	// msgShortCtxGetCmd is the short version description for `context add` command.
	msgShortCtxGetCmd = "Get context info given context name."
)

// CmdContextOptions eases access to storage and console io.
type cmdContextOptions struct {
	*data.DB
	io.Streams
}

// NewCmdContextOptions returns a new Options for context command.
func NewCmdContextOptions(streams io.Streams) (*cmdContextOptions, error) {
	db, err := data.NewDB()
	if err != nil {
		return nil, err
	}
	return &cmdContextOptions{
		DB:      db,
		Streams: streams,
	}, nil
}

// NewCmdContext creates a Context command.
func NewCmdContext(streams io.Streams) *cobra.Command {
	o, err := NewCmdContextOptions(streams)
	if err != nil {
		streams.ErrOut.Write([]byte(err.Error()))
	}

	// Context root command.
	cmd := &cobra.Command{
		Use:                   "context",
		Aliases:               []string{"c", "ctx"},
		DisableFlagsInUseLine: true,
		Short:                 msgShortCtxCmd,
		Example:               msgExamplesCtxCmd,
	}

	// List all the contexts.
	listCmd := &cobra.Command{
		Use:                   "list",
		Aliases:               []string{"ls"},
		DisableFlagsInUseLine: true,
		Short:                 msgShortCtxListCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.listCmdRun())
		},
	}

	// Set current context.
	setCmd := &cobra.Command{
		Use:                   "set",
		DisableFlagsInUseLine: true,
		Short:                 msgShortCtxSetCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.setCmdRun(cmd, args))
		},
	}

	// Add new context.
	addCmd := &cobra.Command{
		Use:                   "add",
		DisableFlagsInUseLine: true,
		Short:                 msgShortCtxAddCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.addCmdRun(cmd))
		},
	}

	// Get context info.
	getCmd := &cobra.Command{
		Use:                   "get",
		DisableFlagsInUseLine: true,
		Short:                 msgShortCtxGetCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.getCmdRun(cmd, args))
		},
	}

	// Add commands.
	cmd.AddCommand(listCmd)
	cmd.AddCommand(setCmd)
	cmd.AddCommand(addCmd)
	cmd.AddCommand(getCmd)

	// Add flags.
	addCmd.Flags().StringP("name", "n", "", "context name (required)")
	addCmd.Flags().StringP("server", "s", "", "registry server address (required)")
	addCmd.Flags().BoolP("verify", "v", false, "insecure skip TLS verify, default is false")
	addCmd.Flags().StringP("user", "u", "", "registry username")
	addCmd.Flags().StringP("password", "p", "", "registry password")

	// Set required options.
	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("server")

	return cmd
}

// listCmdRun lists all the registries.
func (o *cmdContextOptions) listCmdRun() error {
	// Get current registry.
	current, err := o.Current()
	if err != nil {
		return err
	}

	// Get all cached registries.
	registries, err := o.DB.List()
	if err != nil {
		return err
	}

	// Display all cached registries.
	o.Out.Write([]byte("\nContexts:\n"))
	if len(registries) == 0 {
		o.Out.Write([]byte("\nno registries added\n"))
		return nil
	}

	for _, v := range registries {
		hl := ""
		// Point out the current context.
		if v.Name == current.Name {
			hl = " <---"
		}
		o.Out.Write([]byte(fmt.Sprintf("- %s%s\n", v.Name, hl)))
	}

	return nil
}

// setCmdRun sets current registry.
func (o *cmdContextOptions) setCmdRun(cmd *cobra.Command, args []string) error {
	tips := fmt.Sprintf(">> tips：please use '%s -h' to get for information about the command.", cmd.CommandPath())

	if len(args) == 0 {
		return errors.Errorf("context name is missing\n%s", tips)
	}

	name := args[0]
	// Check against cache registries. If the given name does not match any entries, return an error.
	reg, err := o.Get(name)
	if err != nil {
		return err
	}

	if reg == nil {
		return errors.New("context name does not match cached ones")
	}

	// Otherwise, set current context to the given name.
	err = o.Set(name)
	if err != nil {
		return err
	}

	o.Out.Write([]byte(fmt.Sprintf("Context switched. Current context: %s\n", args[0])))

	return nil
}

// getCmdRun get info about current registry.
func (o *cmdContextOptions) getCmdRun(cmd *cobra.Command, args []string) error {
	tips := fmt.Sprintf(">> tips：please use '%s -h' to get for information about the command.", cmd.CommandPath())

	ctxName := args[0]
	if len(ctxName) == 0 {
		return errors.Errorf("context name is not specified, %s", tips)
	}

	// Check against cache registries. If the given name does not match any entries, return an error.
	reg, err := o.DB.Get(ctxName)
	if err != nil {
		return err
	}

	if reg == nil {
		return errors.Errorf("context %q not found", ctxName)
	}

	// Display info.
	o.Out.Write([]byte("\n" + fmt.Sprintf(`Context %q:
- server: %s
- insecure skip TLS verify: %v
- user: %s
- password: ***
`, reg.Name, reg.Server, reg.InsecureSkipTLSVerify, reg.User)))
	return nil
}

// addCmdRun add new registry.
func (o *cmdContextOptions) addCmdRun(cmd *cobra.Command) error {
	// Context name.
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	// Registry server address.
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	// Insecure skip TLS verify.
	verify, err := cmd.Flags().GetBool("verify")
	if err != nil {
		return err
	}

	// Registry username.
	user, err := cmd.Flags().GetString("user")
	if err != nil {
		return err
	}

	// Registry password.
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return err
	}

	// Add new context.
	ok, err := o.DB.Add(name, server, user, password, verify)
	if err != nil {
		return err
	}

	if !ok {
		return errors.Errorf("fail to add context, duplicated entry for context %q is not allowed", name)
	}

	o.Out.Write([]byte(fmt.Sprintf(`Context added:
- name: %s
- server: %s
- insecure skip TLS verify: %v
- user: %s
- password: ***
`, name, server, verify, user)))
	return nil
}
