package domain

type (
	Rank struct {
		ID           int      `json:"id"`
		OriginID     int      `json:"origin_id"`
		Name         string   `json:"name"`
		Avatar       string   `json:"avatar"`
		T            string   `json:"-"`
		Types        []string `json:"types"`
		TotalBattles int      `json:"total_battles"`
		WinBattles   int      `json:"win_battles"`
		LoseBattle   int      `json:"lose_battles"`
		Points       int      `json:"points"`
	}

	IRankRepository interface {
		IReadAllRepository[Rank]
	}
)
