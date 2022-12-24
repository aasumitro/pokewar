package rest_test

import (
	"fmt"
	"github.com/aasumitro/pokewar/internal/repository/rest"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/aasumitro/pokewar/pkg/httpclient"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type pokeapiRESTRepositoryTestSuite struct {
	suite.Suite
	// mock func(url string, v interface{}) error
}

func (suite *pokeapiRESTRepositoryTestSuite) SetupSuite() {

}

func (suite *pokeapiRESTRepositoryTestSuite) TestRepository_Pokemon_() {
	viper.SetConfigFile("../../../.example.env")
	appconfigs.LoadEnv()
	appconfigs.Instance.PokeapiUrl = "http://example.com"

	testCases := []struct {
		offset      int
		limit       int
		response    string
		expectedErr error
	}{
		{
			offset:      0,
			limit:       10,
			response:    `{"results": [{"url": "invalid url"}]}`,
			expectedErr: fmt.Errorf("invalid character '<' looking for beginning of value"),
		},
	}

	for _, tc := range testCases {
		suite.T().Run("", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tc.response))
			}))
			defer server.Close()

			repo := rest.NewPokeapiRESTRepository()
			client := httpclient.HttpClient{}
			client.Endpoint = server.URL
			client.Timeout = 10 * time.Second
			client.Method = http.MethodGet
			monsters, err := repo.Pokemon(tc.offset, tc.limit)

			if tc.expectedErr != nil {
				if err == nil {
					t.Errorf("expected error %q, got nil", tc.expectedErr)
				} else if err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error %q, got %q", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected nil, got error %q", err)
				} else if len(monsters) != 1 {
					t.Errorf("expected 1 monster, got %d", len(monsters))
				}
			}
		})
	}
}

// ============
// TODO: HERE
// ============

func TestPokeapiRESTRepository(t *testing.T) {
	suite.Run(t, new(pokeapiRESTRepositoryTestSuite))
}
