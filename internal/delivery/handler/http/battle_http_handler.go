package http

import (
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BattleHTTPHandler struct {
	Svc domain.IPokewarService
}

func (handler *BattleHTTPHandler) Fetch(ctx *gin.Context) {
	paging, args := utils.ParseParam(ctx)

	data, err := handler.Svc.FetchBattles(args...)
	if err != nil {
		utils.NewHttpRespond(ctx, err.Code, err.Message)
		return
	}

	if len(args) > 0 {
		limit, offset := paging[0], paging[1]
		monsterCount := handler.Svc.BattlesCount()
		host := ctx.Request.Host
		path := "api/v1/battles"
		total, current, next, prev := utils.Paginate(limit, offset, monsterCount, host, path)
		utils.NewHttpRespond(ctx, http.StatusOK, data, total, current, next, prev)
		return
	}

	utils.NewHttpRespond(ctx, http.StatusOK, data)
}

func NewBattleHttpHandler(svc domain.IPokewarService, router *gin.RouterGroup) {
	handler := &BattleHTTPHandler{Svc: svc}
	router.GET("/battles", handler.Fetch)
}
