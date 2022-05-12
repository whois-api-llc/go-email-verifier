package emailverifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// EvapiService is an interface for Email Verification API
type EvapiService interface {
	// Get returns parsed Email Verification API response
	Get(ctx context.Context, emailAddress string, opts ...Option) (*EvapiResponse, *Response, error)

	// GetRaw returns raw Email Verification API response as Response struct with Body saved as a byte slice
	GetRaw(ctx context.Context, emailAddress string, opts ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice
type Response struct {
	*http.Response

	//Body is the byte slice representation of http.Response Body
	Body []byte
}

// emailVerifierServiceOp is the type implementing the EvapiService interface
type emailVerifierServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ EvapiService = &emailVerifierServiceOp{}

// newRequest creates the API request with default parameters and the specified apiKey
func (service *emailVerifierServiceOp) newRequest() (*http.Request, error) {

	req, err := service.client.NewRequest(http.MethodGet, service.baseURL, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// apiResponse is used for parsing Email Verification API response as a model instance
type apiResponse struct {
	EvapiResponse
	ErrorMessage *ErrorMessage `json:"ErrorMessage"`
}

// request returns intermediate EVAPI response for further actions
func (service *emailVerifierServiceOp) request(ctx context.Context, emailAddress string, opts ...Option) (*Response, error) {
	if emailAddress == "" {
		return nil, &ArgError{"emailAddress", "cannot be empty"}
	}

	req, err := service.newRequest()
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("emailAddress", emailAddress)

	for _, opt := range opts {
		opt(q)
	}

	req.URL.RawQuery = q.Encode()

	var b bytes.Buffer
	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw Email Verification API response
func parse(raw []byte) (*apiResponse, error) {

	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// Get returns parsed Email Verification API response
func (service emailVerifierServiceOp) Get(
	ctx context.Context,
	emailAddress string,
	opts ...Option,
) (evapiResponse *EvapiResponse, resp *Response, err error) {

	optsJson := make([]Option, 0, len(opts)+1)
	optsJson = append(optsJson, opts...)
	optsJson = append(optsJson, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, emailAddress, optsJson...)
	if err != nil {
		return nil, resp, err
	}

	evapiResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if evapiResp.ErrorMessage != nil {
		return nil, nil, ErrorMessage{
			evapiResp.ErrorMessage.Message,
		}
	}

	return &evapiResp.EvapiResponse, resp, nil
}

// GetRaw returns raw Email Verification API response as Response struct with Body saved as a byte slice
func (service emailVerifierServiceOp) GetRaw(
	ctx context.Context,
	name string,
	opts ...Option,
) (resp *Response, err error) {

	resp, err = service.request(ctx, name, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
