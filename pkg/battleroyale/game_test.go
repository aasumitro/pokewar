package battleroyale_test

import (
	"github.com/aasumitro/pokewar/pkg/battleroyale"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGameStart(t *testing.T) {
	players := []*battleroyale.Player{
		{ID: 1, Name: "Player 1", Health: 100, Skills: []*battleroyale.Skill{
			{Name: "Kick", Power: 20},
			{Name: "Punch", Power: 10},
		}},
		{ID: 2, Name: "Player 2", Health: 100, Skills: []*battleroyale.Skill{
			{Name: "Kick", Power: 20},
			{Name: "Punch", Power: 10},
		}},
		{ID: 3, Name: "Player 3", Health: 100, Skills: []*battleroyale.Skill{
			{Name: "Kick", Power: 20},
			{Name: "Punch", Power: 10},
		}},
	}

	game := &battleroyale.Game{Players: players}
	result := make(chan *battleroyale.Game)
	battleLogs := make(chan string, 100)
	eliminatedPlayers := make(chan string, 3)
	go game.Start(result, battleLogs, eliminatedPlayers)
	data := <-result
	time.Sleep(1 * time.Second)

	assert.NotEqual(t, data.StartAt.IsZero(), true)
	assert.NotEqual(t, data.EndAt.IsZero(), true)
	assert.NotZero(t, len(data.Logs))
	assert.NotNil(t, data.Winner)
	assert.Equal(t, 3, len(data.Players))
	game.Reset()
	battleroyale.NewGame(nil)
}
