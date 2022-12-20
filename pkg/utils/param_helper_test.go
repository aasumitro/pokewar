package utils_test

import (
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"testing"
)

func TestParseParam(t *testing.T) {
	t.Run("Check that limit and offset are parsed correctly", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?limit=10&offset=20", nil)
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Request = req
		paging, args := utils.ParseParam(ctx)
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
		paging, args := utils.ParseParam(ctx)
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
		paging, args := utils.ParseParam(ctx)
		if len(paging) != 0 {
			t.Errorf("ParseParam: expected empty paging slice, got %v", paging)
		}
		if len(args) != 0 {
			t.Errorf("ParseParam: expected empty args slice, got %v", args)
		}
	})

	//t.Run("Check that between is parsed correctly", func(t *testing.T) {
	//	req, _ := http.NewRequest("GET", "/?between=10-10-2022,11-10-2022", nil)
	//	ctx, _ := gin.CreateTestContext(nil)
	//	ctx.Request = req
	//	_, args := utils.ParseParam(ctx)
	//	if !reflect.DeepEqual(args, []string{"WHERE started_at BETWEEN 10-10-2022 AND 11-10-2022"}) {
	//		t.Errorf("ParseParam: expected args slice [WHERE started_at BETWEEN 10-10-2022 AND 11-10-2022], got %v", args)
	//	}
	//})
}
