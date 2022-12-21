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
	playground  string
	Svc         domain.IPokewarService
	Monsters    []*domain.Monster
	GamePlayers []*battleroyale.Player
	BattleData  *domain.Battle
	wsConn      *websocket.Conn
	mu          sync.Mutex
}

var wsUpgraded = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (handler *MatchWSHandler) battleHistory(msgType int) {
	data, errorData := handler.Svc.FetchBattles("LIMIT 5")
	if errorData != nil {
		handler.mu.Lock()
		message, _ := json.Marshal(map[string]any{
			"status":  "error",
			"message": errorData.Message,
		})
		_ = handler.wsConn.WriteMessage(msgType, message)
		handler.mu.Unlock()
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

	handler.mu.Lock()
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "battle_histories",
		"data":      histories,
	})
	_ = handler.wsConn.WriteMessage(msgType, message)
	handler.mu.Unlock()
}

func (handler *MatchWSHandler) prepareBattle(msgType int) {
	monsterData, errorData := handler.Svc.PrepareMonstersForBattle()
	if errorData != nil {
		handler.mu.Lock()
		message, _ := json.Marshal(map[string]any{
			"status":  "error",
			"message": errorData.Message,
		})
		_ = handler.wsConn.WriteMessage(msgType, message)
		handler.mu.Unlock()
	}

	handler.Monsters = monsterData
	handler.transformMonsterAsPlayer()

	handler.mu.Lock()
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "monsters",
		"data":      monsterData,
	})
	_ = handler.wsConn.WriteMessage(msgType, message)
	handler.mu.Unlock()
}

func (handler *MatchWSHandler) transformMonsterAsPlayer() {
	var players []*battleroyale.Player
	for _, monster := range handler.Monsters {
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
	handler.GamePlayers = players
}

func (handler *MatchWSHandler) startBattle(msgType int) {
	handler.mu.Lock()
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "players",
		"data":      handler.GamePlayers,
	})
	_ = handler.wsConn.WriteMessage(msgType, message)
	handler.mu.Unlock()

	result := make(chan *battleroyale.Game)
	log := make(chan string, 100)
	eliminated := make(chan string, 5)

	game := battleroyale.NewGame(handler.GamePlayers)
	go game.Start(result, log, eliminated)
	gameResult := handler.transformBattleResult(<-result)

	go func() {
		for logMessage := range log {
			handler.mu.Lock()
			message, _ = json.Marshal(map[string]any{
				"status":    "success",
				"data_type": "battle_logs",
				"data":      logMessage,
			})
			_ = handler.wsConn.WriteMessage(msgType, message)
			handler.mu.Unlock()
		}
	}()

	go func() {
		for eliminatedPlayer := range eliminated {
			handler.mu.Lock()
			message, _ = json.Marshal(map[string]any{
				"status":    "success",
				"data_type": "eliminated_player",
				"data":      strings.ToUpper(eliminatedPlayer),
			})
			_ = handler.wsConn.WriteMessage(msgType, message)
			handler.mu.Unlock()
		}
	}()

	handler.BattleData = &gameResult
	handler.mu.Lock()
	message, _ = json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "battle_result",
		"data":      gameResult,
	})
	_ = handler.wsConn.WriteMessage(msgType, message)
	handler.mu.Unlock()

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

func (handler *MatchWSHandler) annulledPlayer(msgType int, data any) {
	var annulledPlayer *domain.Player
	for i := range handler.BattleData.Players {
		if handler.BattleData.Players[i].MonsterID == int(data.(float64)) {
			annulledPlayer = &handler.BattleData.Players[i]
			break
		}
	}

	for i := range handler.BattleData.Players {
		if handler.BattleData.Players[i].Rank > annulledPlayer.Rank {
			handler.BattleData.Players[i].Rank--
			handler.BattleData.Players[i].Point++
		}
	}

	annulledPlayer.AnnulledAt = time.Now().UnixMicro()
	annulledPlayer.Rank = 0
	annulledPlayer.Point = 0

	handler.mu.Lock()
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "eliminated_result",
		"data":      handler.BattleData,
	})
	_ = handler.wsConn.WriteMessage(msgType, message)
	handler.mu.Unlock()
}

func (handler *MatchWSHandler) save(msgType int) {
	err := handler.Svc.AddBattle(handler.BattleData)
	if err != nil {
		fmt.Println(err.Message)
	}
}

func (handler *MatchWSHandler) Run(ctx *gin.Context) {
	idParams := ctx.Param("id")
	handler.playground = idParams
	ws, _ := wsUpgraded.Upgrade(ctx.Writer, ctx.Request, nil)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(ws)
	handler.wsConn = ws

	for {
		mt, message, _ := ws.ReadMessage()
		//if err != nil {
		//handler.mu.Lock()
		//message = []byte("something went wrong")
		//_ = ws.WriteMessage(mt, message)
		//handler.mu.Unlock()
		//}

		var msg map[string]any
		_ = json.Unmarshal(message, &msg)
		switch msg["action"] {
		case "histories":
			handler.battleHistory(mt)
		case "prepare":
			handler.prepareBattle(mt)
		case "start":
			handler.startBattle(mt)
		case "annulled":
			handler.annulledPlayer(mt, msg["data"])
		case "save":
			handler.save(mt)
		}
	}
}

func NewMatchWSHandler(svc domain.IPokewarService, router *gin.RouterGroup) {
	handler := MatchWSHandler{Svc: svc}
	router.GET("/ws/:id", handler.Run)
}
