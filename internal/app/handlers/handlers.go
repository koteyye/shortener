package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/lib/pq"
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
	api := r.Group("/api")
	{
		api.POST("/shorten", h.ShortenerURLJSON)
		api.POST("/shorten/batch", h.Batch)
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

	list, err := h.services.Shortener.Batch(c, input)

	if err != nil {
		newJSONResponse(c, http.StatusBadRequest, err)
		return
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
		var errPQ *pq.Error
		if errors.As(err, &errPQ) {
			if errPQ.Code == storage.PqDuplicateErr {
				result, err = h.services.Shortener.GetShortURLFromOriginal(c, buf.String())
				if err != nil {
					newResponse(c, http.StatusBadRequest, err)
					return
				}
				c.String(http.StatusConflict, result)
				return
			}
		}
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
		var errPQ *pq.Error
		if errors.As(err, &errPQ) {
			if errPQ.Code == storage.PqDuplicateErr {
				result, err = h.services.GetShortURLFromOriginal(c, input.URL)
				if err != nil {
					newJSONResponse(c, http.StatusBadRequest, err)
					return
				}
				mapResponseToJSON(c, http.StatusConflict, result)
				return
			}
			newJSONResponse(c, http.StatusBadRequest, err)
			return
		}
		newJSONResponse(c, http.StatusBadRequest, err)
		return
	}
	mapResponseToJSON(c, http.StatusCreated, result)
}
