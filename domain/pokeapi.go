package domain

type (
	PokemonResult struct {
		Count    int              `json:"count"`
		Next     string           `json:"next"`
		Previous interface{}      `json:"previous"`
		Results  []PokemonSummary `json:"results"`
	}

	PokemonSummary struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	Pokemon struct {
		ID             int     `json:"id"`
		Name           string  `json:"name"`
		Height         int     `json:"height"`
		Weight         int     `json:"weight"`
		BaseExperience int     `json:"base_experience"`
		Sprites        Sprites `json:"sprites"`
		Types          []Types `json:"types"`
		Stats          []Stats `json:"stats"`
		Moves          []Moves `json:"moves"`
	}

	Sprites struct {
		Other Other `json:"other"`
	}

	Other struct {
		DreamWorld DreamWorld `json:"dream_world"`
	}

	DreamWorld struct {
		FrontDefault string `json:"front_default"`
	}

	Types struct {
		Type Type `json:"type"`
	}

	Type struct {
		Name string `json:"name"`
	}

	Stats struct {
		BaseStat int  `json:"base_stat"`
		Stat     Stat `json:"stat"`
	}

	Moves struct {
		Move Move `json:"move"`
	}

	Move struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	IPokeapiRESTRepository interface {
		Pokemon(offset, limit int) ([]*Monster, error)
	}
)
