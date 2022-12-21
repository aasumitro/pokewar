package battleroyale

import (
	"fmt"
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

	IGameAction interface {
		// Start begins the game loop.
		Start(game chan *Game, log chan string, eliminated chan string)
		// Moves do action like attack eliminate etc.
		Moves()
		// Eliminate marks a player as eliminated from the game at the current time.
		// pos: The index of the player to eliminate in the Game.Players slice.
		Eliminate(pos int)
		// Rank calculates the rank for each player in the game.
		// Players who stand in the end will get the biggest rank
		Rank()
		// Point calculates the points for each player based on their rank.
		// Players with a higher rank will get more points.
		Point()
		// Reset the struct
		Reset()
	}
)

func (g *Game) Start(game chan *Game, log chan string, eliminated chan string) {
	g.BattleLogsChannel = log
	g.EliminatedPlayerChannel = eliminated
	g.StartAt = time.Now()
	startLog := fmt.Sprintf("%d - starting battle!\n", g.StartAt.UnixMicro())
	g.Logs = append(g.Logs, Log{Description: startLog})
	g.BattleLogsChannel <- startLog

	for {
		g.Moves()

		// Check if only one player is left alive.
		alivePlayers := 0
		for _, p := range g.Players {
			if p.Health > 0 && p.EliminatedAt == nil {
				alivePlayers++
			}
		}

		if alivePlayers == 1 {
			// Calculate the ranks of the players.
			g.Rank()
			// Calculate the points of the players.
			g.Point()
			// Send the game data
			game <- g
			// break the loop
			break
		}
	}
}

// Moves TODO: REVALIDATE THIS FUNCTION
func (g *Game) Moves() {
	for i, p := range g.Players {
		if p.Health <= 0 && p.EliminatedAt == nil {
			g.Eliminate(i)
		}

		// If the player is not eliminated, choose a random target to attack.
		if p.EliminatedAt == nil {
			// Generate a random index for the target player.
			j := rand.Intn(len(g.Players))

			// Make sure the target player is not the same as the current player
			// and that the target player is not eliminated.
			if i != j && g.Players[j].EliminatedAt == nil {
				data := p.Attack(g.Players[j])
				g.Logs = append(g.Logs, data)
				g.BattleLogsChannel <- data.Description
			}

			alivePlayers := 0
			for _, p := range g.Players {
				if p.Health > 0 && p.EliminatedAt == nil {
					alivePlayers++
				}
			}

			if alivePlayers == 1 && g.Players[j].EliminatedAt == nil {
				// eliminate and get 2nd winner
				g.Eliminate(j)
				// Set the end time of the game
				g.EndAt = time.Now()
				logEnd := fmt.Sprintf("%d - battle end!\n", g.EndAt.UnixMicro())
				g.Logs = append(g.Logs, Log{Description: logEnd})
				g.BattleLogsChannel <- logEnd
				// Set the winner of the game.
				g.Winner = g.Players[i]
				logWinner := fmt.Sprintf("%d - %s win the game!\n", time.Now().UnixMicro(), g.Winner.Name)
				g.Logs = append(g.Logs, Log{Description: logWinner})
				g.BattleLogsChannel <- logWinner
			}
		}
	}
}

func (g *Game) Eliminate(pos int) {
	player := g.Players[pos]
	now := time.Now()
	player.EliminatedAt = &now
	log := fmt.Sprintf("%d - %s eliminated from the game!\n", time.Now().UnixMicro(), player.Name)
	g.Logs = append(g.Logs, Log{Description: log})
	g.BattleLogsChannel <- log
	g.EliminatedPlayerChannel <- player.Name
}

func (g *Game) Rank() {
	sort.Slice(g.Players, func(i, j int) bool {
		if g.Players[i].EliminatedAt != nil && g.Players[j].EliminatedAt != nil {
			return g.Players[j].
				EliminatedAt.
				Before(*g.Players[i].EliminatedAt)
		}
		return true
	})

	ranks := [5]int{1, 2, 3, 4, 5}
	for i, p := range g.Players {
		p.Rank += ranks[i]
	}
}

func (g *Game) Point() {
	points := [5]int{5, 4, 3, 2, 1}
	for i, p := range g.Players {
		p.Score += points[i]
	}
}

func (g *Game) Reset() {
	g.Logs = nil
	g.Winner = nil
	g.StartAt = time.Time{}
	g.EndAt = time.Time{}
	g.Players = nil
}

func NewGame(players []*Player) IGameAction {
	return &Game{
		Players: players,
		Logs:    nil,
		Winner:  nil,
		StartAt: time.Time{},
		EndAt:   time.Time{},
	}
}
