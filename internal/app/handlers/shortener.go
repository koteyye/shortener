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

func (h Handlers) mapParamsGetOriginalURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			mapErrorToResponse(res, r, http.StatusBadRequest, "не указан сокращенный url")
		}
		ctx := context.WithValue(r.Context(), "id", id)
		next.ServeHTTP(res, r.WithContext(ctx))
	})
}

func (h Handlers) ShortenURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	strReqBody, err := mapRequestShortenURL(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, "не удалось прочитать запрос")
		return
	}

	result, err := h.services.AddShortURL(ctx, strReqBody)
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
	return
}

func (h Handlers) Batch(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	input, err := mapRequestBatch(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}

	list, err := h.services.Batch(ctx, input)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}

	if !models.MapBatchConflict(list) {
		mapURLListToJSONResponse(res, http.StatusConflict, list)
		return
	}

	mapURLListToJSONResponse(res, http.StatusCreated, list)
	return
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
	return
}

func (h Handlers) GetOriginalURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	id := r.Context().Value("id").(string)
	originalURL, err := h.services.GetOriginURL(ctx, id)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(res, r, originalURL, http.StatusTemporaryRedirect)
	return
}

func (h Handlers) JSONShortenURL(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	input, err := mapRequestJSONShortenURL(r)
	if err != nil {
		mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.services.AddShortURL(r.Context(), input.URL)
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
	return
}

func (h Handlers) GetUserURLs(res http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	r.WithContext(ctx)
	userId := r.Context().Value(userIdKey).(string)

	allURLs, err := h.services.GetURLByUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			mapErrorToJSONResponse(res, http.StatusNoContent, err.Error())
			return
		}
		mapErrorToJSONResponse(res, http.StatusBadRequest, err.Error())
		return
	}
	mapAllURLsToJSONResponse(res, http.StatusOK, allURLs)
	return
}
