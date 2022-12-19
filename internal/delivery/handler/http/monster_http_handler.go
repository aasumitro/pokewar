package http

import (
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonsterHTTPHandler struct {
	Svc domain.IPokewarService
}

func (handler *MonsterHTTPHandler) Fetch(ctx *gin.Context) {
	paging, args := utils.ParseParam(ctx)

	data, err := handler.Svc.FetchMonsters(args...)
	if err != nil {
		utils.NewHttpRespond(ctx, err.Code, err.Message)
		return
	}

	if len(args) > 0 {
		limit, offset := paging[0], paging[1]
		monsterCount := handler.Svc.MonstersCount()
		host := ctx.Request.Host
		path := "api/v1/monsters"
		total, current, next, prev := utils.Paginate(limit, offset, monsterCount, host, path)
		utils.NewHttpRespond(ctx, http.StatusOK, data, total, current, next, prev)
		return
	}

	utils.NewHttpRespond(ctx, http.StatusOK, data)
}

func (handler *MonsterHTTPHandler) sync(ctx *gin.Context) {
	// TODO SYNC UPDATE OR ADD NEW DATA
}

func NewMonsterHttpHandler(svc domain.IPokewarService, router *gin.RouterGroup) {
	handler := &MonsterHTTPHandler{Svc: svc}
	router.GET("/monsters", handler.Fetch)
	router.GET("/monsters/sync", handler.sync)
}
