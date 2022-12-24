package rest_test

import (
	"encoding/json"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/internal/repository/rest"
	"github.com/aasumitro/pokewar/mocks"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/aasumitro/pokewar/pkg/httpclient"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type pokeapiRESTRepositoryTestSuite struct {
	suite.Suite
	moves   []domain.Moves
	pokemon domain.Pokemon
	monster domain.Monster
}

func (suite *pokeapiRESTRepositoryTestSuite) SetupSuite() {
	suite.moves = []domain.Moves{
		{Move: domain.Move{Name: "a1", Url: "http://example.com/a1"}},
		{Move: domain.Move{Name: "a3", Url: "http://example.com/a3"}},
		{Move: domain.Move{Name: "a4", Url: "http://example.com/a4"}},
		{Move: domain.Move{Name: "a6", Url: "http://example.com/a6"}},
	}
	suite.pokemon = domain.Pokemon{
		ID:             1,
		Name:           "lorem",
		Height:         10,
		Weight:         10,
		BaseExperience: 10,
		Sprites: domain.Sprites{
			Other: domain.Other{
				DreamWorld: domain.DreamWorld{
					FrontDefault: "http://example.com/image.png"}}},
		Types: []domain.Types{
			{Type: domain.Type{Name: "asd"}},
			{Type: domain.Type{Name: "qwe"}},
		},
		Stats: []domain.Stats{
			{BaseStat: 10, Stat: domain.Stat{Name: "asd"}},
			{BaseStat: 5, Stat: domain.Stat{Name: "qwe"}},
			{BaseStat: 1, Stat: domain.Stat{Name: "zxc"}},
		},
		Moves: suite.moves,
	}
	suite.monster = domain.Monster{
		OriginID: suite.pokemon.ID,
		Name:     suite.pokemon.Name,
		BaseExp:  suite.pokemon.BaseExperience,
		Height:   suite.pokemon.Height,
		Weight:   suite.pokemon.Weight,
		Avatar:   suite.pokemon.Sprites.Other.DreamWorld.FrontDefault,
		Types:    []string{"asd", "qwe"},
		Stats: []domain.Stat{
			{BaseStat: 10, Name: "asd"},
			{BaseStat: 5, Name: "qwe"},
			{BaseStat: 1, Name: "zxc"}},
		Skills: []*domain.Skill{nil, nil, nil, nil},
	}
}

func (suite *pokeapiRESTRepositoryTestSuite) TestRepository_Pokemon() {
	viper.SetConfigFile("../../../.example.env")
	appconfigs.LoadEnv()
	appconfigs.Instance.PokeapiUrl = "https://pokeapi.co/api/v2/"
	data, _ := json.Marshal(domain.PokemonResult{
		Results: []domain.PokemonSummary{
			{Name: "bulbasaur", URL: "https://pokeapi.co/api/v2/pokemon/1/"},
		},
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(data)
	}))
	defer server.Close()
	repo := rest.NewPokeapiRESTRepository()
	monsters, err := repo.Pokemon(0, 1)
	require.NotNil(suite.T(), monsters)
	require.Nil(suite.T(), err)
}

func (suite *pokeapiRESTRepositoryTestSuite) TestRepository_ProceedData() {
	m := new(mocks.IHttpClient)
	m.On("NewClient").
		Return(&httpclient.HttpClient{
			Timeout: 10 * time.Second,
			Method:  http.MethodGet,
		}).Once()
	c := m.NewClient()
	c.Endpoint = "https://pokeapi.co/api/v2/pokemon/1/"
	m.On("MakeRequest", mock.Anything).
		Return(domain.PokemonResult{
			Results: []domain.PokemonSummary{
				{Name: "bulbasaur", URL: "https://pokeapi.co/api/v2/pokemon/1/"},
			},
		}).Once()
	monster, err := rest.ProceedData(c, rest.TransformData)
	require.Nil(suite.T(), monster)
	require.Nil(suite.T(), err)
}

func (suite *pokeapiRESTRepositoryTestSuite) TestRepository_TransformData() {
	client := &httpclient.HttpClient{
		Timeout: 10 * time.Second,
		Method:  http.MethodGet,
	}
	monster := rest.TransformData(client, &suite.pokemon)
	require.NotNil(suite.T(), monster)
	require.Equal(suite.T(), monster, &suite.monster)
}

func (suite *pokeapiRESTRepositoryTestSuite) TestRepository_RandomSubset() {
	result := rest.RandomSubset(suite.moves, 4)
	require.NotNil(suite.T(), result)
	require.Contains(suite.T(), result[0], "http://example.com/a")
}

func TestPokeapiRESTRepository(t *testing.T) {
	suite.Run(t, new(pokeapiRESTRepositoryTestSuite))
}
