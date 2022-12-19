package service

import (
	"context"
	"fmt"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/aasumitro/pokewar/pkg/utils"
	"math/rand"
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

func (service *pokewarService) SyncMonsters(args ...string) (data []*domain.Monster, err *utils.ServiceError) {
	//TODO implement me
	panic("implement me")
	// call pokemonRepo -> do update or insert
}

func (service *pokewarService) FetchRanks(args ...string) (ranks []*domain.Rank, error *utils.ServiceError) {
	data, err := service.rankRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Rank](data, err)
}

func (service *pokewarService) BattlesCount() int {
	return service.battleRepo.Count(service.ctx)
}

func (service *pokewarService) FetchBattles(args ...string) (ranks []*domain.Battle, error *utils.ServiceError) {
	// TODO IMPL REPO AND TEST THIS
	//data, err := service.battleRepo.All(service.ctx, args...)
	//
	//return utils.ValidateDataRows[domain.Battle](data, err)
	return nil, nil
}

func (service *pokewarService) PrepareMonstersForBattle() (monsters []*domain.Monster, error *utils.ServiceError) {
	var args []string
	randId := make([]int, 0, 5)
	generatedKey := make(map[int]bool)
	for len(randId) < 5 {
		n := rand.Intn(appconfigs.Instance.TotalMonsterSync-1) + 1
		if _, found := generatedKey[n]; !found {
			randId = append(randId, n)
			generatedKey[n] = true
		}
	}
	args = append(args, fmt.Sprintf(
		"WHERE origin_id IN (%d,%d,%d,%d,%d)",
		randId[0], randId[1], randId[2], randId[3], randId[4]))
	args = append(args, "LIMIT 5")

	data, err := service.monsterRepo.All(service.ctx, args...)

	return utils.ValidateDataRows[domain.Monster](data, err)
}

func (service *pokewarService) AddBattle(param domain.Battle) *utils.ServiceError {
	// TODO
	return nil
}

func (service *pokewarService) AnnulledPlayer(playerId int) (data int64, error *utils.ServiceError) {
	time, err := service.battleRepo.UpdatePlayer(service.ctx, playerId)

	return utils.ValidatePrimitiveValue[int64](time, err)
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
