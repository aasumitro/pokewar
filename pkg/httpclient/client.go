package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	HTTPClient struct {
		Ctx      context.Context
		Endpoint string
		Timeout  time.Duration
		Method   string
	}

	Option func(*HTTPClient)

	IHttpClient interface {
		NewClient(opts ...Option) *HTTPClient
		MakeRequest(obj interface{}) error
	}
)

func Ctx(ctx context.Context) Option {
	return func(httpclient *HTTPClient) {
		httpclient.Ctx = ctx
	}
}

func Endpoint(endpoint string) Option {
	return func(httpclient *HTTPClient) {
		httpclient.Endpoint = endpoint
	}
}

func Timeout(timeout time.Duration) Option {
	return func(httpclient *HTTPClient) {
		httpclient.Timeout = timeout
	}
}

func Method(method string) Option {
	return func(httpclient *HTTPClient) {
		httpclient.Method = method
	}
}

// NewClient - with default value and/or options
func (c *HTTPClient) NewClient(opts ...Option) *HTTPClient {
	httpclient := &HTTPClient{
		Timeout: 10 * time.Second,
		Method:  http.MethodGet,
	}

	for _, opt := range opts {
		opt(httpclient)
	}

	fmt.Println(httpclient.Timeout)
	fmt.Println(httpclient.Method)
	fmt.Println(httpclient.Endpoint)
	return httpclient
}

func (c *HTTPClient) MakeRequest(obj interface{}) error {
	req, err := http.NewRequestWithContext(c.Ctx, c.Method, c.Endpoint, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return err
	}
	data := buf.Bytes()

	return json.Unmarshal(data, &obj)
}
