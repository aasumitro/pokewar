package ws

import (
	"encoding/json"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/battleroyale"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"time"
)

type MatchWSHandler struct {
	Svc         domain.IPokewarService
	Monsters    map[string][]*domain.Monster
	GamePlayers map[string][]*battleroyale.Player
	BattleData  map[string]*domain.Battle
	clients     map[string]*websocket.Conn
	mu          sync.Mutex
}

var isLastBattleSaved map[string]bool

var wsUpgraded = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (handler *MatchWSHandler) battleHistory(msgType int, clientId string) {
	data, errorData := handler.Svc.FetchBattles("LIMIT 5")
	if errorData != nil {
		message, _ := json.Marshal(map[string]any{
			"status":  "error",
			"message": errorData.Message,
		})
		handler.sendMessageToClient(msgType, clientId, message)
	}

	var histories []string
	for _, h := range data {
		var name string
		for _, p := range h.Players {
			if p.Rank == 1 {
				name = p.Name
				break
			}
		}

		history := fmt.Sprintf("BATTLE #%d WINNER [%s]", h.ID, strings.ToUpper(name))
		histories = append(histories, history)
	}

	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "battle_histories",
		"data":      histories,
	})
	handler.sendMessageToClient(msgType, clientId, message)
}

func (handler *MatchWSHandler) prepareBattle(msgType int, clientId string) {
	monsterData, errorData := handler.Svc.PrepareMonstersForBattle()
	if errorData != nil {
		message, _ := json.Marshal(map[string]any{
			"status":  "error",
			"message": errorData.Message,
		})
		handler.sendMessageToClient(msgType, clientId, message)
	}

	handler.Monsters[clientId] = monsterData
	handler.transformMonsterAsPlayer(clientId)

	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "monsters",
		"data":      monsterData,
	})
	handler.sendMessageToClient(msgType, clientId, message)
}

func (handler *MatchWSHandler) transformMonsterAsPlayer(clientId string) {
	var players []*battleroyale.Player
	for _, monster := range handler.Monsters[clientId] {
		var hp int
		for _, stat := range monster.Stats {
			if stat.Name == "hp" {
				hp = stat.BaseStat
			}
		}

		players = append(players, &battleroyale.Player{
			ID:     monster.ID,
			Name:   monster.Name,
			Health: hp,
			Score:  0,
			Rank:   0,
			Skills: []*battleroyale.Skill{
				{Power: monster.Skills[0].PP, Name: monster.Skills[0].Name},
				{Power: monster.Skills[1].PP, Name: monster.Skills[1].Name},
				{Power: monster.Skills[2].PP, Name: monster.Skills[2].Name},
				{Power: monster.Skills[3].PP, Name: monster.Skills[3].Name},
			},
		})
	}
	handler.GamePlayers[clientId] = players
}

func (handler *MatchWSHandler) startBattle(msgType int, clientId string) {
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "players",
		"data":      handler.GamePlayers[clientId],
	})
	handler.sendMessageToClient(msgType, clientId, message)

	result := make(chan *battleroyale.Game)
	log := make(chan string, 100)
	eliminated := make(chan string, 5)

	game := battleroyale.NewGame(handler.GamePlayers[clientId])
	go game.Start(result, log, eliminated)
	gameResult := handler.transformBattleResult(<-result)

	go func() {
		for logMessage := range log {
			message, _ = json.Marshal(map[string]any{
				"status":    "success",
				"data_type": "battle_logs",
				"data":      logMessage,
			})
			handler.sendMessageToClient(msgType, clientId, message)
		}
	}()

	go func() {
		for eliminatedPlayer := range eliminated {
			message, _ = json.Marshal(map[string]any{
				"status":    "success",
				"data_type": "eliminated_player",
				"data":      strings.ToUpper(eliminatedPlayer),
			})
			handler.sendMessageToClient(msgType, clientId, message)
		}
	}()

	handler.BattleData[clientId] = &gameResult
	message, _ = json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "battle_result",
		"data":      gameResult,
	})
	handler.sendMessageToClient(msgType, clientId, message)

	close(result)
	close(log)
	close(eliminated)
	game.Reset()
}

func (handler *MatchWSHandler) transformBattleResult(game *battleroyale.Game) domain.Battle {
	var logs []domain.Log
	for _, log := range game.Logs {
		logs = append(logs, domain.Log{
			Description: log.Description,
		})
	}

	var players []domain.Player
	for _, player := range game.Players {
		players = append(players, domain.Player{
			MonsterID: player.ID,
			Name:      player.Name,
			EliminatedAt: func() int64 {
				if player.EliminatedAt != nil {
					return player.EliminatedAt.UnixMicro()
				}
				return 0
			}(),
			Rank:  player.Rank,
			Point: player.Score,
		})
	}

	return domain.Battle{
		StartedAt: (*game).StartAt.UnixMicro(),
		EndedAt:   (*game).EndAt.UnixMicro(),
		Players:   players,
		Logs:      logs,
	}
}

func (handler *MatchWSHandler) annulledPlayer(msgType int, clientId string, data any) {
	var annulledPlayer *domain.Player
	for i := range handler.BattleData[clientId].Players {
		if handler.BattleData[clientId].Players[i].MonsterID == int(data.(float64)) {
			annulledPlayer = &handler.BattleData[clientId].Players[i]
			break
		}
	}

	for i := range handler.BattleData[clientId].Players {
		if handler.BattleData[clientId].Players[i].Rank > annulledPlayer.Rank {
			handler.BattleData[clientId].Players[i].Rank--
			handler.BattleData[clientId].Players[i].Point++
		}
	}

	annulledPlayer.AnnulledAt = time.Now().UnixMicro()
	annulledPlayer.Rank = 0
	annulledPlayer.Point = 0

	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "eliminated_result",
		"data":      handler.BattleData[clientId],
	})
	handler.sendMessageToClient(msgType, clientId, message)

	isLastBattleSaved[clientId] = false
	time.Sleep(10 * time.Second)
	if !isLastBattleSaved[clientId] && handler.BattleData[clientId] != nil {
		handler.save(clientId)
	}
}

func (handler *MatchWSHandler) save(clientId string) {
	if isLastBattleSaved[clientId] && handler.BattleData[clientId] == nil {
		return
	}

	err := handler.Svc.AddBattle(handler.BattleData[clientId])
	if err != nil {
		fmt.Println(err.Message)
	}

	isLastBattleSaved[clientId] = true
	handler.BattleData[clientId] = nil
	handler.Monsters[clientId] = nil
	handler.GamePlayers[clientId] = nil
}

func (handler *MatchWSHandler) sendMessageToClient(msgType int, clientId string, message []byte) {
	if conn, ok := handler.clients[clientId]; ok {
		handler.mu.Lock()
		_ = conn.WriteMessage(msgType, message)
		handler.mu.Unlock()
	}
}

func (handler *MatchWSHandler) Run(ctx *gin.Context) {
	idParams := ctx.Param("id")
	isLastBattleSaved[idParams] = true

	ws, _ := wsUpgraded.Upgrade(ctx.Writer, ctx.Request, nil)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(ws)
	handler.clients[idParams] = ws

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		var msg map[string]any
		_ = json.Unmarshal(message, &msg)
		switch msg["action"] {
		case "histories":
			handler.battleHistory(mt, msg["id"].(string))
		case "prepare":
			handler.prepareBattle(mt, msg["id"].(string))
		case "start":
			handler.startBattle(mt, msg["id"].(string))
		case "annulled":
			handler.annulledPlayer(mt, msg["id"].(string), msg["data"])
		case "save":
			handler.save(msg["id"].(string))
		}
	}
}

func NewMatchWSHandler(svc domain.IPokewarService, router *gin.RouterGroup) {
	isLastBattleSaved = make(map[string]bool)
	handler := MatchWSHandler{
		Svc:         svc,
		Monsters:    make(map[string][]*domain.Monster),
		GamePlayers: make(map[string][]*battleroyale.Player),
		BattleData:  make(map[string]*domain.Battle),
		clients:     make(map[string]*websocket.Conn),
	}
	router.GET("/ws/:id", handler.Run)
}
