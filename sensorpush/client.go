package sensorpush

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL   = "https://api.sensorpush.com/api/v1/"
	defaultUserAgent = "go-sensorpush v0.1.0"

	mediaTypeJSON = "application/json"
)

type Client struct {
	c *http.Client // The underlying HTTP client

	BaseURL   *url.URL
	UserAgent string

	accessToken AccessToken

	// Services
	Auth   *AuthService
	Sample *SampleService
	Sensor *SensorService
	Status *StatusService
}

type service struct {
	c *Client
}

type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

// NewClient returns a new SensorPush API client. If a nil httpClient
// is provided a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		c:         httpClient,
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
	}

	// Services
	c.Auth = &AuthService{c}
	c.Sample = &SampleService{c}
	c.Sensor = &SensorService{c}
	c.Status = &StatusService{c}

	return c
}

// NewBaseRequest makes an unauthenticated request
func (c *Client) NewBaseRequest(ctx context.Context, method, urlStr string, body any) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have trailing slash, but %q does not", c.BaseURL)
	}

	if strings.HasPrefix(urlStr, "/") {
		return nil, fmt.Errorf("urlStr must not have trailing slash, but %q does", urlStr)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", mediaTypeJSON)
	}

	req.Header.Set("Accept", mediaTypeJSON)

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// NewRequest makes an authenticated request using the access token
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body any) (*http.Request, error) {
	if c.accessToken == "" {
		return nil, fmt.Errorf("access token required")
	}
	req, err := c.NewBaseRequest(ctx, method, urlStr, body)
	if err != nil {
		return req, err
	}
	req.Header.Set("Authorization", string(c.accessToken))
	return req, err
}

func (c *Client) UseAccessToken(tok AccessToken) {
	c.accessToken = tok
}

func (c *Client) BareDo(req *http.Request) (*Response, error) {

	rawResp, err := c.c.Do(req)
	if err != nil {
		defer rawResp.Body.Close()

		return nil, err
	}

	resp := newResponse(rawResp)
	err = CheckResponse(rawResp)
	return resp, err
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode string `json:"statusCode"`
	Type       string `json:"type"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("type=%s, msg=\"%s\"", e.Type, e.Message)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	defer r.Body.Close()
	errResp := ErrorResponse{}
	err := json.NewDecoder(r.Body).Decode(&errResp)
	if err != nil {
		return err
	}
	return errResp
}

func (c *Client) Do(req *http.Request, v any) (*Response, error) {
	var err error

	resp, err := c.BareDo(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF {
			err = nil // Ignore EOF errors caused by empty response
		}
	}
	return resp, err
}
