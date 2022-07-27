package rest

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cliConfig := &ClientConfig{
		Host:    "https://api.github.com",
		APIPath: "repos/iamharvey/regi",
		ContentConfig: &ContentConfig{
			ContentType: "application/json",
		},
		TLSClientConfig: nil,
		Timeout:         time.Second * 3,
	}
	assert.NotNil(t, cliConfig)

	client, err := NewClient(cliConfig)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Examine response
	resp, err := client.
		Verb("GET").
		Do()

	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	rres, err := DecodeResponse(resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, rres)

	assert.Equal(t,
		"https://github.com/iamharvey/regi.git",
		rres["clone_url"])
	assert.Equal(t,
		"regi is a CLI tool for managing your accessibility to multiple Docker registries.",
		rres["description"])
}
