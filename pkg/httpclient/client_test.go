package httpclient_test

import (
	"fmt"
	"github.com/aasumitro/pokewar/pkg/httpclient"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestMakeRequest(t *testing.T) {
	tests := []struct {
		name     string
		response string
		endpoint string
		timeout  time.Duration
		status   int
		expected interface{}
		method   string
		wantErr  bool
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.response))
			}))
			defer server.Close()

			if test.endpoint == "" {
				test.endpoint = server.URL
			}

			c := &httpclient.HttpClient{
				Endpoint: test.endpoint,
				Timeout:  test.timeout,
				Method:   test.method,
			}
			var obj map[string]string
			err := c.MakeRequest(&obj)

			if err != nil {
				fmt.Println(test.name, err)
			}
			if (err != nil) != test.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && !reflect.DeepEqual(obj, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, obj)
			}
		})
	}
}
