package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// ParseParam
// SECURITY ISSUE? - when add query from delivery
func ParseParam(ctx *gin.Context, filterBetween bool) ([]int, []string) {
	params := ctx.Request.URL.Query()

	var paging []int
	var args []string

	if len(params) > 0 {
		if name := params.Get("name"); name != "" {
			args = append(args, fmt.Sprintf(
				"WHERE name LIKE '%%%s%%'", name))
		}

		if filterBetween {
			between := params.Get("between")
			if between != "" {
				date := strings.SplitAfter(between, ",")
				args = append(args, fmt.Sprintf(
					"WHERE b.started_at BETWEEN %s AND %s",
					strings.ReplaceAll(date[0], ",", ""), date[1]))
			}
		}

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
