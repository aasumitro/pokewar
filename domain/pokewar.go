package domain

import "github.com/aasumitro/pokewar/pkg/utils"

type (
	IPokewarService interface {
		MonstersCount() int
		FetchMonsters(args ...string) (data []*Monster, err *utils.ServiceError)
		SyncMonsters(args ...string) (data []*Monster, err *utils.ServiceError)

		FetchRanks(args ...string) (data []*Rank, err *utils.ServiceError)

		BattlesCount() int
		FetchBattles(args ...string) (data []*Battle, err *utils.ServiceError)
		PrepareMonstersForBattle() (data []*Monster, err *utils.ServiceError)
		AddBattle(param *Battle) *utils.ServiceError
	}
)
