// This package contain main structure Client, for working with all methods. Client contain substructures for short way
// access to methods
package httpClient

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

// Error response for all node requests
type ErrorResponse struct {
	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// base client included in all clients
type baseClient struct {
	baseAddress string
}

// create base client
func newBaseClient(baseAddress string) *baseClient {
	return &baseClient{
		baseAddress: baseAddress,
	}
}

// Read response of node request, if error code is not 200, return formatted error
func (bk *baseClient) ReadResponse(resp *resty.Response) ([]byte, error) {
	body := resp.Body()
	if resp.StatusCode() != 200 {
		errResponse := ErrorResponse{}
		if err := json.Unmarshal(body, &errResponse); err != nil {
			return nil, errors.Errorf("cannot read error response: %s", string(body))
		}
		return nil, errors.Errorf(errResponse.ErrorText)
	}
	return body, nil
}
