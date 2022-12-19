package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParseParam(ctx *gin.Context) ([]int, []string) {
	params := ctx.Request.URL.Query()

	var paging []int
	var args []string

	if len(params) > 0 {
		if limit, err := strconv.Atoi(params.Get("limit")); err == nil && limit > 0 {
			args = append(args, fmt.Sprintf("LIMIT %d", limit))
			paging = append(paging, limit)
		}

		if offset, err := strconv.Atoi(params.Get("offset")); err == nil && offset > 0 {
			args = append(args, fmt.Sprintf("OFFSET %d", offset))
			paging = append(paging, offset)
		} else {
			paging = append(paging, 0)
		}
	}

	return paging, args
}
