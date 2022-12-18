package internal

import (
	"context"
	"github.com/gin-gonic/gin"
)

//var (
//	pokeapiRESTRepo domain.PokeapiRESTRepository
//	monsterSQLRepo  domain.ICRUDMonsterRepository[domain.Monster]
//)

func init() {
	//pokeapiRESTRepo = restRepo.NewPokeapiRESTRepository()
	//monsterSQLRepo = sqlRepo.NewMonsterSQlRepository()
	//
	//battleSQLRepo := sqlRepo.NewBattleSQLRepository()
	//rankSQLRepo := sqlRepo.NewRankSQLRepository()
	//pokewarService := service.NewPokewarService()
	//battleHTTPDelivery := httpDelivery.NewBattleHttpHandler()
	//rankHTTPDelivery := httpDelivery.NewRankHttpHandler()
	//matchWSDelivery := wsDelivery.NewMatchWSHandler()
	//
	//fmt.Println(
	//	pokeapiRESTRepo, battleSQLRepo, monsterSQLRepo,
	//	rankSQLRepo, pokewarService, battleHTTPDelivery,
	//	rankHTTPDelivery, matchWSDelivery)
}

// NewApi Inject-Inject Club
func NewApi(ctx context.Context, router *gin.Engine) {
	//if configs.Instance.LastSync == 0 {
	//	doSync(ctx)
	//}

	// TODO ADD HANDLER
}

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
