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

func NewApi(ctx context.Context, router *gin.Engine) {
	pokeapiRESTRepo = restRepo.NewPokeapiRESTRepository()
	monsterSQLRepo = sqlRepo.NewMonsterSQlRepository()
	rankSQLRepo = sqlRepo.NewRankSQLRepository()
	battleSQLRepo = sqlRepo.NewBattleSQLRepository()

	pokewarService := service.NewPokewarService(ctx,
		pokeapiRESTRepo, monsterSQLRepo, rankSQLRepo, battleSQLRepo)

	router.GET("/test", func(c *gin.Context) {
		data, err := pokeapiRESTRepo.Pokemon(10, 25)
		if err != nil {
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}
		c.JSON(200, gin.H{"data": data})
	})

	v1 := router.Group("/api/v1")
	http.NewMonsterHttpHandler(pokewarService, v1)
	http.NewRankHttpHandler(pokewarService, v1)
	http.NewBattleHttpHandler(pokewarService, v1)
	ws.NewMatchWSHandler(pokewarService, v1)
}

// FIRST TIME BOOT
//func doSync(ctx context.Context) {
//	var wg sync.WaitGroup
//
//	data, _ := pokeapiRESTRepo.Pokemon()
//
//	for _, monster := range data {
//		wg.Add(1)
//		go func(monster *domain.Monster) {
//			defer wg.Done()
//			if err := monsterSQLRepo.Update(ctx, monster); err != nil {
//				// todo add data error and re update
//				fmt.Println(err.Error())
//			}
//		}(monster)
//	}
//
//	wg.Wait()
//
//	configs.Instance.UpdateEnv("LAST_SYNC", time.Now().Unix())
//	configs.Instance.UpdateEnv("TOTAL_MONSTER_SYNC", len(data))
//}

// TODO
// MONSTER LIST WITH PAGINATION
// MONSTER SYNC
