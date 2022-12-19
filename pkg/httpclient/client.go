package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	Endpoint string
	Timeout  time.Duration
	Method   string
	// TODO retry mechanism
	// retry? && retryCount!
}

type IHttpClient interface {
	MakeRequest(obj interface{}) error
}

func (c *HttpClient) MakeRequest(obj interface{}) error {
	req, err := http.NewRequest(c.Method, c.Endpoint, nil)
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
