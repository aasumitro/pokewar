package service

import (
	"context"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/aasumitro/pokewar/pkg/utils"
	"math/rand"
	"net/http"
	"time"
)

type pokewarService struct {
	ctx         context.Context
	pokemonRepo domain.IPokeapiRESTRepository
	monsterRepo domain.IMonsterRepository
	rankRepo    domain.IRankRepository
	battleRepo  domain.IBattleRepository
}

func (service *pokewarService) MonstersCount() int {
	return service.monsterRepo.Count(service.ctx)
}

func (service *pokewarService) FetchMonsters(args ...string) (monsters []*domain.Monster, error *utils.ServiceError) {
	data, err := service.monsterRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Monster](data, err)
}

// SyncMonsters
//
//	if lastID == maxID || maxID < lastID {
//		go func() {
//			for _, d := range data {
//				if err := service.monsterRepo.Update(service.ctx, d); err != nil {
//					fmt.Println(err.Error())
//				}
//			}
//			done <- true
//		}()
//	}
func (service *pokewarService) SyncMonsters(updateEnv bool, _ ...string) (data []*domain.Monster, svcErr *utils.ServiceError) {
	offset := appconfigs.Instance.TotalMonsterSync
	limit := appconfigs.Instance.LimitSync
	lastID := appconfigs.Instance.LastMonsterID
	var maxID int
	done := make(chan bool)
	maxRetries, retryCount := 3, 0

	data, err := service.pokemonRepo.Pokemon(offset, limit)
	if err != nil {
		return nil, &utils.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	for _, d := range data {
		if d.OriginID > maxID {
			maxID = d.OriginID
		}
	}

	if maxID > lastID {
		go func() {
			for {
				err := service.monsterRepo.Create(service.ctx, data)
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
				time.Sleep(500 * time.Millisecond)
			}
			done <- true
		}()
	}

	<-done

	if updateEnv {
		appconfigs.Instance.UpdateEnv("LAST_SYNC", time.Now().Unix())
		appconfigs.Instance.UpdateEnv("TOTAL_MONSTER_SYNC", offset+len(data))
		appconfigs.Instance.UpdateEnv("LAST_MONSTER_ID", maxID)
	}

	return utils.ValidateDataRows[domain.Monster](data, err)
}

func (service *pokewarService) FetchRanks(args ...string) (ranks []*domain.Rank, error *utils.ServiceError) {
	data, err := service.rankRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Rank](data, err)
}

func (service *pokewarService) BattlesCount() int {
	return service.battleRepo.Count(service.ctx)
}

func (service *pokewarService) FetchBattles(args ...string) (ranks []*domain.Battle, error *utils.ServiceError) {
	data, err := service.battleRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Battle](data, err)
}

func (service *pokewarService) PrepareMonstersForBattle() (monsters []*domain.Monster, error *utils.ServiceError) {
	var args []string
	randID := make([]int, 0, 5)
	generatedKey := make(map[int]bool)
	for len(randID) < 5 {
		n := rand.Intn(appconfigs.Instance.TotalMonsterSync-1) + 1
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
