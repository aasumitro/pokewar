package battleroyale

import (
	"context"
	"time"
)

type (
	// Game represents the state of a battle royal game.
	Game struct {
		ctx     context.Context
		Players []*Player
		Winner  *Player
		StartAt time.Time
		EndAt   time.Time
		Logs    []Log
	}

	Log struct {
		PlayerID    int
		Description string
	}

	IGameAction interface {
		// Start begins the game loop.
		Start()
		// Moves do action like attack eliminate etc.
		Moves()
		// Eliminate marks a player as eliminated from the game at the current time.
		// pos: The index of the player to eliminate in the Game.Players slice.
		Eliminate()
		// Rank calculates the rank for each player in the game.
		// Players who stand in the end will get the biggest rank
		Rank()
		// Point calculates the points for each player based on their rank.
		// Players with a higher rank will get more points.
		Point()
		// Result return the result of the game
		Result()
	}
)
