package rest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRequest(t *testing.T) {
	cases := []struct {
		name      string
		verb      string
		body      map[string]string
		selectors map[string]string
	}{
		{
			name: "test GET",
			verb: "GET",
			selectors: map[string]string{
				"token":        "1234567",
				"redirectFrom": "api.server",
			},
		},
		{
			name: "test POST",
			verb: "POST",
			body: map[string]string{
				"name":  "Alice",
				"email": "alice@unit.test",
			},
		},
	}

	for _, c := range cases {
		cliConfig := &ClientConfig{
			Host:    "https://api.github.com",
			APIPath: "repos/iamharvey/regi",
		}
		assert.NotNil(t, cliConfig)

		client, err := NewClient(cliConfig)
		assert.NoError(t, err)
		assert.NotNil(t, client)

		r := NewRequest(client)
		r = r.
			Verb(c.verb).
			Body(c.body).
			Selectors(c.selectors)

		assert.NotNil(t, r)
		assert.Equal(t, c.verb, r.verb)
		assert.Equal(t, c.body, r.body)
		assert.Equal(t, c.selectors, r.selectors)
	}

}
