package domain

type (
	Rank struct {
		ID           int
		OriginID     int
		Name         string
		Avatar       string
		T            string   `json:"-"`
		Types        []string `json:"types"`
		TotalBattles int
		WinBattles   int
		LoseBattle   int
		Points       int
	}

	IRankRepository interface {
		IReadAllRepository[Rank]
	}
)
