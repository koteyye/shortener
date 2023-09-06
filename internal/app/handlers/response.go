package handlers

import (
	"github.com/gin-gonic/gin"
)

func newResponse(c *gin.Context, statusCode int, message error) {
	c.AbortWithError(statusCode, message)
}
