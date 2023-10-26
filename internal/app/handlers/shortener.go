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

func (h Handlers) ShortenURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	userId := ctx.Value(userIDKey).(string)
	strReqBody, err := mapRequestShortenURL(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, "не удалось прочитать запрос")
		return
	}

	result, err := h.services.AddShortURL(ctx, strReqBody, userId)
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

func (h Handlers) Batch(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	userId := ctx.Value(userIDKey).(string)
	input, err := mapRequestBatch(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}

	list, err := h.services.Batch(ctx, input, userId)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}

	if !models.MapBatchConflict(list) {
		mapURLListToJSONResponse(res, http.StatusConflict, list)
		return
	}

	mapURLListToJSONResponse(res, http.StatusCreated, list)
}

func (h Handlers) Ping(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := h.services.PingDB(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (h Handlers) GetOriginalURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")
	originalURL, err := h.services.GetOriginURL(ctx, id)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(res, r, originalURL, http.StatusTemporaryRedirect)
}

func (h Handlers) JSONShortenURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	userId := ctx.Value(userIDKey).(string)
	input, err := mapRequestJSONShortenURL(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.services.AddShortURL(r.Context(), input.URL, userId)
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

func (h Handlers) GetURLsByUser(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	r.WithContext(ctx)
	userID := r.Context().Value(userIDKey).(string)

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
