package battleroyale_test

import (
	"github.com/aasumitro/pokewar/pkg/battleroyale"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerAttack(t *testing.T) {
	attacker := &battleroyale.Player{
		ID:     1,
		Name:   "Player 1",
		Health: 100,
		Skills: []*battleroyale.Skill{
			{Name: "Punch", Power: 10},
			{Name: "Kick", Power: 20},
		},
	}
	defender := &battleroyale.Player{
		ID:     2,
		Name:   "Player 2",
		Health: 100,
	}

	log := attacker.Attack(defender)
	assert.NotEqual(t, "", log.Description)
	assert.Contains(t, log.Description, "Player 1 uses")
	assert.Contains(t, log.Description, "to attack Player 2")
}
