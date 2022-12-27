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
		Eliminate() Log
		UpdateRank(rank int)
		UpdateScore(point int)
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
			"%d - %s uses %s to attack %s, reducing %s health to %d\n",
			time.Now().UnixMicro(), p.Name, skill.Name, other.Name, other.Name, other.Health,
		),
	}
}

func (p *Player) Eliminate() Log {
	eliminated := time.Now()
	p.EliminatedAt = &eliminated

	return Log{
		Description: fmt.Sprintf("%d - %s eliminated from the game!\n",
			p.EliminatedAt.UnixMicro(), p.Name),
	}
}

func (p *Player) UpdateRank(rank int) {
	p.Rank += rank
}

func (p *Player) UpdateScore(point int) {
	p.Score += point
}
