package handlers

import (
	"github.com/gin-gonic/gin"
)

type errorJSON struct {
	Message string `json:"Message"`
}

func newResponse(c *gin.Context, statusCode int, message error) {
	c.String(statusCode, message.Error())
}

func newJSONResponse(c *gin.Context, statusCode int, message error) {
	c.AbortWithStatusJSON(statusCode, errorJSON{Message: message.Error()})
}
