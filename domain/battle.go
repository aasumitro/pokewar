package domain

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
		//TODO
		// Create - Create battle result (battle, players, logs)
		// Update - Annulled player?
		// Read - Read all battle with monster rank
	}
)
