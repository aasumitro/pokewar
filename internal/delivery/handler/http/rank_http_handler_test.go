package http_test

import (
	"encoding/json"
	"github.com/aasumitro/pokewar/domain"
	httpHandler "github.com/aasumitro/pokewar/internal/delivery/handler/http"
	"github.com/aasumitro/pokewar/mocks"
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type rankHTTPHandlerTestSuite struct {
	suite.Suite
	ranks []*domain.Rank
}

func (suite *rankHTTPHandlerTestSuite) SetupSuite() {
	suite.ranks = []*domain.Rank{
		{
			ID:           1,
			OriginID:     1,
			Name:         "test",
			Avatar:       "test.png",
			T:            "['test']",
			Types:        []string{"test"},
			TotalBattles: 1,
			WinBattles:   1,
			LoseBattle:   0,
			Points:       5,
		},
	}
}

func (suite *rankHTTPHandlerTestSuite) TestHandler_Fetch_ShouldSuccess() {
	svcMock := new(mocks.IPokewarService)
	svcMock.
		On("FetchRanks", mock.Anything).
		Return(suite.ranks, nil).
		Once()
	svcMock.
		On("MonstersCount").
		Return(50).
		Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/monsters?limit=1", nil)
	ctx.Request = req
	handler := httpHandler.RankHTTPHandler{Svc: svcMock}
	handler.Fetch(ctx)
	var got utils.SuccessRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	assert.Equal(suite.T(), http.StatusOK, writer.Code)
	assert.Equal(suite.T(), http.StatusOK, got.Code)
	assert.Equal(suite.T(), http.StatusText(http.StatusOK), got.Status)
}

func (suite *rankHTTPHandlerTestSuite) TestHandler_Fetch_ShouldError() {
	svcMock := new(mocks.IPokewarService)
	svcMock.
		On("FetchRanks", mock.Anything).
		Return(nil, &utils.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "UNEXPECTED_ERROR",
		}).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/monsters?limit=1", nil)
	ctx.Request = req
	handler := httpHandler.RankHTTPHandler{Svc: svcMock}
	handler.Fetch(ctx)
	var got utils.ErrorRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	assert.Equal(suite.T(), http.StatusInternalServerError, writer.Code)
	assert.Equal(suite.T(), http.StatusInternalServerError, got.Code)
	assert.Equal(suite.T(), http.StatusText(http.StatusInternalServerError), got.Status)
	assert.Equal(suite.T(), "UNEXPECTED_ERROR", got.Data)
}

func TestRankHTTPHandler(t *testing.T) {
	suite.Run(t, new(rankHTTPHandlerTestSuite))
}
