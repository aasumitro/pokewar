package ws_test

import (
	"encoding/json"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/internal/delivery/handler/ws"
	"github.com/aasumitro/pokewar/mocks"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/aasumitro/pokewar/pkg/battleroyale"
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type matchWSHandlerTestSuite struct {
	suite.Suite
	battles  []*domain.Battle
	monsters []*domain.Monster
	players  []*battleroyale.Player
}

func (suite *matchWSHandlerTestSuite) SetupSuite() {
	viper.SetConfigFile("../../../../.example.env")

	appconfigs.LoadEnv()

	appconfigs.Instance.TotalMonsterSync = 10

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
			Stats:    []domain.Stat{{BaseStat: 1, Name: "hp"}},
			Skills:   []*domain.Skill{{PP: 1, Name: "test"}, {PP: 1, Name: "test"}, {PP: 1, Name: "test"}, {PP: 1, Name: "test"}},
		},
	}
	suite.players = []*battleroyale.Player{
		{ID: 1, Name: "Player 1", Health: 100, Skills: []*battleroyale.Skill{
			{Name: "Kick", Power: 20},
			{Name: "Punch", Power: 10},
		}},
		{ID: 2, Name: "Player 2", Health: 100, Skills: []*battleroyale.Skill{
			{Name: "Kick", Power: 20},
			{Name: "Punch", Power: 10},
		}},
		{ID: 3, Name: "Player 3", Health: 100, Skills: []*battleroyale.Skill{
			{Name: "Kick", Power: 20},
			{Name: "Punch", Power: 10},
		}},
	}
}

func (suite *matchWSHandlerTestSuite) TestHandler_ActionHistory_ShouldSuccess() {
	svc := new(mocks.IPokewarService)
	router := gin.New()
	ws.NewMatchWSHandler(svc, router.Group(""))
	server := httptest.NewServer(router)
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.Nil(suite.T(), err)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)

	svc.On("FetchBattles", mock.Anything).
		Return(suite.battles, nil).Once()
	if err := wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"action": "histories", "id": "1"}`),
	); err != nil {
		suite.T().Fatalf("failed to write message: %v", err)
	}
	_, message, err := wsConn.ReadMessage()
	require.Nil(suite.T(), err)
	var msg map[string]interface{}
	err = json.Unmarshal(message, &msg)
	require.Nil(suite.T(), err)

	require.Equal(suite.T(), msg["status"], "success")
	require.Equal(suite.T(), msg["data_type"], "battle_histories")
}
func (suite *matchWSHandlerTestSuite) TestHandler_ActionHistory_ShouldErrorFromService() {
	svc := new(mocks.IPokewarService)
	router := gin.New()
	ws.NewMatchWSHandler(svc, router.Group(""))
	server := httptest.NewServer(router)
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.Nil(suite.T(), err)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)

	svc.On("FetchBattles", mock.Anything).
		Return(nil, &utils.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "UNEXPECTED_ERROR",
		}).Once()
	if err := wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"action": "histories", "id": "1"}`),
	); err != nil {
		suite.T().Fatalf("failed to write message: %v", err)
	}
	_, message, err := wsConn.ReadMessage()
	require.Nil(suite.T(), err)
	var msg map[string]interface{}
	err = json.Unmarshal(message, &msg)
	require.Nil(suite.T(), err)

	require.Equal(suite.T(), msg["status"], "error")
	require.Equal(suite.T(), msg["message"], "UNEXPECTED_ERROR")
}

func (suite *matchWSHandlerTestSuite) TestHandler_ActionPrepare_ShouldSuccess() {
	svc := new(mocks.IPokewarService)
	router := gin.New()
	ws.NewMatchWSHandler(svc, router.Group(""))
	server := httptest.NewServer(router)
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.Nil(suite.T(), err)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)

	svc.On("PrepareMonstersForBattle").
		Return(suite.monsters, nil).Once()
	if err := wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"action": "prepare", "id": "1"}`),
	); err != nil {
		suite.T().Fatalf("failed to write message: %v", err)
	}
	_, message, err := wsConn.ReadMessage()
	require.Nil(suite.T(), err)
	var msg map[string]interface{}
	err = json.Unmarshal(message, &msg)
	require.Nil(suite.T(), err)

	require.Equal(suite.T(), msg["status"], "success")
	require.Equal(suite.T(), msg["data_type"], "monsters")
}
func (suite *matchWSHandlerTestSuite) TestHandler_ActionPrepare_ShouldError() {
	svc := new(mocks.IPokewarService)
	router := gin.New()
	ws.NewMatchWSHandler(svc, router.Group(""))
	server := httptest.NewServer(router)
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.Nil(suite.T(), err)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)

	svc.On("PrepareMonstersForBattle").
		Return(nil, &utils.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "UNEXPECTED_ERROR",
		}).Once()
	if err := wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"action": "prepare", "id": "1"}`),
	); err != nil {
		suite.T().Fatalf("failed to write message: %v", err)
	}
	_, message, err := wsConn.ReadMessage()
	require.Nil(suite.T(), err)
	var msg map[string]interface{}
	err = json.Unmarshal(message, &msg)
	require.Nil(suite.T(), err)

	require.Equal(suite.T(), msg["status"], "error")
	require.Equal(suite.T(), msg["message"], "UNEXPECTED_ERROR")
}

func (suite *matchWSHandlerTestSuite) TestHandler_ActionStart_ShouldSuccess() {
	suite.T().Skip()
	svc := new(mocks.IPokewarService)
	router := gin.New()
	ws.NewMatchWSHandler(svc, router.Group(""))
	server := httptest.NewServer(router)
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.Nil(suite.T(), err)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)

	if err := wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"action": "start", "id": "1"}`),
	); err != nil {
		suite.T().Fatalf("failed to write message: %v", err)
	}
	_, message, err := wsConn.ReadMessage()
	require.Nil(suite.T(), err)
	var msg map[string]interface{}
	err = json.Unmarshal(message, &msg)
	require.Nil(suite.T(), err)

	require.Equal(suite.T(), msg["status"], "success")
	require.Equal(suite.T(), msg["data_type"], "players")
}

//func (suite *matchWSHandlerTestSuite) TestHandler_Action_() {
//svc := new(mocks.IPokewarService)
//router := gin.New()
//ws.NewMatchWSHandler(svc, router.Group(""))
//server := httptest.NewServer(router)
//defer server.Close()
//
//u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
//wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
//require.Nil(suite.T(), err)
//defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)
//
//svc.On("FetchBattles", mock.Anything).Return(suite.battles, nil).Once()
//if err := wsConn.WriteMessage(
//	websocket.TextMessage,
//	[]byte(`{"action": "histories", "id": "1"}`),
//); err != nil {
//	suite.T().Fatalf("failed to write message: %v", err)
//}
//_, message, err := wsConn.ReadMessage()
//require.Nil(suite.T(), err)
//var msg map[string]interface{}
//err = json.Unmarshal(message, &msg)
//require.Nil(suite.T(), err)
//require.Equal(suite.T(), msg["status"], "success")
//require.Equal(suite.T(), msg["data_type"], "battle_histories")
//}

// TODO UPDATE COVERAGE

func TestMatchWSHandler(t *testing.T) {
	suite.Run(t, new(matchWSHandlerTestSuite))
}
