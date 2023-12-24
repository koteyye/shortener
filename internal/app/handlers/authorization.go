package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/koteyye/shortener/internal/app/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type ctxUserKey string

const userIDKey ctxUserKey = "user_id"

const (
	TokenExp = time.Hour * 12
)

func (h Handlers) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authorization")
		if err != nil {
			//проверяем: ошибка из-за отсутствия куки
			if errors.Is(err, http.ErrNoCookie) {
				newToken, err := h.buildJWTString()
				if err != nil {
					mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
					return
				}
				cookie := &http.Cookie{
					Name:  "authorization",
					Value: newToken,
					Path:  "/",
				}
				//берем из куки userId
				userID, err := h.getUserID(newToken)
				if err != nil {
					mapErrorToResponse(res, r, http.StatusBadRequest, fmt.Errorf("ошибка при получении userid из новой куки: %v", err).Error())
					return
				}
				http.SetCookie(res, cookie)
				ctx := context.WithValue(r.Context(), userIDKey, userID)
				next.ServeHTTP(res, r.WithContext(ctx))
			} else {
				mapErrorToResponse(res, r, http.StatusBadRequest, fmt.Errorf("возника ошибка при получении cookie: %v", err).Error())
				return
			}
		} else {
			//если кука есть, проверяем у нее userId
			userID, err := h.getUserID(cookie.Value)
			if err != nil {
				//проверяем валиден токен или просроченный
				var jwtErr *jwt.ValidationError
				errors.As(err, &jwtErr)
				if errors.Is(err, models.ErrInvalidToken) || jwtErr.Errors == models.JWTExpiredToken {
					newToken, err := h.buildJWTString()
					if err != nil {
						mapErrorToResponse(res, r, http.StatusBadRequest, err.Error())
						return
					}
					cookie := &http.Cookie{
						Name:  "authorization",
						Value: newToken,
						Path:  "/",
					}
					http.SetCookie(res, cookie)
					mapErrorToResponse(res, r, http.StatusUnauthorized, fmt.Errorf("выпущен новый токен, текущий: %v", err).Error())
					return
				}
				mapErrorToResponse(res, r, http.StatusBadRequest, fmt.Errorf("возникла ошибка при получении пользователя по токену: %v", err).Error())
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(res, r.WithContext(ctx))
		}
	})
}

func (h Handlers) buildJWTString() (string, error) {
	newUserID, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации нового uuid для userid: %v", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: newUserID.String(),
	})

	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка при получении токена: %v", err)
	}

	return tokenString, nil
}

func (h Handlers) getUserID(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("не указан заголовок: %v", t.Header["alg"])
		}
		return []byte(h.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", models.ErrInvalidToken
	}

	return claims.UserID, nil
}
