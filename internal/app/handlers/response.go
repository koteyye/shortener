package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/koteyye/shortener/internal/app/models"
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

func mapResponseToJson(c *gin.Context, code int, url string) {
	shortURL, err := json.Marshal(&models.ShortURL{Result: url})
	if err != nil {
		newJSONResponse(c, code, err)
		return
	}
	c.Header("Content-type", "application/json")
	c.String(code, string(shortURL))
}
