package battleroyale_test

import (
	"github.com/aasumitro/pokewar/pkg/battleroyale"
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

	if data.StartAt.IsZero() {
		t.Errorf("expected StartAt to be set, but it is zero")
	}
	if data.Winner == nil {
		t.Errorf("expected a winner, but got nil")
	}
	if data.EndAt.IsZero() {
		t.Errorf("expected EndAt to be set, but it is zero")
	}
	if len(data.Logs) == 0 {
		t.Errorf("expected at least one log, but got none")
	}
	if len(data.Players) != 3 {
		t.Errorf("expected 3 players, but got %d", len(data.Players))
	}

	game.Reset()
	battleroyale.NewGame(nil)
}
