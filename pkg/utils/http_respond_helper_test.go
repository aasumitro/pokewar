package utils_test

import (
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHttpRespond(t *testing.T) {
	tests := []struct {
		name          string
		code          int
		data          interface{}
		args          []any
		expected      interface{}
		expectedError bool
	}{
		{
			name:     "success with no pagination",
			code:     http.StatusOK,
			data:     []string{"foo", "bar"},
			args:     nil,
			expected: utils.SuccessRespond{Code: http.StatusOK, Status: "OK", Data: []string{"foo", "bar"}},
		},
		{
			name: "success with pagination",
			code: http.StatusOK,
			data: []string{"foo", "bar"},
			args: []any{
				2,
				1,
				utils.Paging{Url: "http://example.com/next", Path: "/next"},
				utils.Paging{Url: "http://example.com/prev", Path: "/prev"},
			},
			expected: utils.SuccessRespondWithPagination{
				Code:     http.StatusOK,
				Status:   "OK",
				Total:    2,
				Current:  1,
				Next:     utils.Paging{Url: "http://example.com/next", Path: "/next"},
				Previous: utils.Paging{Url: "http://example.com/prev", Path: "/prev"},
				Data:     []string{"foo", "bar"},
			},
		},
		{
			name:     "error with data",
			code:     http.StatusBadRequest,
			data:     "invalid request",
			expected: utils.ErrorRespond{Code: http.StatusBadRequest, Status: "Bad Request", Data: "invalid request"},
		},
		{
			name:          "error with no data",
			code:          http.StatusBadRequest,
			expected:      utils.ErrorRespond{Code: http.StatusBadRequest, Status: "Bad Request", Data: "something went wrong with the request"},
			expectedError: false,
		},
		{
			name:          "error with no data and server error code",
			code:          http.StatusInternalServerError,
			expected:      utils.ErrorRespond{Code: http.StatusInternalServerError, Status: "Internal Server Error", Data: "something went wrong with the server"},
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(writer)
			utils.NewHttpRespond(c, test.code, test.data, test.args...)

			if test.expectedError {
				t.Error("expected error but got none")
				return
			}
		})
	}
}
