package ws

import (
	"encoding/json"
	"fmt"
	"github.com/aasumitro/pokewar/constants"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/battleroyale"
	"github.com/aasumitro/pokewar/pkg/datatransform"
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
}

var (
	isLastBattleSaved map[string]bool
	clients           map[string]*websocket.Conn
	wsUpgraded        = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	mu sync.Mutex
)

// battleHistory handles the "histories" request message type from the client.
// Proceed an error message or the list of battles history that extracted to
// the relevant information and format, and sends it back to the client.
func (handler *MatchWSHandler) battleHistory(msgType int, clientID string) {
	// Fetch the last 5 battles from the database via service
	data, errorData := handler.Svc.FetchBattles("LIMIT 5")
	if errorData != nil {
		// If there was an error fetching the battles,
		// send an error message to the client
		handler.sendMessageToClient(
			msgType, clientID, "error",
			errorData.Message.(string), nil)
	}
	// Pre-allocate capacity for the histories slice using make
	histories := make([]string, len(data))
	// Iterate over the battles and create a history string for each battle
	for i, h := range data {
		// Find the name of the player who won the battle
		var name string
		for _, p := range h.Players {
			if p.Rank == 1 {
				name = strings.ToUpper(p.Name)
				break
			}
		}
		// Calculate the duration of the battle
		endedAt := time.Duration(h.EndedAt) * time.Microsecond
		startedAt := time.Duration(h.StartedAt) * time.Microsecond
		diff := (endedAt - startedAt).String()
		// Create the history string for the current battle
		template := "BATTLE#%d - %s <br>WINNER [%s]<br><hr class='py-2'>"
		history := fmt.Sprintf(template, h.ID, diff, name)
		histories[i] = history
	}
	// Send the histories to the client
	handler.sendMessageToClient(
		msgType, clientID, "success",
		"battle_histories", histories)
}

// prepareBattle handles the "prepare" request message type from the client.
// Proceed an error message or the list of monsters then transform the list of monsters
// into a list of players and sends the list of monsters back to the client.
func (handler *MatchWSHandler) prepareBattle(msgType int, clientID string) {
	// Fetch the list of monsters from the database
	monsterData, errorData := handler.Svc.PrepareMonstersForBattle()
	if errorData != nil {
		// If there was an error fetching the monsters,
		// send an error message to the client
		handler.sendMessageToClient(
			msgType, clientID, "error",
			errorData.Message.(string), nil)
	}
	// Store the list of monsters and transform
	// the data to battleroyale.players
	handler.Monsters[clientID] = monsterData
	players := datatransform.TransformMonstersAsGamePlayers(monsterData)
	handler.GamePlayers[clientID] = players
	// Send the list of monsters to the client
	handler.sendMessageToClient(
		msgType, clientID, "success",
		"monsters", monsterData)
}

// startBattle handles the "start" request message type from the client.
// It starts a new battle game for a specified client, sends data (logs, eliminated players, and result)
// back to the client after transforming it into the specified format, and resets the game.
func (handler *MatchWSHandler) startBattle(msgType int, clientID string) {
	if handler.GamePlayers[clientID] == nil {
		// if players not set don't play the game
		// and send notify to user load random monster
		// to play the game and start the match
		handler.sendMessageToClient(
			msgType, clientID, "error",
			"Please press random button again!", nil)
		return
	}
	// Create a buffered channel with a capacity of 10 to store log updates and eliminated players
	updateBuffer := make(chan map[string]any, constants.MaxWSUpdateBufferSize)
	// Start a goroutine to handle sending updates to the client
	go func() {
		for update := range updateBuffer {
			handler.sendMessageToClient(
				msgType, clientID, update["status"].(string),
				update["data_type"].(string), update["data"])
		}
	}()
	// Use a fixed-size buffer channel for the result, log, and eliminated channels
	result := make(chan *battleroyale.Game, 1)
	log := make(chan string, constants.MaxGameLogSize)
	eliminated := make(chan string, constants.MaxPlayerSize)
	// Start a new game and transform the result to domain.Battle
	game := battleroyale.NewGame(handler.GamePlayers[clientID])
	go game.Start(result, log, eliminated)
	// Start a goroutine to handle log updates and eliminated players
	go func() {
		for {
			select {
			case logMessage := <-log:
				updateBuffer <- map[string]any{
					"status":    "success",
					"data_type": "battle_logs",
					"data":      logMessage,
				}
			case eliminatedPlayer := <-eliminated:
				updateBuffer <- map[string]any{
					"status":    "success",
					"data_type": "eliminated_player",
					"data":      strings.ToUpper(eliminatedPlayer),
				}
			case <-result:
				// Close the updateBuffer channel when the result channel is closed
				close(updateBuffer)
				return
			}
		}
	}()
	// Send game result to client
	gameResult := datatransform.TransformGameResultToBattle(<-result)
	handler.BattleData[clientID] = gameResult
	handler.sendMessageToClient(
		msgType, clientID, "success",
		"battle_result", gameResult)
	// Reset the game
	game.Reset()
	// force safe data after 5 seconds
	isLastBattleSaved[clientID] = false
	time.AfterFunc(constants.SaveDuration, func() {
		handler.save(clientID)
	})
}

// annulledPlayer handles the "annulled" request message type from the client.
// It finds the index of the player to be annulled increment/decrement-ing the rank and points of other players
// after that make new log, sends the updated data to the client and schedules the save function
// to be called after 10 seconds using time.AfterFunc
func (handler *MatchWSHandler) annulledPlayer(msgType int, clientID string, data any) {
	// Find the index of the player to be annulled
	var playerIndex int
	for i, player := range handler.BattleData[clientID].Players {
		if player.MonsterID == int(data.(float64)) {
			playerIndex = i
			break
		}
	}
	// Get a reference to the player to be annulled
	annulledPlayer := &handler.BattleData[clientID].Players[playerIndex]
	// Increment the rank and points of players with a higher rank than the annulled player
	for i := range handler.BattleData[clientID].Players {
		if handler.BattleData[clientID].Players[i].Rank > annulledPlayer.Rank {
			handler.BattleData[clientID].Players[i].Rank--
			handler.BattleData[clientID].Players[i].Point++
		}
	}
	// Update the annulled player's data
	annulledPlayer.AnnulledAt = time.Now().UnixMicro()
	annulledPlayer.Rank = 0
	annulledPlayer.Point = 0
	tplMsg := "%d - player %s was annulled from the game and " +
		"their rank and points were reset to 0!\n"
	logMsg := fmt.Sprintf(tplMsg, annulledPlayer.AnnulledAt, annulledPlayer.Name)
	handler.BattleData[clientID].Logs = append(
		handler.BattleData[clientID].Logs,
		domain.Log{Description: logMsg},
	)
	// Send log message to the client
	handler.sendMessageToClient(
		msgType, clientID, "success",
		"battle_logs", logMsg)
	// Send updated data to the client
	handler.sendMessageToClient(
		msgType, clientID, "success",
		"eliminated_result", handler.BattleData[clientID])
	// Schedule the save function to be called after 10 seconds
	isLastBattleSaved[clientID] = false
	time.AfterFunc(1*time.Second, func() {
		handler.save(clientID)
	})
}

// save handles the "save" request message type from the client, or Schedule event from annulledPlayer
// It attempts to store the battle data. If the database connection is unavailable,
// the function will retry for a maximum of 5 times with a 500ms delay between each retry.
// If the data stored or the maximum number of retries is reached,
// the battle data and other related data for the client will be reset.
func (handler *MatchWSHandler) save(clientID string) {
	if isLastBattleSaved[clientID] && handler.BattleData[clientID] == nil {
		return
	}
	// Set the maximum number of retries and
	// the initial retry counter
	maxRetries, retryCount := 3, 0
	for {
		err := handler.Svc.AddBattle(handler.BattleData[clientID])
		// If the function succeeds,
		// break out of the loop
		if err == nil {
			break
		}
		retryCount++
		// If the maximum number of retries is reached,
		// break out of the loop
		if retryCount >= maxRetries {
			break
		}
		// Sleep for a short period before retrying
		time.Sleep(constants.SleepDuration)
	}

	// reset the data
	isLastBattleSaved[clientID] = true
	handler.BattleData[clientID] = nil
	handler.Monsters[clientID] = nil
	handler.GamePlayers[clientID] = nil
}

// sendMessageToClient helper function to send message to specified client by given id
func (handler *MatchWSHandler) sendMessageToClient(
	msgType int,
	clientID, status, dt string,
	data any,
) {
	if conn, ok := clients[clientID]; ok {
		mu.Lock()
		var message []byte
		if status == "error" {
			message, _ = json.Marshal(map[string]any{
				"status":  status,
				"message": dt,
			})
		}
		if status == "success" {
			message, _ = json.Marshal(map[string]any{
				"status":    status,
				"data_type": dt,
				"data":      data,
			})
		}
		_ = conn.WriteMessage(msgType, message)
		mu.Unlock()
	}
}

// Run function is a WebSocket handler for a Battleroyale game.
// It handles incoming messages from clients and performs
// the appropriate action based on the value of the "action" field in the message.
func (handler *MatchWSHandler) Run(ctx *gin.Context) {
	// get client/playground id
	idParams := ctx.Param("id")
	isLastBattleSaved[idParams] = true
	// make a new connection and store the specified client connection
	ws, _ := wsUpgraded.Upgrade(ctx.Writer, ctx.Request, nil)
	defer func(ws *websocket.Conn) { _ = ws.Close() }(ws)
	clients[idParams] = ws
	// starts a loop that listens for incoming messages from the client and processes them.
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		// When the client sends a message with the registered case
		// call the specified function and do the action
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
	// initialize struct and variable value
	isLastBattleSaved = make(map[string]bool)
	clients = make(map[string]*websocket.Conn)
	handler := MatchWSHandler{
		Svc:         svc,
		Monsters:    make(map[string][]*domain.Monster),
		GamePlayers: make(map[string][]*battleroyale.Player),
		BattleData:  make(map[string]*domain.Battle),
	}
	// register the user
	router.GET("/ws/:id", handler.Run)
}
