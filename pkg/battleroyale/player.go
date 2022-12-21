package battleroyale

import (
	"fmt"
	"math/rand"
	"time"
)

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

func (p *Player) Attack(other *Player) Log {
	// Choose a random skill to use in the attack.
	var skill *Skill
	if len(p.Skills) > 0 {
		idx := rand.Intn(len(p.Skills))
		skill = p.Skills[idx]
	}

	other.Health -= skill.Power

	return Log{
		Description: fmt.Sprintf(
			"%d - %s uses %s to attack %s, reducing their health to %d\n",
			time.Now().UnixMicro(), p.Name, skill.Name, other.Name, other.Health,
		),
	}
}
