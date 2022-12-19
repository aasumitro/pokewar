package utils_test

import (
	"bytes"
	"encoding/json"
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMockJsonRequest(t *testing.T) {
	writer := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{Header: make(http.Header)}
	expectedContent := map[string]string{"foo": "foo"}
	utils.MockJsonRequest(ctx, http.MethodPost, "application/json", map[string]interface{}{
		"foo": "foo",
	})
	assert.Equal(t, http.MethodPost, ctx.Request.Method)
	assert.Equal(t, "application/json", ctx.Request.Header.Get("Content-Type"))
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(ctx.Request.Body)
	actualContent := make(map[string]string)
	_ = json.Unmarshal(buf.Bytes(), &actualContent)
	assert.Equal(t, expectedContent, actualContent)
}
