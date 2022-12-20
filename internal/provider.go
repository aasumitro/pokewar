package internal

import (
	"context"
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/internal/delivery/handler/http"
	"github.com/aasumitro/pokewar/internal/delivery/handler/ws"
	restRepo "github.com/aasumitro/pokewar/internal/repository/rest"
	sqlRepo "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/aasumitro/pokewar/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	pokeapiRESTRepo domain.IPokeapiRESTRepository
	monsterSQLRepo  domain.IMonsterRepository
	rankSQLRepo     domain.IRankRepository
	battleSQLRepo   domain.IBattleRepository
)

func NewApiProvider(ctx context.Context, router *gin.Engine) {
	pokeapiRESTRepo = restRepo.NewPokeapiRESTRepository()
	monsterSQLRepo = sqlRepo.NewMonsterSQlRepository()
	rankSQLRepo = sqlRepo.NewRankSQLRepository()
	battleSQLRepo = sqlRepo.NewBattleSQLRepository()

	pokewarService := service.NewPokewarService(ctx,
		pokeapiRESTRepo, monsterSQLRepo, rankSQLRepo, battleSQLRepo)

	v1 := router.Group("/api/v1")
	http.NewMonsterHttpHandler(pokewarService, v1)
	http.NewRankHttpHandler(pokewarService, v1)
	http.NewBattleHttpHandler(pokewarService, v1)
	ws.NewMatchWSHandler(pokewarService, v1)
}
