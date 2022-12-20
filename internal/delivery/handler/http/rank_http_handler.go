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

// Fetch godoc
// @Schemes
// @Summary 	 Rank List
// @Description  Get Rank List.
// @Tags 		 Ranks
// @Accept       json
// @Produce      json
// @Param        limit    query     string  false  "data limit"
// @Param        offset   query    string  false  "data offset"
// @Success 200 {object} utils.SuccessRespond{data=[]domain.Rank} "BASIC RESPOND"
// @Success 200 {object} utils.SuccessRespondWithPagination{data=[]domain.Rank} "PAGINATION RESPOND"
// @Failure 404 {object} utils.ErrorRespond "NOT FOUND"
// @Failure 500 {object} utils.ErrorRespond "INTERNAL SERVER ERROR RESPOND"
// @Router /api/v1/ranks [GET]
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
