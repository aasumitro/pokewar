package battleroyale

import "time"

type (
	// Skill represents a skill that a player can use in the game.
	Skill struct {
		Name  string
		Power int
	}

	// Player represents a player in the game.
	Player struct {
		ID           int
		Name         string
		Health       int
		Score        int
		Rank         int
		Skills       []*Skill
		EliminatedAt *time.Time
	}

	IPlayerAction interface {
		Attack(other *Player) Log
	}
)
