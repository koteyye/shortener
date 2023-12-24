package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
)

func ExampleShortenerService_AddShortURL() {
	repo := storage.NewURLMap()
	s := &ShortenerService{storage: repo, shortener: &config.Shortener{Listen: "localhost:8081"}}

	ctx := context.Background()
	userID := uuid.NewString()

	url, _ := s.AddShortURL(ctx, "http://yandexpracticum.ru", userID)

	fmt.Println(s.GetOriginURL(ctx, url))
}