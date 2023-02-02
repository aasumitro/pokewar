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

// Fetch godoc
// @Schemes
// @Summary 	 Battle List
// @Description  Get Battle List.
// @Tags 		 Battles
// @Accept       json
// @Produce      json
// @Param        limit    query     string  false  "data limit"
// @Param        offset   query    string  false  "data offset"
// @Success 200 {object} utils.SuccessRespond{data=[]domain.Battle} "BASIC RESPOND"
// @Success 200 {object} utils.SuccessRespondWithPagination{data=[]domain.Battle} "PAGINATION RESPOND"
// @Failure 404 {object} utils.ErrorRespond "NOT FOUND"
// @Failure 500 {object} utils.ErrorRespond "INTERNAL SERVER ERROR RESPOND"
// @Router /api/v1/battles [GET]
func (handler *BattleHTTPHandler) Fetch(ctx *gin.Context) {
	paging, args := utils.ParseParam(ctx, true)

	data, err := handler.Svc.FetchBattles(args...)
	if err != nil {
		utils.NewHTTPRespond(ctx, err.Code, err.Message)
		return
	}

	if len(paging) > 0 && paging[0] != 0 {
		limit, offset := paging[0], paging[1]
		battleCount := handler.Svc.BattlesCount()
		host := ctx.Request.Host
		path := "api/v1/battles"
		total, current, next, prev := utils.Paginate(limit, offset, battleCount, host, path)
		utils.NewHTTPRespond(ctx, http.StatusOK, data, total, current, next, prev, battleCount)
		return
	}

	utils.NewHTTPRespond(ctx, http.StatusOK, data)
}

func NewBattleHTTPHandler(svc domain.IPokewarService, router *gin.RouterGroup) {
	handler := &BattleHTTPHandler{Svc: svc}
	router.GET("/battles", handler.Fetch)
}
