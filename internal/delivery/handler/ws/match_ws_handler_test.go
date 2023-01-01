package ws_test

import (
	"encoding/json"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/internal/delivery/handler/ws"
	"github.com/aasumitro/pokewar/mocks"
	"github.com/aasumitro/pokewar/pkg/appconfig"
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
}

func (suite *matchWSHandlerTestSuite) SetupSuite() {
	viper.SetConfigFile("../../../../.example.env")

	appconfig.LoadEnv()

	appconfig.Instance.TotalMonsterSync = 10

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
			OriginID: 1,
			Name:     "test",
			BaseExp:  1,
			Height:   1,
			Weight:   1,
			Avatar:   "test.png",
			Types:    []string{"test", "test"},
			Stats:    []domain.Stat{{BaseStat: 1, Name: "hp"}},
			Skills:   []*domain.Skill{{PP: 1, Name: "test"}, {PP: 1, Name: "test"}, {PP: 1, Name: "test"}, {PP: 1, Name: "test"}},
		},
		{
			ID:       2,
			OriginID: 2,
			Name:     "test2",
			BaseExp:  1,
			Height:   1,
			Weight:   1,
			Avatar:   "test2.png",
			Types:    []string{"test2", "test2"},
			Stats:    []domain.Stat{{BaseStat: 1, Name: "hp"}},
			Skills:   []*domain.Skill{{PP: 1, Name: "test"}, {PP: 1, Name: "test"}, {PP: 1, Name: "test"}, {PP: 1, Name: "test"}},
		},
		{
			ID:       3,
			OriginID: 23,
			Name:     "test3",
			BaseExp:  1,
			Height:   1,
			Weight:   1,
			Avatar:   "test3.png",
			Types:    []string{"test3", "test3"},
			Stats:    []domain.Stat{{BaseStat: 1, Name: "hp"}},
			Skills:   []*domain.Skill{{PP: 1, Name: "test"}, {PP: 1, Name: "test"}, {PP: 1, Name: "test"}, {PP: 1, Name: "test"}},
		},
	}
}

func (suite *matchWSHandlerTestSuite) TestHandler_ActionHistory() {
	svc := new(mocks.IPokewarService)
	router := gin.New()
	ws.NewMatchWSHandler(svc, router.Group(""))
	server := httptest.NewServer(router)
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.Nil(suite.T(), err)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)

	ttERROR := []bool{true, false}
	for _, t := range ttERROR {
		if t {
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

		if !t {
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
	}
}

func (suite *matchWSHandlerTestSuite) TestHandler_ActionPrepare() {
	svc := new(mocks.IPokewarService)
	router := gin.New()
	ws.NewMatchWSHandler(svc, router.Group(""))
	server := httptest.NewServer(router)
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.Nil(suite.T(), err)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)

	ttERROR := []bool{true, false}
	for _, t := range ttERROR {
		if t {
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

		if !t {
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
	}
}

func (suite *matchWSHandlerTestSuite) TestHandler_ActionStart() {
	svc := new(mocks.IPokewarService)
	router := gin.New()
	ws.NewMatchWSHandler(svc, router.Group(""))
	server := httptest.NewServer(router)
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/1"
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	require.Nil(suite.T(), err)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(wsConn)
	ttCASE := []string{"NIL_GAME_PLAYER", "NOT_NIL_GAME_PLAYER"}

	for _, t := range ttCASE {
		if t == "NIL_GAME_PLAYER" {
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
			require.Equal(suite.T(), msg["status"], "error")
			require.Equal(suite.T(), msg["message"], "Please press random button again!")
		}

		if t == "NOT_NIL_GAME_PLAYER" {
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
			if err := wsConn.WriteMessage(
				websocket.TextMessage,
				[]byte(`{"action": "start", "id": "1"}`),
			); err != nil {
				suite.T().Fatalf("failed to write message: %v", err)
			}
			for {
				_, message2, err := wsConn.ReadMessage()
				require.Nil(suite.T(), err)
				var msg2 map[string]interface{}
				err = json.Unmarshal(message2, &msg2)
				require.Nil(suite.T(), err)
				require.Equal(suite.T(), msg2["status"], "success")
				if !utils.InArray[string](msg2["data_type"].(string), []string{"battle_logs", "eliminated_player", "battle_result"}) {
					suite.T().Errorf("not expected data type %s", msg2["data_type"])
				}
				if msg2["data_type"] == "battle_result" {
					break
				}
			}
		}
	}
}

func (suite *matchWSHandlerTestSuite) TestHandler_Annulled_Save() {
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
	if err := wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"action": "start", "id": "1"}`),
	); err != nil {
		suite.T().Fatalf("failed to write message: %v", err)
	}
	for {
		_, message2, err := wsConn.ReadMessage()
		require.Nil(suite.T(), err)
		var msg2 map[string]interface{}
		err = json.Unmarshal(message2, &msg2)
		require.Nil(suite.T(), err)
		require.Equal(suite.T(), msg2["status"], "success")
		if !utils.InArray[string](msg2["data_type"].(string), []string{"battle_logs", "eliminated_player", "battle_result"}) {
			suite.T().Errorf("not expected data type %s", msg2["data_type"])
		}
		if msg2["data_type"] == "battle_result" {
			break
		}
	}

	if err := wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"action": "annulled", "id": "1", "data": 1}`),
	); err != nil {
		suite.T().Fatalf("failed to write message: %v", err)
	}
	for i := 1; i <= 2; i++ {
		_, message3, err := wsConn.ReadMessage()
		require.Nil(suite.T(), err)
		var msg3 map[string]interface{}
		err = json.Unmarshal(message3, &msg3)
		require.Nil(suite.T(), err)
		require.Equal(suite.T(), msg3["status"], "success")
	}

	svc.On("AddBattle", mock.Anything).
		Return(nil).Once()
	if err := wsConn.WriteMessage(
		websocket.TextMessage,
		[]byte(`{"action": "save", "id": "1"}`),
	); err != nil {
		suite.T().Fatalf("failed to write message: %v", err)
	}
}

func TestMatchWSHandler(t *testing.T) {
	suite.Run(t, new(matchWSHandlerTestSuite))
}
