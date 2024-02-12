package grpchandlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/koteyye/shortener/internal/app/models"
)

// Claims структура для требований к токену.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// ctxUserKey контекстный ключ пользователя.
type ctxUserKey string

// userIDKey значение контекстного ключа.
const userIDKey ctxUserKey = "user_id"

// TokenExp время жизни токена.
const (
	TokenExp = time.Hour * 12
)

// AuthInterceptor авторизация пользователя
func (g *GRPCHandlers) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Проверяется, входит ли метод в число приватных (для них не требуется авторизация)
	var count int
	for i := range subnetMethods {
		if strings.Contains(info.FullMethod, subnetMethods[i]) {
			count += 1
		}
	}
	if count > 0 {
		return handler(ctx, req)
	}

	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		val := md.Get("token")
		if len(val) > 0 {
			token = val[0]
		}
	}
	if len(token) == 0 {
		newToken, err := g.buildJWTString()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Errorf(codes.Unauthenticated, "выпущен новый токен: %v", newToken)
	}
	userID, err := g.getUserID(token)
	if err != nil {
		var jwtErr *jwt.ValidationError
		errors.As(err, &jwtErr)
		if errors.Is(err, models.ErrInvalidToken) || jwtErr.Errors == models.JWTExpiredToken {
			newToken, err := g.buildJWTString()
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return nil, status.Errorf(codes.Unauthenticated, "выпущен новый токен: %v вместо текущего %v", newToken, token)
		}
		return nil, status.Errorf(codes.Internal, "возникла ошибка пр получении пользователя по токену: %s", err.Error())
	}
	ctx = context.WithValue(ctx, userIDKey, userID)
	return handler(ctx, req)
}

func (g *GRPCHandlers) buildJWTString() (string, error) {
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

	tokenString, err := token.SignedString([]byte(g.secretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка при получении токена: %v", err)
	}

	return tokenString, nil
}

func (g *GRPCHandlers) getUserID(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("не указан заголовок: %v", t.Header["alg"])
		}
		return []byte(g.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", models.ErrInvalidToken
	}

	return claims.UserID, nil
}
