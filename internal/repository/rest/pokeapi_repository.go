package rest

import (
	"context"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/appconfig"
	"github.com/aasumitro/pokewar/pkg/constant"
	"github.com/aasumitro/pokewar/pkg/httpclient"
	"math/rand"
	"sync"
)

type pokeapiRESTRepository struct {
	client httpclient.IHttpClient
}

// Pokemon retrieves a list of monsters from the PokeAPI REST API.
func (repo *pokeapiRESTRepository) Pokemon(ctx context.Context, offset, limit int) ([]*domain.Monster, error) {
	client := repo.client.NewClient(
		httpclient.Ctx(ctx),
		httpclient.Timeout(constant.TimeoutDuration),
		httpclient.Endpoint(fmt.Sprintf(
			"%spokemon?offset=%d&limit=%d",
			appconfig.Instance.PokeapiURL, offset, limit,
		)))

	monsters, err := ProceedData(client, TransformData)
	if err != nil {
		return nil, err
	}

	return monsters, nil
}

// ProceedData - helper function to proceed data
func ProceedData(
	client *httpclient.HTTPClient,
	transformData func(*httpclient.HTTPClient, *domain.Pokemon) *domain.Monster,
) ([]*domain.Monster, error) {
	var pokemons *domain.PokemonResult
	var monsters []*domain.Monster

	if err := client.MakeRequest(&pokemons); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, pokemon := range pokemons.Results {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			var pokemon *domain.Pokemon
			client.Endpoint = url
			if err := client.MakeRequest(&pokemon); err != nil {
				monsters = append(monsters, nil)
			} else {
				monsters = append(monsters, transformData(client, pokemon))
			}
		}(pokemon.URL)
	}
	wg.Wait()

	return monsters, nil
}

// TransformData function to transform origin data
func TransformData(
	client *httpclient.HTTPClient,
	pokemon *domain.Pokemon,
) *domain.Monster {
	types := make([]string, 0, len(pokemon.Types))
	for _, t := range pokemon.Types {
		types = append(types, t.Type.Name)
	}

	stats := make([]domain.Stat, 0, len(pokemon.Stats))
	for _, s := range pokemon.Stats {
		stats = append(stats, domain.Stat{
			BaseStat: s.BaseStat,
			Name:     s.Stat.Name,
		})
	}

	var skills []*domain.Skill
	var wgMove sync.WaitGroup
	for _, moveURL := range RandomSubset(pokemon.Moves, constant.MaxMoveSize) {
		wgMove.Add(1)
		go func(moveURL string) {
			defer wgMove.Done()
			var skill *domain.Skill
			client.Endpoint = moveURL
			if err := client.MakeRequest(&skill); err != nil {
				skills = append(skills, nil)
			} else {
				skills = append(skills, skill)
			}
		}(moveURL)
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

// RandomSubset helper function to get random skills/moves
func RandomSubset(slice []domain.Moves, size int) []string {
	max := len(slice)
	result := make([]string, 0, size)
	generatedKey := make(map[int]bool)
	for len(result) < size {
		n := rand.Intn(max)
		if _, found := generatedKey[n]; !found {
			result = append(result, slice[n].Move.URL)
			generatedKey[n] = true
		}
	}
	return result
}

// NewPokeapiRESTRepository use in main app
func NewPokeapiRESTRepository() domain.IPokeapiRESTRepository {
	return &pokeapiRESTRepository{client: &httpclient.HTTPClient{}}
}
