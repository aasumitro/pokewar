package httpclient_test

import (
	"context"
	"github.com/aasumitro/pokewar/pkg/httpclient"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestMakeRequest(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		endpoint    string
		timeout     time.Duration
		status      int
		expected    interface{}
		method      string
		retry       bool
		retryMax    int
		retryWait   time.Duration
		wantErr     bool
		setupServer func(*httptest.Server, string)
	}{
		{
			name:     "valid request",
			response: `{"foo": "bar"}`,
			status:   200,
			method:   http.MethodGet,
			expected: map[string]string{"foo": "bar"},
			wantErr:  false,
		},
		{
			method:   http.MethodGet,
			name:     "invalid endpoint",
			endpoint: "invalid endpoint",
			wantErr:  true,
		},
		{
			method:  http.MethodGet,
			name:    "short timeout",
			timeout: time.Nanosecond,
			wantErr: true,
		},
		{
			method:   http.MethodGet,
			name:     "nil endpoint",
			endpoint: "",
			wantErr:  true,
		},
		{
			method:   http.MethodGet,
			name:     "invalid response",
			response: "invalid response",
			wantErr:  true,
		},
		{
			method:   http.MethodGet,
			name:     "large response",
			response: string(make([]byte, 10485760)),
			wantErr:  true,
		},
		{
			method:   http.MethodGet,
			name:     "invalid JSON response",
			response: `{"foo": "bar`,
			wantErr:  true,
		},
		{
			method:   "lorem_ipsum",
			name:     "invalid method",
			response: `{"foo": "bar`,
			wantErr:  true,
		},
		{
			name:      "valid request with retry",
			response:  `{"foo": "bar"}`,
			status:    200,
			method:    http.MethodGet,
			expected:  map[string]string{"foo": "bar"},
			retry:     true,
			retryMax:  3,
			retryWait: 1 * time.Second,
			wantErr:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var server *httptest.Server
			if test.setupServer != nil {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					test.setupServer(server, test.response)
				}))
			} else {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(test.status)
					_, _ = w.Write([]byte(test.response))
				}))
			}
			defer server.Close()

			if test.endpoint == "" {
				test.endpoint = server.URL
			}

			c := &httpclient.HTTPClient{}
			c = c.NewClient(
				httpclient.Timeout(test.timeout),
				httpclient.Endpoint(test.endpoint),
				httpclient.Method(test.method),
				httpclient.Ctx(context.TODO()))
			var obj map[string]string
			err := c.MakeRequest(&obj)

			if (err != nil) != test.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && !reflect.DeepEqual(obj, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, obj)
			}
		})
	}
}
