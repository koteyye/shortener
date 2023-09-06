package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/koteyye/shortener/internal/app/service"
	"net/http"
)

type Handlers struct {
	services *service.Service
}

func NewHandlers(services *service.Service) *Handlers {
	return &Handlers{services: services}
}

func (h Handlers) InitRoutes() *gin.Engine {
	r := gin.New()
	r.POST("/", h.ShortenerURL)
	r.GET("/:id", h.LongerURL)
	return r
}

func (h Handlers) ShortenerURL(c *gin.Context) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err)
	}
	result, err := h.services.ShortURL(buf.String())
	if err != nil {
		newResponse(c, http.StatusBadRequest, err)
	}
	c.String(http.StatusCreated, result)
}

func (h Handlers) LongerURL(c *gin.Context) {
	id := c.Param("id")
	resURL, err := h.services.LongURL(id)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, resURL)
	}

}
