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

type pokeapiRESTRepository struct{}

// Pokemon retrieves a list of monsters from the PokeAPI REST API.
func (repo *pokeapiRESTRepository) Pokemon(offset, limit int) ([]*domain.Monster, error) {
	client := &httpclient.HttpClient{
		Endpoint: fmt.Sprintf(
			"%s/pokemon?offset=%d&limit=%d",
			appconfigs.Instance.PokeapiUrl, offset, limit,
		),
		Timeout: 10 * time.Second,
		Method:  http.MethodGet,
	}

	monsters, err := proceedData(client, transformData)
	if err != nil {
		return nil, err
	}

	return monsters, nil
}

// helper function to proceed data
func proceedData(
	client *httpclient.HttpClient,
	transformData func(*httpclient.HttpClient, *domain.Pokemon) *domain.Monster,
) ([]*domain.Monster, error) {
	var pokemons *domain.PokemonResult
	var wg sync.WaitGroup
	var monsters []*domain.Monster

	if err := client.MakeRequest(&pokemons); err != nil {
		return nil, err
	}

	for _, pokemon := range pokemons.Results {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			var pokemon *domain.Pokemon

			client.Endpoint = url
			err := client.MakeRequest(&pokemon)

			if err != nil {
				monsters = append(monsters, nil)
			} else {
				monsters = append(monsters, transformData(client, pokemon))
			}
		}(pokemon.URL)
	}

	wg.Wait()

	return monsters, nil
}

// helper function to transform origin data
func transformData(
	client *httpclient.HttpClient,
	pokemon *domain.Pokemon,
) *domain.Monster {
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

	movesUrls := randomSubset(pokemon.Moves, 4)
	var wgMove sync.WaitGroup
	var skills []*domain.Skill
	for _, moveUrl := range movesUrls {
		wgMove.Add(1)
		go func(moveUrl string) {
			defer wgMove.Done()
			var skill *domain.Skill
			client.Endpoint = moveUrl
			if err := client.MakeRequest(&skill); err != nil {
				skills = append(skills, nil)
			} else {
				skills = append(skills, skill)
			}
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

// helper function to get random skills/moves
func randomSubset(slice []domain.Moves, size int) []string {
	max := len(slice)
	result := make([]string, 0, 4)
	generatedKey := make(map[int]bool)
	for len(result) < size {
		n := rand.Intn(max)
		if _, found := generatedKey[n]; !found {
			result = append(result, slice[n].Move.Url)
			generatedKey[n] = true
		}
	}
	return result
}

// NewPokeapiRESTRepository use in main app
func NewPokeapiRESTRepository() domain.IPokeapiRESTRepository {
	return &pokeapiRESTRepository{}
}
