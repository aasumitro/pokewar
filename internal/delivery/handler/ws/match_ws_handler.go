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
func (handler *MatchWSHandler) battleHistory(msgType int, clientId string) {
	// Fetch the last 5 battles from the database via service
	data, errorData := handler.Svc.FetchBattles("LIMIT 5")
	if errorData != nil {
		// If there was an error fetching the battles,
		// send an error message to the client
		message, _ := json.Marshal(map[string]any{
			"status":  "error",
			"message": errorData.Message,
		})
		handler.sendMessageToClient(msgType, clientId, message)
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
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "battle_histories",
		"data":      histories,
	})
	handler.sendMessageToClient(msgType, clientId, message)
}

// prepareBattle handles the "prepare" request message type from the client.
// Proceed an error message or the list of monsters then transform the list of monsters
// into a list of players and sends the list of monsters back to the client.
func (handler *MatchWSHandler) prepareBattle(msgType int, clientId string) {
	// Fetch the list of monsters from the database
	monsterData, errorData := handler.Svc.PrepareMonstersForBattle()
	if errorData != nil {
		// If there was an error fetching the monsters,
		// send an error message to the client
		message, _ := json.Marshal(map[string]any{
			"status":  "error",
			"message": errorData.Message,
		})
		handler.sendMessageToClient(msgType, clientId, message)
	}
	// Store the list of monsters and transform
	// the data to battleroyale.players
	handler.Monsters[clientId] = monsterData
	handler.transformMonsterAsPlayer(clientId)
	// Send the list of monsters to the client
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "monsters",
		"data":      monsterData,
	})
	handler.sendMessageToClient(msgType, clientId, message)
}

// transformMonsterAsPlayer transforms a list of monsters into a list of players
// by extracting relevant information from the monsters.
func (handler *MatchWSHandler) transformMonsterAsPlayer(clientId string) {
	// Pre-allocate capacity for the players slice using make
	players := make([]*battleroyale.Player, len(handler.Monsters[clientId]))
	// Iterate through each monster and extract the necessary information to create a player
	for i, monster := range handler.Monsters[clientId] {
		// Find the monster's HP stat
		var hp int
		for _, stat := range monster.Stats {
			if stat.Name == "hp" {
				hp = stat.BaseStat
				break
			}
		}
		// Create a player using the extracted information
		players[i] = &battleroyale.Player{
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
		}
	}
	// Save the list of players to the GamePlayers map
	handler.GamePlayers[clientId] = players
}

// startBattle handles the "start" request message type from the client.
// It starts a new battle game for a specified client, sends data (logs, eliminated players, and result)
// back to the client after transforming it into the specified format, and resets the game.
func (handler *MatchWSHandler) startBattle(msgType int, clientId string) {
	// Create a buffered channel with a capacity of 10 to store log updates and eliminated players
	updateBuffer := make(chan map[string]any, 10)
	// Start a goroutine to handle sending updates to the client
	go func() {
		for update := range updateBuffer {
			message, _ := json.Marshal(update)
			handler.sendMessageToClient(msgType, clientId, message)
		}
	}()
	// Use a fixed-size buffer channel for the result, log, and eliminated channels
	result := make(chan *battleroyale.Game, 1)
	log := make(chan string, 100)
	eliminated := make(chan string, 5)
	// Start a new game and transform the result
	game := battleroyale.NewGame(handler.GamePlayers[clientId])
	go game.Start(result, log, eliminated)
	gameResult := handler.transformBattleResult(<-result)
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
	handler.BattleData[clientId] = &gameResult
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "battle_result",
		"data":      gameResult,
	})
	handler.sendMessageToClient(msgType, clientId, message)
	// Reset the game
	game.Reset()
}

// transformBattleResult converts the data from the battleroyale.Game struct to the domain.Battle struct format.
func (handler *MatchWSHandler) transformBattleResult(game *battleroyale.Game) domain.Battle {
	// pre-allocates the slices for logs sing the length of game.Logs
	logs := make([]domain.Log, len(game.Logs))
	for i, log := range game.Logs {
		// populates and transform data.
		logs[i] = domain.Log{Description: log.Description}
	}
	// pre-allocates the slices for logs sing the length of game.Players
	players := make([]domain.Player, len(game.Players))
	for i, player := range game.Players {
		// populates and transform data.
		players[i] = domain.Player{
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
		}
	}
	// returns the transformed data.
	return domain.Battle{
		StartedAt: (*game).StartAt.UnixMicro(),
		EndedAt:   (*game).EndAt.UnixMicro(),
		Players:   players,
		Logs:      logs,
	}
}

// annulledPlayer handles the "annulled" request message type from the client.
// It finds the index of the player to be annulled increment/decrement-ing the rank and points of other players
// after that make new log, sends the updated data to the client and schedules the save function
// to be called after 10 seconds using time.AfterFunc
func (handler *MatchWSHandler) annulledPlayer(msgType int, clientId string, data any) {
	// Find the index of the player to be annulled
	var playerIndex int
	for i, player := range handler.BattleData[clientId].Players {
		if player.MonsterID == int(data.(float64)) {
			playerIndex = i
			break
		}
	}
	// Get a reference to the player to be annulled
	annulledPlayer := &handler.BattleData[clientId].Players[playerIndex]
	// Increment the rank and points of players with a higher rank than the annulled player
	for i := range handler.BattleData[clientId].Players {
		if handler.BattleData[clientId].Players[i].Rank > annulledPlayer.Rank {
			handler.BattleData[clientId].Players[i].Rank--
			handler.BattleData[clientId].Players[i].Point++
		}
	}
	// Update the annulled player's data
	annulledPlayer.AnnulledAt = time.Now().UnixMicro()
	annulledPlayer.Rank = 0
	annulledPlayer.Point = 0
	tplMsg := "%d - player %s was annulled from the game and " +
		"their rank and points were reset to 0!\n"
	logMsg := fmt.Sprintf(tplMsg, annulledPlayer.AnnulledAt, annulledPlayer.Name)
	handler.BattleData[clientId].Logs = append(
		handler.BattleData[clientId].Logs,
		domain.Log{Description: logMsg},
	)
	// Send log message to the client
	message, _ := json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "battle_logs",
		"data":      logMsg,
	})
	handler.sendMessageToClient(msgType, clientId, message)
	// Send updated data to the client
	message, _ = json.Marshal(map[string]any{
		"status":    "success",
		"data_type": "eliminated_result",
		"data":      handler.BattleData[clientId],
	})
	handler.sendMessageToClient(msgType, clientId, message)
	// Schedule the save function to be called after 10 seconds
	isLastBattleSaved[clientId] = false
	time.AfterFunc(10*time.Second, func() {
		handler.save(clientId)
	})
}

// save handles the "save" request message type from the client, or Schedule event from annulledPlayer
// It attempts to store the battle data. If the database connection is unavailable,
// the function will retry for a maximum of 5 times with a 500ms delay between each retry.
// If the data stored or the maximum number of retries is reached,
// the battle data and other related data for the client will be reset.
func (handler *MatchWSHandler) save(clientId string) {
	if isLastBattleSaved[clientId] && handler.BattleData[clientId] == nil {
		return
	}
	// Set the maximum number of retries and
	// the initial retry counter
	maxRetries, retryCount := 5, 0
	for {
		err := handler.Svc.AddBattle(handler.BattleData[clientId])
		// If the function succeeds,
		// break out of the loop
		if err == nil {
			break
		}
		fmt.Println("gagal", err.Message)
		retryCount++
		// If the maximum number of retries is reached,
		// break out of the loop
		if retryCount >= maxRetries {
			break
		}
		// Sleep for a short period before retrying
		time.Sleep(500 * time.Millisecond)
	}
	// reset the data
	isLastBattleSaved[clientId] = true
	handler.BattleData[clientId] = nil
	handler.Monsters[clientId] = nil
	handler.GamePlayers[clientId] = nil
}

// sendMessageToClient helper function to send message to specified client by given id
func (handler *MatchWSHandler) sendMessageToClient(msgType int, clientId string, message []byte) {
	if conn, ok := clients[clientId]; ok {
		mu.Lock()
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
	//starts a loop that listens for incoming messages from the client and processes them.
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
