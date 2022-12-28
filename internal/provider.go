package internal

import (
	"context"
	"github.com/aasumitro/pokewar/internal/delivery/handler/http"
	"github.com/aasumitro/pokewar/internal/delivery/handler/ws"
	restRepo "github.com/aasumitro/pokewar/internal/repository/rest"
	sqlRepo "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/aasumitro/pokewar/internal/service"
	"github.com/aasumitro/pokewar/pkg/appconfigs"
	"github.com/aasumitro/pokewar/pkg/consts"
	"github.com/gin-gonic/gin"
)

func NewAPIProvider(ctx context.Context, router *gin.Engine) {
	pokeapiRESTRepo := restRepo.NewPokeapiRESTRepository()
	monsterSQLRepo := sqlRepo.NewMonsterSQLRepository()
	rankSQLRepo := sqlRepo.NewRankSQLRepository()
	battleSQLRepo := sqlRepo.NewBattleSQLRepository()
	pokewarService := service.NewPokewarService(ctx,
		pokeapiRESTRepo, monsterSQLRepo, rankSQLRepo, battleSQLRepo)

	if shouldSyncMonsters() {
		pokewarService.SyncMonsters(true)
	}

	v1 := router.Group("/api/v1")
	http.NewMonsterHTTPHandler(pokewarService, v1)
	http.NewRankHTTPHandler(pokewarService, v1)
	http.NewBattleHTTPHandler(pokewarService, v1)
	ws.NewMatchWSHandler(pokewarService, v1)
}

func shouldSyncMonsters() bool {
	return appconfigs.Instance.LastSync <= consts.SyncThreshold &&
		appconfigs.Instance.TotalMonsterSync <= consts.SyncThreshold &&
		appconfigs.Instance.LastMonsterID <= consts.SyncThreshold
}
