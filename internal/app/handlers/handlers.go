package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/service"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Handlers struct {
	services *service.Service
	logger   zap.SugaredLogger
}

func NewHandlers(services *service.Service, logger zap.SugaredLogger) *Handlers {
	return &Handlers{services: services, logger: logger}
}

func (h Handlers) InitRoutes(baseURL string) *gin.Engine {
	r := gin.New()
	r.Use(h.WithLogging(), Compress())
	r.POST(baseURL, h.ShortenerURL)
	r.GET(baseURL+":id", h.LongerURL)
	r.GET(baseURL+"/ping", h.Ping)
	r.GET(baseURL+"/batch", h.Batch)
	api := r.Group("/api")
	{
		api.POST("/shorten", h.ShortenerURLJSON)
		api.POST(baseURL+"/shorten/batch", h.Batch)
	}
	return r
}

func (h Handlers) Batch(c *gin.Context) {
	var input []*models.OriginURLList
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		newJSONResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(body, &input); err != nil {
		newJSONResponse(c, http.StatusBadRequest, err)
		return
	}

	h.logger.Info(input)

	list, err := h.services.Shortener.Batch(c, input)

	if err != nil {
		newJSONResponse(c, http.StatusBadRequest, err)
	}
	//здесь лучше использовать c.JSON, но по заданию надо задействовать encoding/json
	c.Header("Content-type", "application/json")
	c.JSON(http.StatusCreated, list)
}

func (h Handlers) Ping(c *gin.Context) {
	err := h.services.Ping(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (h Handlers) ShortenerURL(c *gin.Context) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err)
		return
	}
	result, err := h.services.ShortURL(c, buf.String())
	if err != nil {
		newResponse(c, http.StatusBadRequest, err)
		return
	}
	c.String(http.StatusCreated, result)
}

func (h Handlers) LongerURL(c *gin.Context) {
	id := c.Param("id")
	resURL, err := h.services.LongURL(c, id)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err)
		return
	} else {
		c.Redirect(http.StatusTemporaryRedirect, resURL)
	}
}

func (h Handlers) ShortenerURLJSON(c *gin.Context) {
	var input models.LongURL
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		newJSONResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(body, &input); err != nil {
		newJSONResponse(c, http.StatusBadRequest, err)
		return
	}

	result, err := h.services.ShortURL(c, input.URL)
	if err != nil {
		newJSONResponse(c, http.StatusBadRequest, err)
		return
	}
	shortURL, err := json.Marshal(models.ShortURL{Result: result})
	if err != nil {
		newJSONResponse(c, http.StatusBadRequest, err)
		return
	}
	//здесь лучше использовать c.JSON, но по заданию надо задействовать encoding/json
	c.Header("Content-type", "application/json")
	c.String(http.StatusCreated, string(shortURL))
}
