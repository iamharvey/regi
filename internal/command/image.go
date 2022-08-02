package command

import (
	"fmt"
	"github.com/iamharvey/regi/internal/pkg/data"
	rio "github.com/iamharvey/regi/internal/pkg/io"
	"github.com/iamharvey/regi/internal/pkg/rest"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
	"time"
)

const (
	// msgShortImageCmd is the short version description for image root command.
	msgShortImageCmd = "Pull, push, delete and list images over Docker registry"

	// msgShortImageListCmd is the short version description for 'image list' command.
	msgShortImgListCmd = "ListContexts images on current registry."

	// msgShortImgPullCmd is the short version description for 'image pull' command.
	msgShortImgPullCmd = "Pull image from current registry."

	// msgShortImgPushCmd is the short version description for 'image push' command.
	msgShortImgPushCmd = "Push image to current registry."

	// msgShortImgDelCmd is the short version description for 'image delete' command.
	msgShortImgDelCmd = "DeleteContext image from current registry."

	// termDelSuccess is the key term for successful deletion operation.
	termDelSuccess = "202 accepted"

	// termDelFailure is the key term for failed deletion operation.
	// termDelFailure = "404 not found"
)

// cmdImageOptions eases access to storage and console io.
type cmdImageOptions struct {
	*data.DB
	rio.Streams
}

// NewCmdImageOptions returns a new Options for image command.
func NewCmdImageOptions(streams rio.Streams) (*cmdImageOptions, error) {
	db, err := data.NewDB()
	if err != nil {
		return nil, err
	}
	return &cmdImageOptions{
		DB:      db,
		Streams: streams,
	}, nil
}

// NewCmdImage creates an image command.
func NewCmdImage(streams rio.Streams) *cobra.Command {
	o, err := NewCmdImageOptions(streams)
	if err != nil {
		streams.ErrOut.Write([]byte(err.Error()))
	}

	// Context root command.
	cmd := &cobra.Command{
		Use:                   "image",
		Aliases:               []string{"i", "im", "img"},
		DisableFlagsInUseLine: true,
		Short:                 msgShortImageCmd,
	}

	// ListContexts all the contexts.
	listCmd := &cobra.Command{
		Use:                   "list",
		Aliases:               []string{"ls"},
		DisableFlagsInUseLine: true,
		Short:                 msgShortImgListCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.listCmdRun(cmd))
		},
	}

	// Pull image from remote registry.
	pullCmd := &cobra.Command{
		Use:                   "pull",
		DisableFlagsInUseLine: true,
		Short:                 msgShortImgPullCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.pullCmdRun(args))
		},
	}

	// Push image to remote registry.
	pushCmd := &cobra.Command{
		Use:                   "push",
		DisableFlagsInUseLine: true,
		Short:                 msgShortImgPushCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.pushCmdRun(args))
		},
	}

	// DeleteContext image from remote registry.
	delCmd := &cobra.Command{
		Use:                   "delete",
		Aliases:               []string{"d", "del"},
		DisableFlagsInUseLine: true,
		Short:                 msgShortImgDelCmd,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.delCmdRun(args))
		},
	}

	cmd.AddCommand(listCmd)
	cmd.AddCommand(pullCmd)
	cmd.AddCommand(pushCmd)
	cmd.AddCommand(delCmd)
	listCmd.Flags().BoolP("withTag", "t", true, "show tags")

	return cmd
}

// imageCmdRun lists all the images with/without tags for current registry.
func (o *cmdImageOptions) listCmdRun(cmd *cobra.Command) error {
	showTags, err := cmd.Flags().GetBool("withTag")
	if err != nil {
		return err
	}

	// GetContext current registry.
	current, err := o.CurrentContext()
	if err != nil {
		return err
	}

	if current == nil {
		return errors.New("context is not set, please set current context with 'regi ctx set <name>' first")
	}

	// Create REST client for accessing APIs.
	cliConfig := &rest.ClientConfig{
		Host:    current.Server,
		APIPath: "v2/_catalog",
		ContentConfig: &rest.ContentConfig{
			ContentType: "application/json",
		},
		TLSClientConfig: nil,
		Timeout:         time.Second * 3,
	}

	client, err := rest.NewClient(cliConfig)
	if err != nil {
		return err
	}

	// Execute API call.
	resp, err := client.
		Verb("GET").
		Do()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response body.
	rres, err := rest.DecodeResponse(resp)
	if err != nil {
		return err
	}

	// Display all the images.
	fmt.Println("\nImages:")
	for _, repo := range rres["repositories"].([]interface{}) {
		fmt.Printf("- %s ", repo)

		// Query image tags.
		if showTags {
			cliConfig := &rest.ClientConfig{
				Host:    current.Server,
				APIPath: fmt.Sprintf("v2/%s/tags/list", repo),
				ContentConfig: &rest.ContentConfig{
					ContentType: "application/json",
				},
				TLSClientConfig: nil,
				Timeout:         time.Second * 3,
			}

			client, err := rest.NewClient(cliConfig)
			if err != nil {
				return err
			}

			resp, err := client.
				Verb("GET").
				Do()
			if err != nil {
				return err
			}

			tres, err := rest.DecodeResponse(resp)
			if err != nil {
				return err
			}

			resp.Body.Close()

			fmt.Printf(" %s", tres["tags"])
		}

		fmt.Println()
	}

	return nil
}

// pullCmdRun pull image from remote registry.
func (o *cmdImageOptions) pullCmdRun(args []string) error {
	// Verify arguments.
	if len(args) == 0 {
		return errors.New("image and tag is not specified")
	}

	if len(args) == 1 {
		return errors.New("image tag is not specified")
	}

	name := args[0]
	tag := args[1]

	// GetContext current registry.
	current, err := o.CurrentContext()
	if err != nil {
		return err
	}

	// Remove prefix for registry server.
	server := strings.TrimPrefix(strings.TrimPrefix(current.Server, "https://"), "http://")

	// Execute command.
	c := exec.Command(
		"docker",
		"pull",
		fmt.Sprintf("%s/%s:%s", server, name, tag),
	)

	out, err := c.CombinedOutput()
	if err != nil {
		return err
	}

	o.Streams.Out.Write([]byte(fmt.Sprintf("%s", out)))

	return nil
}

// pushCmdRun tag and push image to remote registry.
func (o *cmdImageOptions) pushCmdRun(args []string) error {
	// Verify arguments.
	if len(args) == 0 {
		return errors.New("image and tag is not specified")
	}

	if len(args) == 1 {
		return errors.New("image tag is not specified")
	}

	name := args[0]
	tag := args[1]

	// GetContext current registry.
	current, err := o.CurrentContext()
	if err != nil {
		return err
	}

	// Remove prefix for registry server.
	server := strings.TrimPrefix(strings.TrimPrefix(current.Server, "https://"), "http://")

	// Tag image.
	c := exec.Command(
		"docker",
		"image",
		"tag",
		fmt.Sprintf("%s:%s", name, tag),
		fmt.Sprintf("%s/%s:%s", server, name, tag))

	out, err := c.CombinedOutput()
	if err != nil {
		return err
	}

	o.Streams.Out.Write([]byte(fmt.Sprintf("%s", out)))

	// Push image.
	c = exec.Command(
		"docker",
		"image",
		"push",
		fmt.Sprintf("%s/%s:%s", server, name, tag))

	out, err = c.CombinedOutput()
	if err != nil {
		return err
	}

	o.Streams.Out.Write([]byte(fmt.Sprintf("\n%s", out)))

	return nil
}

// delCmdRun delete image from remote registry.
func (o *cmdImageOptions) delCmdRun(args []string) error {
	// Verify arguments.
	if len(args) == 0 {
		return errors.New("image and tag is not specified")
	}

	if len(args) == 1 {
		return errors.New("image tag is not specified")
	}

	name := args[0]
	tag := args[1]

	// GetContext current registry.
	current, err := o.CurrentContext()
	if err != nil {
		return err
	}

	// First, get the manifest info with desired tag.
	cliConfig := &rest.ClientConfig{
		Host:    current.Server,
		APIPath: fmt.Sprintf("v2/%s/manifests/%s", name, tag),
		ContentConfig: &rest.ContentConfig{
			AcceptContentTypes: "application/vnd.docker.distribution.manifest.v2+json",
		},
		TLSClientConfig: nil,
		Timeout:         time.Second * 3,
	}

	client, err := rest.NewClient(cliConfig)
	if err != nil {
		return err
	}

	// Execute API call.
	resp, err := client.
		Verb("GET").
		Do()
	if err != nil {
		return err
	}

	// Second, we delete that image using the obtained manifest digest.
	// In most cases, digest can be obtained by resp.Header.GetContext("Docker-Content-Digest").
	digest := resp.Header.Get("Docker-Content-Digest")
	cliConfig = &rest.ClientConfig{
		Host:    current.Server,
		APIPath: fmt.Sprintf("v2/%s/manifests/%s", name, digest),
		ContentConfig: &rest.ContentConfig{
			AcceptContentTypes: "application/vnd.docker.distribution.manifest.v2+json",
		},
		TLSClientConfig: nil,
		Timeout:         time.Second * 3,
	}

	client, err = rest.NewClient(cliConfig)
	if err != nil {
		return err
	}

	// Execute API call.
	resp, err = client.
		Verb("DELETE").
		Do()
	if err != nil {
		return err
	}
	resp.Body.Close()

	if strings.Contains(strings.ToLower(fmt.Sprintf("%v", resp)), termDelSuccess) {
		o.Streams.Out.Write([]byte(fmt.Sprintf("image %s:%s is deleted\n", name, tag)))
		return nil
	} else {
		return errors.Errorf(
			"unable to perform deletion on %s:%s, the manifest is either not found or has been deleted already\n",
			name, tag)
	}
}
