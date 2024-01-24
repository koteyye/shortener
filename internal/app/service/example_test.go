package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
	"go.uber.org/zap"
)

func ExampleService_AddShortURL() {
	repo := storage.NewURLMap()
	s := NewService(repo, &config.Shortener{Listen: "localhost:8081"}, &zap.SugaredLogger{})

	ctx := context.Background()
	userID := uuid.NewString()

	url, _ := s.AddShortURL(ctx, "http://yandexpracticum.ru", userID)

	fmt.Println(s.GetOriginURL(ctx, url))
}
