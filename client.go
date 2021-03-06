package emailverifier

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	libraryVersion = "0.1.0"
	userAgent      = "go-email-verifier/" + libraryVersion
	mediaType      = "application/json"
)

// defaultEvapiURL is the default Email Verification API URL
const defaultEvapiURL = `https://emailverification.whoisxmlapi.com/api/v2`

// ClientParams is used to create Client. None of parameters are mandatory and
// leaving this struct empty works just fine for most cases
type ClientParams struct {
	// HTTPClient is the client used to access API endpoint
	// If it's nil then value API client uses http.DefaultClient
	HTTPClient *http.Client

	// EvapiBaseURL is the endpoint for 'Email Verification API' service
	EvapiBaseURL *url.URL
}

// NewBasicClient creates Client with recommended parameters
func NewBasicClient(APIKey string) *Client {
	return NewClient(APIKey, ClientParams{})
}

// NewClient creates Client with specified parameters
func NewClient(apiKey string, params ClientParams) *Client {

	var err error

	evapiBaseURL := params.EvapiBaseURL
	if evapiBaseURL == nil {
		evapiBaseURL, err = url.Parse(defaultEvapiURL)
		if err != nil {
			panic(err)
		}
	}

	httpClient := http.DefaultClient
	if params.HTTPClient != nil {
		httpClient = params.HTTPClient
	}

	client := &Client{
		client:    httpClient,
		userAgent: userAgent,
		apiKey:    apiKey,
	}

	client.EvapiService = &emailVerifierServiceOp{client: client, baseURL: evapiBaseURL}

	return client
}

// Client is the client for Email Verification API services
type Client struct {
	client *http.Client

	userAgent string
	apiKey    string

	// EmailVerifierService is an interface for Email Verification API
	EvapiService
}

// NewRequest creates a basic API request
func (c *Client) NewRequest(method string, u *url.URL, body io.Reader) (*http.Request, error) {

	var err error
	var req *http.Request

	req, err = http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

// Do sends the API request and returns the API response
func (c *Client) Do(ctx context.Context, req *http.Request, v io.Writer) (response *http.Response, err error) {

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot execute request: %w", err)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil && rerr != nil {
			err = fmt.Errorf("cannot close response: %w", rerr)
		}
	}()

	_, err = io.Copy(v, resp.Body)
	if err != nil {
		return resp, fmt.Errorf("cannot read response: %w", err)
	}

	return resp, err
}

// ErrorResponse is returned when the response status code is not 2xx
type ErrorResponse struct {
	Response *http.Response
	Message  string
}

// Error returns error message as a string
func (e ErrorResponse) Error() string {
	if e.Message != "" {
		return "API failed with status code: " + strconv.Itoa(e.Response.StatusCode) + " (" + e.Message + ")"
	}
	return "API failed with status code: " + strconv.Itoa(e.Response.StatusCode)
}

// checkResponse checks if the response status code is not 2xx
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	var errorResponse = ErrorResponse{
		Response: r,
	}

	return errorResponse
}
