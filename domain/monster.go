package domain

type (
	MonsterEntity struct {
		ID       int
		OriginID int
		Name     string
		BaseExp  int
		Height   int
		Weight   int
		Avatar   string
		Types    string
		Stats    string
		Skills   string
	}

	Monster struct {
		ID       int      `json:"id"`
		OriginID int      `json:"origin_id"`
		Name     string   `json:"name"`
		BaseExp  int      `json:"base_exp"`
		Height   int      `json:"height"`
		Weight   int      `json:"weight"`
		Avatar   string   `json:"avatar"`
		Types    []string `json:"types"`
		Stats    []Stat   `json:"stats"`
		Skills   []*Skill `json:"skills"`
	}

	Stat struct {
		BaseStat int    `json:"base_stat"`
		Name     string `json:"name"`
	}

	Skill struct {
		PP   int    `json:"pp"` // Power Points
		Name string `json:"name"`
	}

	IMonsterRepository interface {
		ICreateRepository[Monster]
		IReadAllRepository[Monster]
		IReadAllWhereInRepository[Monster]
		IUpdateRepository[Monster]
	}
)
