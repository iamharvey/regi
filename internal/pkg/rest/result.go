package rest

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

// DecodeResponse decodes response body into a map.
func DecodeResponse(resp *http.Response) (map[string]interface{}, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errors.Wrap(err, "decoding response failed")
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.Wrap(err, "fail to encode to JSON string")
	}

	return m, nil
}
