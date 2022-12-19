package domain

import "context"

type (
	Battle struct {
		ID        int
		StartedAt int
		EndedAt   int
		Players   []Player
		Logs      []Log
	}

	Player struct {
		ID           int
		BattleID     int
		MonsterID    int
		EliminatedAt int
		AnnulledAt   int
		Rank         int
		Point        int
	}

	Log struct {
		ID          int
		BattleID    int
		Description string
		CreatedAt   int
	}

	IBattleRepository interface {
		ICreateRepository[Battle]
		IReadAllRepository[Battle]
		ICountRowRepository
		UpdatePlayer(ctx context.Context, id int) (annulledTime int64, err error)
	}
)
