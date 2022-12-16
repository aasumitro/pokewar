package internal

import (
	"fmt"
	httpDelivery "github.com/aasumitro/pokewar/internal/delivery/handler/http"
	wsDelivery "github.com/aasumitro/pokewar/internal/delivery/handler/ws"
	restRepo "github.com/aasumitro/pokewar/internal/repository/rest"
	sqlRepo "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/aasumitro/pokewar/internal/service"
)

// NewApi Inject-Inject Club
func NewApi() {
	// TODO
	pokeapiRESTRepo := restRepo.NewPokeapiRESTRepository()
	battleSQLRepo := sqlRepo.NewBattleSQLRepository()
	monsterSQLRepo := sqlRepo.NewMonsterSQlRepository()
	rankSQLRepo := sqlRepo.NewRankSQLRepository()
	pokewarService := service.NewPokewarService()
	battleHTTPDelivery := httpDelivery.NewBattleHttpHandler()
	rankHTTPDelivery := httpDelivery.NewRankHttpHandler()
	matchWSDelivery := wsDelivery.NewMatchWSHandler()

	fmt.Println(
		pokeapiRESTRepo, battleSQLRepo, monsterSQLRepo,
		rankSQLRepo, pokewarService, battleHTTPDelivery,
		rankHTTPDelivery, matchWSDelivery)
}
