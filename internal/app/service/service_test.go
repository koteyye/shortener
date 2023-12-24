package service_test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koteyye/shortener/internal/app/service"
)

func ExampleShortenerService_AddShortURL() {
	var s service.ShortenerService
	ctx := context.Background()
	userID := uuid.NewString()

	url, _ := s.AddShortURL(ctx, "http://yandexpracticum.ru", userID)

	fmt.Println(s.GetOriginURL(ctx, url))
	// Output:
	// http://localhost:8080/fewf23f23....
}