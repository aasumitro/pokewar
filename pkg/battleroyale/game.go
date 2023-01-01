package battleroyale

import (
	"fmt"
	"github.com/aasumitro/pokewar/pkg/constant"
	"math/rand"
	"sort"
	"time"
)

type (
	// Game represents the state of a battle royal game.
	Game struct {
		Players                 []*Player
		Winner                  *Player
		StartAt                 time.Time
		EndAt                   time.Time
		Logs                    []Log
		BattleLogsChannel       chan string
		EliminatedPlayerChannel chan string
	}

	// Log represent of player action to eliminate other player
	Log struct {
		Description string
	}

	// IGameAction contract that can be accessed from outer universe #lol
	IGameAction interface {
		// Start begins the game loop.
		Start(game chan *Game, log chan string, eliminated chan string)
		// Reset the game
		Reset()
	}
)

// Start begins the game loop. It takes three channels as arguments:
// game: A channel used to send the final game data when the game ends.
// log: A channel used to send battle logs.
// eliminated: A channel used to send the names of eliminated players.
func (g *Game) Start(game chan *Game, log chan string, eliminated chan string) {
	// Assign the log and eliminated channels to the game struct.
	g.BattleLogsChannel = log
	g.EliminatedPlayerChannel = eliminated
	// Set the start time of the game.
	g.StartAt = time.Now()
	// Create a log message for the start of the game.
	startLog := fmt.Sprintf(
		"%d - starting battle!\n",
		g.StartAt.UnixMicro())
	// Append the log message to the Logs slice.
	g.Logs = append(g.Logs, Log{Description: startLog})
	// Send the log message to the log channel.
	g.BattleLogsChannel <- startLog
	// Start an infinite loop.
	for {
		// Call this function to perform actions
		// like attacking and eliminating players.
		g.performPlayerActions()
		time.Sleep(constant.SleepDuration) // todo: lets think about this
		// Check if there is only one player left alive.
		if g.alivePlayers() == 1 {
			// Calculate the ranks of the players.
			g.calculatePlayersRank()
			// Calculate the points of the players.
			g.calculatePlayersPoint()
			// Send the game data
			game <- g
			// break the loop
			break
		}
	}
}

// performPlayerActions is responsible for executing the moves in the game.
func (g *Game) performPlayerActions() {
	// Iterate over all the players in the game.
	for i, p := range g.Players {
		// If the player eliminated, continue.
		if p.EliminatedAt != nil {
			continue
		}
		// if player not eliminated,
		// Generate a random index for the target player.
		j := rand.Intn(len(g.Players))
		// Make sure the target player is not the same as the current player
		// and that the target player is not eliminated.
		if i != j && g.Players[j].EliminatedAt == nil {
			// Attack the target player.
			attack := p.Attack(g.Players[j])
			// Append the attack log to the game logs.
			g.Logs = append(g.Logs, attack)
			// Send the attack log to the battle logs channel.
			g.BattleLogsChannel <- attack.Description
		}
		// If the target player's health is zero or below
		// and the player has not been eliminated,
		// mark the player as eliminated.
		if g.Players[j].Health <= 0 && g.Players[j].EliminatedAt == nil {
			// eliminate and get 2nd winner
			g.eliminatePlayer(j)
		}
		// If there is only one player left alive, marks the end of a game,
		// sets the winner, and logs information
		if g.alivePlayers() == 1 {
			g.markEnd(i)
		}
	}
}

// alivePlayers checks the number of players who are still alive in the game.
// criteria:
// 1. have a positive Health value,
// 2. have not been eliminated (EliminatedAt is nil)
func (g *Game) alivePlayers() int {
	alivePlayers := 0
	for _, p := range g.Players {
		if p.Health > 0 && p.EliminatedAt == nil {
			alivePlayers++
		}
	}
	return alivePlayers
}

// eliminatePlayer marks a player as eliminated from the game at the current time.
// pos: The index of the player to eliminate in the Game.Players slice.
func (g *Game) eliminatePlayer(pos int) {
	// Eliminate selected player.
	p := g.Players[pos]
	e := p.Eliminate()
	// Append the log message to the Logs slice.
	g.Logs = append(g.Logs, e)
	// Send the log message to the battle logs channel.
	g.BattleLogsChannel <- e.Description
	// Send the player's name to the eliminated players channel.
	g.EliminatedPlayerChannel <- p.Name
}

// calculatePlayersRank calculates the rank for each player in the game.
// Players who stand in the end will get the biggest rank
func (g *Game) calculatePlayersRank() {
	sort.Slice(g.Players, func(i, j int) bool {
		if g.Players[i].EliminatedAt != nil &&
			g.Players[j].EliminatedAt != nil {
			return g.Players[j].
				EliminatedAt.
				Before(*g.Players[i].EliminatedAt)
		}
		return true
	})

	rank := 1
	for _, p := range g.Players {
		if p.EliminatedAt != nil {
			rank += 1
		}
		p.UpdateRank(rank)
	}
}

// calculatePlayersPoint calculates the points for each player based on their rank.
// Players with a higher rank will get more points.
func (g *Game) calculatePlayersPoint() {
	for i, p := range g.Players {
		p.UpdateScore(constant.MaxPoint - i)
	}
}

// markEnd marks the end of a game, sets the winner, and logs information
func (g *Game) markEnd(pos int) {
	// Set the end time of the game
	g.EndAt = time.Now()
	// Set the winner of the game.
	g.Winner = g.Players[pos]
	// send winner log
	logWinner := fmt.Sprintf(
		"%d - %s win the game!\n",
		g.EndAt.UnixMicro(), g.Winner.Name)
	g.Logs = append(g.Logs, Log{Description: logWinner})
	g.BattleLogsChannel <- logWinner
	// send end battle log
	logEnd := fmt.Sprintf("%d - battle end!\n",
		g.EndAt.UnixMicro())
	g.Logs = append(g.Logs, Log{Description: logEnd})
	g.BattleLogsChannel <- logEnd
}

// Reset the game make the struct value fresh again
func (g *Game) Reset() {
	g.Logs = nil
	g.Winner = nil
	g.Players = nil
	g.StartAt = time.Time{}
	g.EndAt = time.Time{}
}

func NewGame(players []*Player) IGameAction {
	return &Game{
		Logs:    nil,
		Winner:  nil,
		Players: players,
		StartAt: time.Time{},
		EndAt:   time.Time{},
	}
}
