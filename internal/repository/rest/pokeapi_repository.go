package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/configs"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type pokeapiRESTRepository struct{}

// Pokemon
// TODO: Optimize
func (repo *pokeapiRESTRepository) Pokemon() ([]*domain.Monster, error) {
	var pokemons *domain.PokemonResult
	var wg sync.WaitGroup
	var monsters []*domain.Monster

	if err := repo.MakeRequest(fmt.Sprintf(
		"%s/pokemon?offiset=50&limit=50",
		configs.Instance.PokeapiUrl,
	), &pokemons); err != nil {
		return nil, err
	}

	for _, pokemon := range pokemons.Results {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			var pokemon *domain.Pokemon
			err := repo.MakeRequest(url, &pokemon)
			if err != nil {
				monsters = append(monsters, nil)
			} else {
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
						err := repo.MakeRequest(moveUrl, &skill)
						if err != nil {
							skills = append(skills, nil)
						} else {
							skills = append(skills, skill)
						}
					}(moveUrl)
				}
				wgMove.Wait()

				monsters = append(monsters, &domain.Monster{
					OriginID: pokemon.ID,
					Name:     pokemon.Name,
					BaseExp:  pokemon.BaseExperience,
					Height:   pokemon.Height,
					Weight:   pokemon.Weight,
					Avatar:   pokemon.Sprites.Other.DreamWorld.FrontDefault,
					Types:    types,
					Stats:    stats,
					Skills:   skills,
				})
			}
		}(pokemon.URL)
	}

	wg.Wait()

	return monsters, nil
}

func (repo *pokeapiRESTRepository) MakeRequest(endpoint string, obj interface{}) error {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return err
	}
	data := buf.Bytes()

	return json.Unmarshal(data, &obj)
}

func NewPokeapiRESTRepository() domain.PokeapiRESTRepository {
	return &pokeapiRESTRepository{}
}
