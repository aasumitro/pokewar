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

type battleHTTPHandlerTestSuite struct {
	suite.Suite
	battles []*domain.Battle
}

func (suite *battleHTTPHandlerTestSuite) SetupSuite() {
	suite.battles = []*domain.Battle{
		{
			ID:        1,
			StartedAt: 1,
			EndedAt:   1,
			Players: []domain.Player{
				{
					ID:           1,
					Name:         "asd",
					BattleID:     1,
					MonsterID:    1,
					EliminatedAt: 1,
					AnnulledAt:   1,
					Rank:         1,
					Point:        1,
				},
			},
			Logs: []domain.Log{
				{
					ID:          1,
					BattleID:    1,
					Description: "asd",
					CreatedAt:   1,
				},
			},
		},
	}
}

func (suite *battleHTTPHandlerTestSuite) TestHandler_Fetch_ShouldSuccess() {
	svcMock := new(mocks.IPokewarService)
	svcMock.
		On("FetchBattles", mock.Anything).
		Return(suite.battles, nil).
		Once()
	svcMock.
		On("BattlesCount").
		Return(50).
		Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/battles?limit=1", nil)
	ctx.Request = req
	handler := httpHandler.BattleHTTPHandler{Svc: svcMock}
	handler.Fetch(ctx)
	var got utils.SuccessRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	assert.Equal(suite.T(), http.StatusOK, writer.Code)
	assert.Equal(suite.T(), http.StatusOK, got.Code)
	assert.Equal(suite.T(), http.StatusText(http.StatusOK), got.Status)
}

func (suite *battleHTTPHandlerTestSuite) TestHandler_Fetch_ShouldError() {
	svcMock := new(mocks.IPokewarService)
	svcMock.
		On("FetchBattles", mock.Anything).
		Return(nil, &utils.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "UNEXPECTED_ERROR",
		}).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/battles?limit=1", nil)
	ctx.Request = req
	handler := httpHandler.BattleHTTPHandler{Svc: svcMock}
	handler.Fetch(ctx)
	var got utils.ErrorRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	assert.Equal(suite.T(), http.StatusInternalServerError, writer.Code)
	assert.Equal(suite.T(), http.StatusInternalServerError, got.Code)
	assert.Equal(suite.T(), http.StatusText(http.StatusInternalServerError), got.Status)
	assert.Equal(suite.T(), "UNEXPECTED_ERROR", got.Data)
}

func TestBattleHTTPHandler(t *testing.T) {
	suite.Run(t, new(battleHTTPHandlerTestSuite))
}
