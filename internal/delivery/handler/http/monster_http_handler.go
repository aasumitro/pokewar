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

// Fetch godoc
// @Schemes
// @Summary 	 Monster List
// @Description  Get Monster List.
// @Tags 		 Monsters
// @Accept       json
// @Produce      json
// @Param        limit    query     string  false  "data limit"
// @Param        offset   query    string  false  "data offset"
// @Success 200 {object} utils.SuccessRespond{data=[]domain.Monster} "BASIC RESPOND"
// @Success 200 {object} utils.SuccessRespondWithPagination{data=[]domain.Monster} "PAGINATION RESPOND"
// @Failure 404 {object} utils.ErrorRespond "NOT FOUND"
// @Failure 500 {object} utils.ErrorRespond "INTERNAL SERVER ERROR RESPOND"
// @Router /api/v1/monsters [GET]
func (handler *MonsterHTTPHandler) Fetch(ctx *gin.Context) {
	paging, args := utils.ParseParam(ctx, false)

	data, err := handler.Svc.FetchMonsters(args...)
	if err != nil {
		utils.NewHttpRespond(ctx, err.Code, err.Message)
		return
	}

	if len(paging) > 0 && paging[0] != 0 {
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

// Sync godoc
// @Schemes
// @Summary 	 Sync Local Monster
// @Description  Get Monster From Pokeapi and Sync with Local Data.
// @Tags 		 Monsters
// @Accept       json
// @Produce      json
// @Success 200 {object} utils.SuccessRespond{data=[]domain.Monster} "BASIC RESPOND"
// @Failure 404 {object} utils.ErrorRespond "NOT FOUND"
// @Failure 500 {object} utils.ErrorRespond "INTERNAL SERVER ERROR RESPOND"
// @Router /api/v1/monsters/sync [GET]
// TODO: validate request -> update current data | add new data
func (handler *MonsterHTTPHandler) Sync(ctx *gin.Context) {
	data, err := handler.Svc.SyncMonsters()
	if err != nil {
		utils.NewHttpRespond(ctx, err.Code, err.Message)
		return
	}

	utils.NewHttpRespond(ctx, http.StatusOK, data)
}

func NewMonsterHttpHandler(svc domain.IPokewarService, router *gin.RouterGroup) {
	handler := &MonsterHTTPHandler{Svc: svc}
	router.GET("/monsters", handler.Fetch)
	router.GET("/monsters/sync", handler.Sync)
}
