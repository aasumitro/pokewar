package internal

import (
	"context"
	"github.com/aasumitro/pokewar/internal/delivery/handler/http"
	"github.com/aasumitro/pokewar/internal/delivery/handler/ws"
	restRepo "github.com/aasumitro/pokewar/internal/repository/rest"
	sqlRepo "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/aasumitro/pokewar/internal/service"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/gin-gonic/gin"
)

func NewAPIProvider(ctx context.Context, router *gin.Engine) {
	pokeapiRESTRepo := restRepo.NewPokeapiRESTRepository()
	monsterSQLRepo := sqlRepo.NewMonsterSQLRepository()
	rankSQLRepo := sqlRepo.NewRankSQLRepository()
	battleSQLRepo := sqlRepo.NewBattleSQLRepository()
	pokewarService := service.NewPokewarService(ctx,
		pokeapiRESTRepo, monsterSQLRepo, rankSQLRepo, battleSQLRepo)

	if appconfigs.Instance.LastSync == 0 &&
		appconfigs.Instance.TotalMonsterSync == 0 &&
		appconfigs.Instance.LastMonsterID == 0 {
		pokewarService.SyncMonsters(true)
	}

	v1 := router.Group("/api/v1")
	http.NewMonsterHttpHandler(pokewarService, v1)
	http.NewRankHttpHandler(pokewarService, v1)
	http.NewBattleHttpHandler(pokewarService, v1)
	ws.NewMatchWSHandler(pokewarService, v1)
}
