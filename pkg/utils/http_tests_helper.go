package utils

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

// MockJsonRequest
// accepted method GET, POST, PUT, DELETE
func MockJsonRequest(c *gin.Context, method string, cType string, content interface{}) {
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", cType)

	b, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
}
