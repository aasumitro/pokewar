package domain

type (
	BattleEntity struct {
		ID        int
		StartedAt int64
		EndedAt   int64
		Players   string
		Logs      string
	}

	Battle struct {
		ID        int      `json:"id"`
		StartedAt int64    `json:"started_at"`
		EndedAt   int64    `json:"ended_at"`
		Players   []Player `json:"players"`
		Logs      []Log    `json:"logs"`
	}

	Player struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Avatar       string `json:"avatar"`
		BattleID     int    `json:"battle_id"`
		MonsterID    int    `json:"monster_id"`
		EliminatedAt int64  `json:"eliminated_at"`
		AnnulledAt   int64  `json:"annulled_at"`
		Rank         int    `json:"rank"`
		Point        int    `json:"point"`
	}

	Log struct {
		ID          int    `json:"id"`
		BattleID    int    `json:"battle_id"`
		Description string `json:"description"`
		CreatedAt   int64  `json:"created_at"`
	}

	IBattleRepository interface {
		ICreateRepository[Battle]
		IReadAllRepository[Battle]
		ICountRowRepository
	}
)
