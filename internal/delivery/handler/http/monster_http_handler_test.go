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

type monsterHTTPHandlerTestSuite struct {
	suite.Suite
	monsters []*domain.Monster
}

func (suite *monsterHTTPHandlerTestSuite) SetupSuite() {
	suite.monsters = []*domain.Monster{
		{
			ID:       1,
			OriginID: 2,
			Name:     "test",
			BaseExp:  1,
			Height:   1,
			Weight:   1,
			Avatar:   "test.png",
			Types:    []string{"test", "test"},
			Stats:    []domain.Stat{{BaseStat: 1, Name: "test"}},
			Skills:   []*domain.Skill{{PP: 1, Name: "test"}},
		},
	}
}

func (suite *monsterHTTPHandlerTestSuite) TestHandler_Fetch_ShouldSuccess() {
	svcMock := new(mocks.IPokewarService)
	svcMock.
		On("FetchMonsters", mock.Anything).
		Return(suite.monsters, nil).
		Once()
	svcMock.
		On("MonstersCount").
		Return(50).
		Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/ranks?limit=1", nil)
	ctx.Request = req
	handler := httpHandler.MonsterHTTPHandler{Svc: svcMock}
	handler.Fetch(ctx)
	var got utils.SuccessRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	assert.Equal(suite.T(), http.StatusOK, writer.Code)
	assert.Equal(suite.T(), http.StatusOK, got.Code)
	assert.Equal(suite.T(), http.StatusText(http.StatusOK), got.Status)
}

func (suite *monsterHTTPHandlerTestSuite) TestHandler_Fetch_ShouldError() {
	svcMock := new(mocks.IPokewarService)
	svcMock.
		On("FetchMonsters", mock.Anything).
		Return(nil, &utils.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "UNEXPECTED_ERROR",
		}).Once()
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	req, _ := http.NewRequest("GET", "/api/v1/monsters?limit=1", nil)
	ctx.Request = req
	handler := httpHandler.MonsterHTTPHandler{Svc: svcMock}
	handler.Fetch(ctx)
	var got utils.ErrorRespond
	_ = json.Unmarshal(writer.Body.Bytes(), &got)
	assert.Equal(suite.T(), http.StatusInternalServerError, writer.Code)
	assert.Equal(suite.T(), http.StatusInternalServerError, got.Code)
	assert.Equal(suite.T(), http.StatusText(http.StatusInternalServerError), got.Status)
	assert.Equal(suite.T(), "UNEXPECTED_ERROR", got.Data)
}

func TestMonsterHTTPHandler(t *testing.T) {
	suite.Run(t, new(monsterHTTPHandlerTestSuite))
}
