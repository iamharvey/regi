package action

import (
	"fmt"
	"github.com/iamharvey/regi/internal/pkg/cache"
	"github.com/iamharvey/regi/internal/pkg/rest"
	"github.com/iamharvey/regi/internal/pkg/settings"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"time"
)

type repoResult struct {
	repositories []string `json:"repositories"`
}

type tagResult struct {
	name string   `json:"name"`
	tags []string `json:"tags"`
}

func List(c *cli.Context) error {
	ca, err := cache.Load()
	if err != nil {
		return errors.Wrapf(err, "fail to load cache data")
	}

	registry, err := ca.Get(settings.CacheKeyRegistry)
	if err != nil {
		return errors.Wrapf(err, "fail to get `registry`")
	}

	fmt.Fprintf(c.App.Writer, "registry: %s\n", registry)

	cliConfig := &rest.ClientConfig{
		Host:    string(registry),
		APIPath: "v2/_catalog",
		ContentConfig: &rest.ContentConfig{
			ContentType: "application/json",
		},
		TLSClientConfig: nil,
		Timeout:         time.Second * 3,
	}
	client, err := rest.NewClient(cliConfig)
	if err != nil {
		return errors.Wrap(err, "fail to create new HTTP client")
	}
	// Make HTTP call to the remote API server.
	resp, err := client.
		Verb("GET").
		Do()
	if err != nil || resp == nil {
		return errors.Wrap(err, "fail to get response")
	}
	defer resp.Body.Close()

	withTags := c.Value("tags").(bool)

	rres, err := rest.DecodeResponse(resp)
	if err != nil {
		return errors.Wrap(err, "fail to decode response data for repositories")
	}

	for _, repo := range rres["repositories"].([]interface{}) {
		if withTags {
			cliConfig := &rest.ClientConfig{
				Host:    string(registry),
				APIPath: fmt.Sprintf("v2/%s/tags/list", repo),
				ContentConfig: &rest.ContentConfig{
					ContentType: "application/json",
				},
				TLSClientConfig: nil,
				Timeout:         time.Second * 3,
			}

			client, err := rest.NewClient(cliConfig)
			// Make HTTP call to the remote API server.
			resp, err := client.
				Verb("GET").
				Do()

			if err != nil || resp == nil {
				return errors.Wrap(err, "fail to get response")
			}

			tres, err := rest.DecodeResponse(resp)
			if err != nil {
				return errors.Wrap(err, "fail to decode response data for tags")
			}

			resp.Body.Close()

			fmt.Fprintf(c.App.Writer, "%s:	%s\n", tres["name"], tres["tags"])
		} else {
			fmt.Fprintf(c.App.Writer, "%s\n", repo)
		}
	}

	return nil
}
