package service

import (
	"context"
	"fmt"
	"github.com/aasumitro/pokewar/configs"
	"github.com/aasumitro/pokewar/constants"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/utils"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

type pokewarService struct {
	ctx         context.Context
	pokemonRepo domain.IPokeapiRESTRepository
	monsterRepo domain.IMonsterRepository
	rankRepo    domain.IRankRepository
	battleRepo  domain.IBattleRepository
}

// MonstersCount returns the total count of monsters stored in the database.
func (service *pokewarService) MonstersCount() int {
	return service.monsterRepo.Count(service.ctx)
}

// FetchMonsters retrieves a list of monsters from the database.
// The list can be filtered by providing arguments such as limit and offset.
// The function returns a slice of pointers to domain.Monster structs and an error.
func (service *pokewarService) FetchMonsters(
	args ...string,
) (monsters []*domain.Monster, error *utils.ServiceError) {
	data, err := service.monsterRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Monster](data, err)
}

// SyncMonsters synchronizes the data for the monsters in the local database
// with the data from the remote API (https://pokeapi.co).
func (service *pokewarService) SyncMonsters(
	updateEnv bool,
	_ ...string,
) (data []*domain.Monster, svcErr *utils.ServiceError) {
	offset := configs.Instance.TotalMonsterSync
	limit := configs.Instance.LimitSync
	lastID := configs.Instance.LastMonsterID
	maxRetries, retryCount := 3, 0
	// get data from pokeapi.co
	data, err := service.pokemonRepo.Pokemon(service.ctx, offset, limit)
	if err != nil {
		return nil, &utils.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	// get highest id from result
	sort.Slice(data, func(i, j int) bool {
		return data[i].OriginID > data[j].OriginID
	})
	maxID := data[0].OriginID
	// store data from pokeapi.co
	if maxID > lastID {
		for {
			err = service.monsterRepo.Create(service.ctx, data)
			// Data was successfully inserted,
			// so break out of the loop
			if err == nil {
				break
			}
			retryCount++
			// If the maximum number of retries is reached,
			// break out of the loop
			if retryCount >= maxRetries {
				break
			}
			// Data was not successfully inserted, so sleep for the specified delay before trying again
			time.Sleep(constants.SleepDuration)
		}
	}
	// when success update env
	if updateEnv && err == nil {
		configs.Instance.UpdateEnv("LAST_SYNC", time.Now().Unix())
		configs.Instance.UpdateEnv("TOTAL_MONSTER_SYNC", offset+len(data))
		configs.Instance.UpdateEnv("LAST_MONSTER_ID", maxID)
	}
	// return data to handler
	return utils.ValidateDataRows[domain.Monster](data, err)
}

// FetchRanks retrieves a list of monsters rank from the database.
// The list can be filtered by providing arguments such as limit and offset.
// The function returns a slice of pointers to domain.Rank structs and an error.
func (service *pokewarService) FetchRanks(
	args ...string,
) (ranks []*domain.Rank, error *utils.ServiceError) {
	data, err := service.rankRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Rank](data, err)
}

// BattlesCount returns the total count of battles stored in the database.
func (service *pokewarService) BattlesCount() int {
	return service.battleRepo.Count(service.ctx)
}

// FetchBattles retrieves a list of battles from the database.
// The list can be filtered by providing arguments such as limit and offset.
// also between with unix timestamp in millisecond format
// The function returns a slice of pointers to domain.Battle structs and an error.
func (service *pokewarService) FetchBattles(
	args ...string,
) (ranks []*domain.Battle, error *utils.ServiceError) {
	data, err := service.battleRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Battle](data, err)
}

// PrepareMonstersForBattle retrieves a list of monsters from the database limited to 5.
func (service *pokewarService) PrepareMonstersForBattle() (
	monsters []*domain.Monster,
	error *utils.ServiceError,
) {
	var args []string
	randID := make([]int, 0, constants.MaxPlayerSize)
	generatedKey := make(map[int]bool)
	for len(randID) < 5 {
		n := rand.Intn(configs.Instance.TotalMonsterSync-1) + 1
		if _, found := generatedKey[n]; !found {
			randID = append(randID, n)
			generatedKey[n] = true
		}
	}
	args = append(args, fmt.Sprintf(
		"WHERE origin_id IN (%d,%d,%d,%d,%d)",
		randID[0], randID[1], randID[2], randID[3], randID[4]))
	args = append(args, "LIMIT 5")

	data, err := service.monsterRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Monster](data, err)
}

// AddBattle store/insert the latest match to database
func (service *pokewarService) AddBattle(param *domain.Battle) *utils.ServiceError {
	err := service.battleRepo.Create(service.ctx, param)
	if err != nil {
		return &utils.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func NewPokewarService(
	ctx context.Context,
	pokemonRepo domain.IPokeapiRESTRepository,
	monsterRepo domain.IMonsterRepository,
	rankRepo domain.IRankRepository,
	battleRepo domain.IBattleRepository,
) domain.IPokewarService {
	return &pokewarService{
		ctx:         ctx,
		pokemonRepo: pokemonRepo,
		monsterRepo: monsterRepo,
		rankRepo:    rankRepo,
		battleRepo:  battleRepo,
	}
}
