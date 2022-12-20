package http

import (
	"github.com/aasumitro/pokewar/domain"
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RankHTTPHandler struct {
	Svc domain.IPokewarService
}

func (handler *RankHTTPHandler) Fetch(ctx *gin.Context) {
	paging, args := utils.ParseParam(ctx)

	data, err := handler.Svc.FetchRanks(args...)
	if err != nil {
		utils.NewHttpRespond(ctx, err.Code, err.Message)
		return
	}

	if len(args) > 0 {
		limit, offset := paging[0], paging[1]
		monsterCount := handler.Svc.MonstersCount()
		host := ctx.Request.Host
		path := "api/v1/ranks"
		total, current, next, prev := utils.Paginate(limit, offset, monsterCount, host, path)
		utils.NewHttpRespond(ctx, http.StatusOK, data, total, current, next, prev)
		return
	}

	utils.NewHttpRespond(ctx, http.StatusOK, data)
}

func NewRankHttpHandler(svc domain.IPokewarService, router *gin.RouterGroup) {
	handler := &RankHTTPHandler{Svc: svc}
	router.GET("/ranks", handler.Fetch)
}
