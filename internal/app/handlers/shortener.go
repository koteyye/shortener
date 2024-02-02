package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/koteyye/shortener/internal/app/models"
)

// @Title Shortener
// @Description Сервис для сокращения URL.
// @Version 1.0

// @Contact.email koteyye@yandex.ru

// @BasePath /
// @Host localhost:8081

// @Tag.name Info
// @Tag.description "Группа запросов состояния сервиса"

// @Tag.name Shortener
// @Tag.desctiption "Группа запросов для сокращения URL"

// ShortenURL godoc
// @Tags Shortener
// @Summary Запрос на сокращение URL
// @Success 201 {string} string "http://localhost:8081/nmgvwemvgpwemv"
// @Failure 400 {string} string Некорректный запрос"
// @Router / [post]
func (h Handlers) ShortenURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	userID := ctx.Value(userIDKey).(string)
	strReqBody, err := mapRequestShortenURL(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, "не удалось прочитать запрос")
		return
	}

	result, err := h.services.AddShortURL(ctx, strReqBody, userID)
	if err != nil {
		if models.MapConflict(err) {
			result, err = h.services.GetShortURLFromOriginal(r.Context(), strReqBody)
			if err != nil {
				h.logger.Errorw(err.Error(), "event", "shortURL")
				mapErrorToResponse(res, r, http.StatusBadRequest, models.ErrDB.Error())
				return
			}
			mapToStringResponse(res, http.StatusConflict, result)
			return
		}
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}

	mapToStringResponse(res, http.StatusCreated, result)
}

// Batch godoc
// @Tags Shortener
// @Summary Запрос на множественное сокращение URL
// @Accept json
// @Produce json
// @Success 201 {array} models.URLList
// @Failure 409 {array} models.URLList
// @Failure 400 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /batch [post]
func (h Handlers) Batch(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	userID := ctx.Value(userIDKey).(string)
	input, err := mapRequestBatch(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}

	list, err := h.services.Batch(ctx, input, userID)
	if err != nil {
		if !errors.Is(err, models.ErrDuplicate) {
			mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !models.MapBatchConflict(list) {
		mapURLListToJSONResponse(res, http.StatusConflict, list)
		return
	}

	mapURLListToJSONResponse(res, http.StatusCreated, list)
}

// Ping godoc
// @Tags Info
// @Summary Запрос подключения к БД
// @Success 200 {string} string "Подключение установлено"
// @Failure 500 {string} string "Ошибка подключения"
// @Router /ping [get]
func (h Handlers) Ping(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := h.services.GetDBPing(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// GetOriginalURL godoc
// @Tags Shortener
// @Summary Запрос на получение оригинального URL
// @Success 307
// @Failure 400 {string} string "Некорректный запрос"
// @Router /{shortURL} [get]
func (h Handlers) GetOriginalURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")
	originalURL, err := h.services.GetOriginURL(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrDeleted) {
			res.WriteHeader(http.StatusGone)
			return
		}
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(res, r, originalURL, http.StatusTemporaryRedirect)
}

// JSONShortenURL godoc
// @Tags Shortener
// @Summary Запрос на сокращение URL с JSON телом
// @Accept json
// @Produce json
// @Success 201 {string} string "http://localhost:8081/powsevgpoewkvewv"
// @Failure 400 {object} errorJSON
// @Failure 409 {string} string "http://localhost:8081/pojmpogvkewpove"
// @Failure 500 {object} errorJSON
// @Router /api/shorten [post]
func (h Handlers) JSONShortenURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	userID := ctx.Value(userIDKey).(string)
	input, err := mapRequestJSONShortenURL(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.services.AddShortURL(r.Context(), input.URL, userID)
	if err != nil {
		if models.MapConflict(err) {
			result, err := h.services.GetShortURLFromOriginal(ctx, input.URL)
			if err != nil {
				h.logger.Errorw("HTTP Request", "error", err.Error())
				mapErrorToResponse(res, r, http.StatusBadRequest, "непредвиденная ошибка в бд")
				return
			}
			mapShortURLToJSONResponse(res, http.StatusConflict, result)
			return
		}
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}
	mapShortURLToJSONResponse(res, http.StatusCreated, result)
}

// GetURLsByUser godoc
// @Tags Shortener
// @Summary Запрос на получение всех сокращенных URL текущего пользователя
// @Produce json
// @Success 200 {array}  models.AllURLs
// @Failure 204 {object} errorJSON
// @Failure 400 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/user/urls [get]
func (h Handlers) GetURLsByUser(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	userID := ctx.Value(userIDKey).(string)

	allURLs, err := h.services.GetURLByUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			mapErrorToJSONResponse(res, http.StatusNoContent, err.Error())
			return
		}
		mapErrorToJSONResponse(res, http.StatusBadRequest, err.Error())
		return
	}
	if allURLs == nil {
		mapErrorToJSONResponse(res, http.StatusNoContent, "у данного пользователя нет сокращенных url")
	}
	mapAllURLsToJSONResponse(res, http.StatusOK, allURLs)
}

// DeleteURLsByUser godoc
// @Tags Shortener
// @Summary Запрос на удаление сокращенных URL по списку
// @Accept json
// @Produce json
// @Success 202
// @Failure 500 {object} errorJSON
// @Router /api/user/urls [delete]
func (h Handlers) DeleteURLsByUser(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	userID := ctx.Value(userIDKey).(string)

	urls, _ := mapRequestDeleteByUser(r)

	go h.services.DeleteURLByUser(context.Background(), urls, userID)

	res.WriteHeader(http.StatusAccepted)
}
