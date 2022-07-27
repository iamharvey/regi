package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	debug                  = true
	defaultPostContentType = "application/json"
)

// Request allows for building up a request to a server in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type Request struct {
	client    *Client
	verb      string
	url       string
	body      map[string]string
	selectors map[string]string
}

// NewRequest returns an new Request.
func NewRequest(c *Client) *Request {
	return &Request{
		client: c,
		url:    fmt.Sprintf("%s%s", c.base, c.versionedAPIPath),
	}
}

// Verb receives a verb that indicates which HTTP method should be invoked.
func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

// Selectors receives selectors that will be used to make query string.
func (r *Request) Selectors(selectors map[string]string) *Request {
	r.selectors = selectors
	return r
}

// Body receives body that will be used to make POST call.
func (r *Request) Body(body map[string]string) *Request {
	r.body = body
	return r
}

// Do does the real dirty job.
func (r *Request) Do() (*http.Response, error) {
	var (
		body io.Reader
		err  error
	)

	if r.selectors != nil {
		r.url = fmt.Sprintf("%s%s", r.url, r.makeQueryStrings())
	}

	if r.body != nil {
		body, err = r.makePostBody()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(r.verb, r.url, body)

	contentType := r.client.contentConfig.ContentType
	acceptType := r.client.contentConfig.AcceptContentTypes
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}

	if len(acceptType) > 0 {
		req.Header.Set("Accept", acceptType)
	}

	if err != nil {
		return nil, err
	}
	return r.client.Do(req)

}

func (r *Request) makeQueryStrings() string {
	q := ""
	first := true
	for k, v := range r.selectors {
		if first {
			q += fmt.Sprintf("?%s=%s", k, v)
			first = false
		} else {
			q += fmt.Sprintf("&%s=%s", k, v)
		}
	}
	return q
}

func (r *Request) makePostBody() (io.Reader, error) {
	postBody, err := json.Marshal(r.body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(postBody), nil
}
