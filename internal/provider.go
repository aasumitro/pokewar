package internal

import (
	"context"
	"github.com/aasumitro/pokewar/internal/delivery/handler/http"
	"github.com/aasumitro/pokewar/internal/delivery/handler/ws"
	restRepo "github.com/aasumitro/pokewar/internal/repository/rest"
	sqlRepo "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/aasumitro/pokewar/internal/service"
	"github.com/aasumitro/pokewar/pkg/appconfig"
	"github.com/aasumitro/pokewar/pkg/constant"
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
	return appconfig.Instance.LastSync <= constant.SyncThreshold &&
		appconfig.Instance.TotalMonsterSync <= constant.SyncThreshold &&
		appconfig.Instance.LastMonsterID <= constant.SyncThreshold
}
