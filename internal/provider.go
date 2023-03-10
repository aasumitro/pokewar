package internal

import (
	"context"
	"github.com/aasumitro/pokewar/configs"
	"github.com/aasumitro/pokewar/constants"
	"github.com/aasumitro/pokewar/internal/delivery/handler/http"
	"github.com/aasumitro/pokewar/internal/delivery/handler/ws"
	"github.com/aasumitro/pokewar/internal/delivery/middleware"
	restRepo "github.com/aasumitro/pokewar/internal/repository/rest"
	sqlRepo "github.com/aasumitro/pokewar/internal/repository/sql"
	"github.com/aasumitro/pokewar/internal/service"
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

	router.Use(middleware.CORS())
	v1 := router.Group("/api/v1")
	http.NewMonsterHTTPHandler(pokewarService, v1)
	http.NewRankHTTPHandler(pokewarService, v1)
	http.NewBattleHTTPHandler(pokewarService, v1)
	ws.NewMatchWSHandler(pokewarService, v1)
}

func shouldSyncMonsters() bool {
	return configs.Instance.LastSync <= constants.SyncThreshold &&
		configs.Instance.TotalMonsterSync <= constants.SyncThreshold &&
		configs.Instance.LastMonsterID <= constants.SyncThreshold
}
