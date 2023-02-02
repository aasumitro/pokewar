package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var allowHeaders = []string{
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"X-CSRF-Token",
	"Authorization",
	"accept",
	"origin",
	"Cache-Control",
	"X-Requested-With",
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
