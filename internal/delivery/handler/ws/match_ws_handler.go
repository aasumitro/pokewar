package ws

import (
	"encoding/json"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MatchWSHandler struct {
	Svc      domain.IPokewarService
	Monsters []*domain.Monster
	wsConn   *websocket.Conn
}

var wsUpgraded = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (handler *MatchWSHandler) prepareBattle(msgType int) {
	monsterData, errorData := handler.Svc.PrepareMonstersForBattle()
	if errorData != nil {
		message, _ := json.Marshal(map[string]any{
			"status":  "error",
			"message": errorData.Message,
		})
		_ = handler.wsConn.WriteMessage(msgType, message)
	}

	handler.Monsters = monsterData

	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "monsters",
		"data":      monsterData,
	})
	_ = handler.wsConn.WriteMessage(msgType, message)
}

func (handler *MatchWSHandler) startBattle(msgType int) {
	fmt.Println(handler.Monsters)
	err := handler.wsConn.WriteMessage(msgType, []byte("start"))
	if err != nil {
		fmt.Println(err)
	}
}

func (handler *MatchWSHandler) annulledPlayer(playerId int) {
	fmt.Println(playerId)
}

func (handler *MatchWSHandler) Run(ctx *gin.Context) {
	//idParams := ctx.Param("id")
	ws, err := wsUpgraded.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) { _ = ws.Close() }(ws)
	handler.wsConn = ws

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			message = []byte("something went wrong")
			_ = ws.WriteMessage(mt, message)
		}

		var msg map[string]any
		_ = json.Unmarshal(message, &msg)
		switch msg["action"] {
		case "prepare":
			handler.prepareBattle(mt)
			break
		case "start":
			handler.startBattle(mt)
			break
		case "annulled":
			handler.annulledPlayer(msg["data"].(int))
			break
		}
	}
}

func NewMatchWSHandler(svc domain.IPokewarService, router *gin.RouterGroup) {
	handler := MatchWSHandler{Svc: svc}
	router.GET("/ws/:id", handler.Run)
}
