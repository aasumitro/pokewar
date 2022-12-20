package rest

import (
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/aasumitro/pokewar/pkg/httpclient"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type pokeapiRESTRepository struct {
	httpClient httpclient.IHttpClient
}

// Pokemon
// TODO: Optimize
// BENCH RESULT wg is better than ch
// STILL TO SLOW SOMETIMES TOOK 3s TO LOAD AND PROCEED THE DATA
func (repo *pokeapiRESTRepository) Pokemon(offset int) ([]*domain.Monster, error) {
	var pokemons *domain.PokemonResult
	var wg sync.WaitGroup
	var monsters []*domain.Monster

	client := &httpclient.HttpClient{
		Endpoint: fmt.Sprintf(
			"%s/pokemon?offset=%d&limit=25",
			appconfigs.Instance.PokeapiUrl, offset,
		),
		Timeout: 10 * time.Second,
		Method:  http.MethodGet,
	}
	if err := client.MakeRequest(&pokemons); err != nil {
		return nil, err
	}
	httpClientTest(repo) // FOR TEST PURPOSE (MOCK)

	for _, pokemon := range pokemons.Results {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			var pokemon *domain.Pokemon

			client = &httpclient.HttpClient{
				Endpoint: url,
				Timeout:  10 * time.Second,
				Method:   http.MethodGet,
			}
			err := client.MakeRequest(&pokemon)
			httpClientTest(repo) // FOR TEST PURPOSE (MOCK)

			if err != nil {
				monsters = append(monsters, nil)
			} else {
				monsters = append(monsters, transformData(repo, pokemon))
			}
		}(pokemon.URL)
	}

	wg.Wait()

	return monsters, nil
}

// helper function to transform origin data
func transformData(repo *pokeapiRESTRepository, pokemon *domain.Pokemon) *domain.Monster {
	types := make([]string, 0, len(pokemon.Types))
	for _, pokemon := range pokemon.Types {
		types = append(types, pokemon.Type.Name)
	}

	stats := make([]domain.Stat, 0, len(pokemon.Stats))
	for _, pokemon := range pokemon.Stats {
		stats = append(stats, domain.Stat{
			BaseStat: pokemon.BaseStat,
			Name:     pokemon.Stat.Name,
		})
	}

	max := len(pokemon.Moves)
	movesUrls := make([]string, 0, 4)
	generatedKey := make(map[int]bool)
	for len(movesUrls) < 4 {
		n := rand.Intn(max)
		if _, found := generatedKey[n]; !found {
			movesUrls = append(movesUrls, pokemon.Moves[n].Move.Url)
			generatedKey[n] = true
		}
	}

	var wgMove sync.WaitGroup
	var skills []*domain.Skill
	for _, moveUrl := range movesUrls {
		wgMove.Add(1)
		go func(moveUrl string) {
			defer wgMove.Done()
			var skill *domain.Skill
			client := httpclient.HttpClient{
				Endpoint: moveUrl,
				Timeout:  10 * time.Second,
				Method:   http.MethodGet,
			}
			if err := client.MakeRequest(&skill); err != nil {
				skills = append(skills, nil)
			} else {
				skills = append(skills, skill)
			}
			httpClientTest(repo) // FOR TEST PURPOSE (MOCK)
		}(moveUrl)
	}
	wgMove.Wait()

	return &domain.Monster{
		OriginID: pokemon.ID,
		Name:     pokemon.Name,
		BaseExp:  pokemon.BaseExperience,
		Height:   pokemon.Height,
		Weight:   pokemon.Weight,
		Avatar:   pokemon.Sprites.Other.DreamWorld.FrontDefault,
		Types:    types,
		Stats:    stats,
		Skills:   skills,
	}
}

// NewPokeapiRESTRepository use in main app
func NewPokeapiRESTRepository() domain.IPokeapiRESTRepository {
	return &pokeapiRESTRepository{}
}

// helper function to help unit test
func httpClientTest(repo *pokeapiRESTRepository) {
	if repo.httpClient != nil {
		var obj interface{}
		_ = repo.httpClient.MakeRequest(obj)
	}
}

// NewPokeapiRESTRepositoryTest - use for tests with httpclient mock (inject)
func NewPokeapiRESTRepositoryTest(httpclient httpclient.IHttpClient) domain.IPokeapiRESTRepository {
	return &pokeapiRESTRepository{httpClient: httpclient}
}
