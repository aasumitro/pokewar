package datatransform

import (
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/battleroyale"
)

// TransformMonstersAsGamePlayers transforms a list of monsters into a list of players
// by extracting relevant information from the monsters.
func TransformMonstersAsGamePlayers(monsters []*domain.Monster) []*battleroyale.Player {
	// Pre-allocate capacity for the players slice using make
	players := make([]*battleroyale.Player, len(monsters))
	// Iterate through each monster and extract the necessary information to create a player
	for i, monster := range monsters {
		// Find the monster's HP stat
		var hp int
		for _, stat := range monster.Stats {
			if stat.Name == "hp" {
				hp = stat.BaseStat
				break
			}
		}
		// Create a player using the extracted information
		players[i] = &battleroyale.Player{
			ID:     monster.ID,
			Name:   monster.Name,
			Health: hp,
			Score:  0,
			Rank:   0,
			Skills: []*battleroyale.Skill{
				{Power: monster.Skills[0].PP, Name: monster.Skills[0].Name},
				{Power: monster.Skills[1].PP, Name: monster.Skills[1].Name},
				{Power: monster.Skills[2].PP, Name: monster.Skills[2].Name},
				{Power: monster.Skills[3].PP, Name: monster.Skills[3].Name},
			},
		}
	}
	// Save the list of players to the GamePlayers map
	return players
}

// TransformGameResultToBattle converts the data from the battleroyale.Game struct to the domain.Battle struct format.
func TransformGameResultToBattle(game *battleroyale.Game) *domain.Battle {
	// pre-allocates the slices for logs sing the length of game.Logs
	logs := make([]domain.Log, len(game.Logs))
	for i, log := range game.Logs {
		// populates and transform data.
		logs[i] = domain.Log{Description: log.Description}
	}
	// pre-allocates the slices for logs sing the length of game.Players
	players := make([]domain.Player, len(game.Players))
	for i, player := range game.Players {
		// populates and transform data.
		players[i] = domain.Player{
			MonsterID: player.ID,
			Name:      player.Name,
			EliminatedAt: func() int64 {
				if player.EliminatedAt != nil {
					return player.EliminatedAt.UnixMicro()
				}
				return 0
			}(),
			Rank:  player.Rank,
			Point: player.Score,
		}
	}
	// returns the transformed data.
	return &domain.Battle{
		StartedAt: (*game).StartAt.UnixMicro(),
		EndedAt:   (*game).EndAt.UnixMicro(),
		Players:   players,
		Logs:      logs,
	}
}
