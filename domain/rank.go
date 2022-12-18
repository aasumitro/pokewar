package domain

type (
	Rank struct {
		ID           int
		OriginID     int
		Name         string
		Avatar       string
		Types        string
		TotalBattles int
		WinBattles   int
		LoseBattle   int
		Points       int
	}

	IRankRepository interface {
		IReadAllRepository[Rank]
	}
)
