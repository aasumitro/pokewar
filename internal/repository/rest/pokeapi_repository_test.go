package rest_test

import (
	"encoding/json"
	"github.com/aasumitro/pokewar/domain"
	"github.com/stretchr/testify/suite"
	"testing"
)

type pokeapiRESTRepositoryTestSuite struct {
	suite.Suite
	mock func(url string, v interface{}) error
}

func (suite *pokeapiRESTRepositoryTestSuite) SetupSuite() {
	//goland:noinspection ALL
	suite.mock = func(url string, v interface{}) error {
		switch url {
		case "http://example.com/pokemon?offset=50&limit=50":
			result := &domain.PokemonResult{
				Results: []domain.PokemonSummary{
					{Name: "Bulbasaur", URL: "http://example.com/pokemon/1"},
				},
			}
			b, _ := json.Marshal(result)
			json.Unmarshal(b, v)
			return nil
		case "http://example.com/pokemon/1":
			p := &domain.Pokemon{
				ID:             1,
				Name:           "Bulbasaur",
				BaseExperience: 63,
				Height:         7,
				Weight:         69,
				Sprites: domain.Sprites{
					Other: domain.Other{
						DreamWorld: domain.DreamWorld{
							FrontDefault: "http://example.com/sprites/bulbasaur.png",
						},
					},
				},
				Types: []domain.Types{
					{Type: domain.Type{Name: "grass"}},
					{Type: domain.Type{Name: "poison"}},
				},
				Stats: []domain.Stats{
					{Stat: domain.Stat{Name: "hp"}, BaseStat: 45},
					{Stat: domain.Stat{Name: "attack"}, BaseStat: 49},
					{Stat: domain.Stat{Name: "defense"}, BaseStat: 49},
					{Stat: domain.Stat{Name: "special-attack"}, BaseStat: 65},
					{Stat: domain.Stat{Name: "special-defense"}, BaseStat: 65},
					{Stat: domain.Stat{Name: "speed"}, BaseStat: 45},
				},
				Moves: []domain.Moves{
					{Move: domain.Move{Name: "tackle", Url: "http://example.com/move/1"}},
					{Move: domain.Move{Name: "vine-whip", Url: "http://example.com/move/2"}},
					{Move: domain.Move{Name: "razor-leaf", Url: "http://example.com/move/3"}},
					{Move: domain.Move{Name: "sludge-bomb", Url: "http://example.com/move/4"}},
				},
			}
			b, _ := json.Marshal(p)
			json.Unmarshal(b, v)
			return nil
		}
		return nil
	}
}

// ============
// TODO: HERE
// ============

func TestPokeapiRESTRepository(t *testing.T) {
	suite.Run(t, new(pokeapiRESTRepositoryTestSuite))
}
