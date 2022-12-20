package battleroyale_test

import (
	"github.com/aasumitro/pokewar/pkg/battleroyale"
	"testing"
)

func TestPlayerAttack(t *testing.T) {
	attacker := &battleroyale.Player{
		ID:     1,
		Name:   "Player 1",
		Health: 100,
		Skills: []*battleroyale.Skill{
			{Name: "Punch", Power: 10},
			// ADD MORE HMM #LOL
		},
	}
	defender := &battleroyale.Player{
		ID:     2,
		Name:   "Player 2",
		Health: 100,
	}

	log := attacker.Attack(defender)
	if log.Description != "Player 1 uses Punch to attack Player 2, reducing their health to 90\n" {
		t.Errorf("unexpected attack log: %s", log.Description)
	}
}
