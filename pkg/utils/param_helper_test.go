package utils_test

import (
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestParseParam(t *testing.T) {
	t.Run("Check that limit and offset are parsed correctly", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?limit=10&offset=20", nil)
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Request = req
		paging, args := utils.ParseParam(ctx, false)
		if !reflect.DeepEqual(paging, []int{10, 20}) {
			t.Errorf("ParseParam: expected paging slice [10, 20], got %v", paging)
		}
		if !reflect.DeepEqual(args, []string{"LIMIT 10", "OFFSET 20"}) {
			t.Errorf("ParseParam: expected args slice [LIMIT 10, OFFSET 20], got %v", args)
		}
	})

	t.Run("Check that limit is parsed correctly and offset is defaulted to 0", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?limit=10", nil)
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Request = req
		paging, args := utils.ParseParam(ctx, false)
		if !reflect.DeepEqual(paging, []int{10, 0}) {
			t.Errorf("ParseParam: expected paging slice [10, 0], got %v", paging)
		}
		if !reflect.DeepEqual(args, []string{"LIMIT 10"}) {
			t.Errorf("ParseParam: expected args slice [LIMIT 10], got %v", args)
		}
	})

	t.Run("Check that limit and offset are not included", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Request = req
		paging, args := utils.ParseParam(ctx, false)
		assert.Zero(t, len(paging))
		assert.Zero(t, len(args))
	})

	t.Run("Check that between is parsed correctly", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?between=1671552000,1671724800", nil)
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Request = req
		_, args := utils.ParseParam(ctx, true)
		if !reflect.DeepEqual(args, []string{"WHERE b.started_at BETWEEN 1671552000 AND 1671724800"}) {
			t.Errorf("ParseParam: expected args slice [WHERE started_at BETWEEN 1671552000 AND 1671724800], got %v", args)
		}
	})

	t.Run("Check that name is parsed correctly", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?name=lorem", nil)
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Request = req
		_, args := utils.ParseParam(ctx, true)
		if !reflect.DeepEqual(args, []string{"WHERE name LIKE '%lorem%'"}) {
			t.Errorf("ParseParam: WHERE name LIKE %%lorem%%, got %v", args)
		}
	})
}
