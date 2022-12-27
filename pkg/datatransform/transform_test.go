package datatransform_test

import (
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/battleroyale"
	"github.com/aasumitro/pokewar/pkg/datatransform"
	"reflect"
	"testing"
	"time"
)

func TestTransformMonstersAsGamePlayers(t *testing.T) {
	monsters := []*domain.Monster{
		{
			ID:   1,
			Name: "Charizard",
			Stats: []domain.Stat{
				{Name: "hp", BaseStat: 78},
				{Name: "attack", BaseStat: 84},
				{Name: "defense", BaseStat: 78},
			},
			Skills: []*domain.Skill{
				{Name: "Flame Thrower", PP: 15},
				{Name: "Dragon Claw", PP: 15},
				{Name: "Fire Blast", PP: 5},
				{Name: "Blast Burn", PP: 5},
			},
		},
		{
			ID:   2,
			Name: "Blastoise",
			Stats: []domain.Stat{
				{Name: "hp", BaseStat: 79},
				{Name: "attack", BaseStat: 83},
				{Name: "defense", BaseStat: 100},
			},
			Skills: []*domain.Skill{
				{Name: "Water Pulse", PP: 20},
				{Name: "Hydro Pump", PP: 5},
				{Name: "Surf", PP: 15},
				{Name: "Aqua Tail", PP: 10},
			},
		},
	}
	players := []*battleroyale.Player{
		{
			ID:     1,
			Name:   "Charizard",
			Health: 78,
			Score:  0,
			Rank:   0,
			Skills: []*battleroyale.Skill{
				{Power: 15, Name: "Flame Thrower"},
				{Power: 15, Name: "Dragon Claw"},
				{Power: 5, Name: "Fire Blast"},
				{Power: 5, Name: "Blast Burn"},
			},
		},
		{
			ID:     2,
			Name:   "Blastoise",
			Health: 79,
			Score:  0,
			Rank:   0,
			Skills: []*battleroyale.Skill{
				{Power: 20, Name: "Water Pulse"},
				{Power: 5, Name: "Hydro Pump"},
				{Power: 15, Name: "Surf"},
				{Power: 10, Name: "Aqua Tail"},
			},
		},
	}

	result := datatransform.TransformMonstersAsGamePlayers(monsters)

	if !reflect.DeepEqual(result, players) {
		t.Errorf("expected %v, got %v", result, players)
	}
}

func TestTransformGameResultToBattle(t *testing.T) {
	gamesResult := &battleroyale.Game{
		StartAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		EndAt:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
		Logs: []battleroyale.Log{
			{Description: "log 1"},
			{Description: "log 2"},
		},
		Players: []*battleroyale.Player{
			{
				ID:           1,
				Name:         "player 1",
				Health:       100,
				Score:        10,
				Rank:         1,
				EliminatedAt: &time.Time{},
			},
			{
				ID:     2,
				Name:   "player 2",
				Health: 50,
				Score:  5,
				Rank:   2,
			},
		},
	}
	expectedBattle := &domain.Battle{
		StartedAt: gamesResult.StartAt.UnixMicro(),
		EndedAt:   gamesResult.EndAt.UnixMicro(),
		Logs: []domain.Log{
			{Description: "log 1"},
			{Description: "log 2"},
		},
		Players: []domain.Player{
			{
				MonsterID: 1,
				Name:      "player 1",
				EliminatedAt: func() int64 {
					if gamesResult.Players[0].EliminatedAt != nil {
						return gamesResult.Players[0].EliminatedAt.UnixMicro()
					}
					return 0
				}(),
				Rank:  1,
				Point: 10,
			},
			{
				MonsterID: 2,
				Name:      "player 2",
				EliminatedAt: func() int64 {
					if gamesResult.Players[1].EliminatedAt != nil {
						return gamesResult.Players[1].EliminatedAt.UnixMicro()
					}
					return 0
				}(),
				Rank:  2,
				Point: 5,
			},
		},
	}

	battleResult := datatransform.TransformGameResultToBattle(gamesResult)

	if !reflect.DeepEqual(battleResult, expectedBattle) {
		t.Errorf("expected %v, got %v", battleResult, expectedBattle)
	}
}
