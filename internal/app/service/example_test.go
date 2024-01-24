package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
)

func ExampleService_AddShortURL() {
	repo := storage.NewURLMap()
	s := NewService(repo, &config.Shortener{Listen: "localhost:8081"}, &zap.SugaredLogger{})

	ctx := context.Background()
	userID := uuid.NewString()

	url, _ := s.AddShortURL(ctx, "http://yandexpracticum.ru", userID)

	fmt.Println(s.GetOriginURL(ctx, url))
}
