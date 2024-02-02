package service

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
)

const (
	testListenURL = "localhost:8081"
)

func BenchmarkAddURL(b *testing.B) {
	repo := storage.NewURLMap()
	s := &Service{storage: repo, shortener: &config.Shortener{Listen: "localhost:8081"}}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		s.AddShortURL(ctx, "http://testURL", uuid.NewString())
	}
}

func BenchmarkGetURL(b *testing.B) {
	repo := storage.NewURLMap()
	s := &Service{storage: repo, shortener: &config.Shortener{Listen: "localhost:8081"}}

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		s.GetOriginURL(ctx, "http://ojgpoewrogkwegewg")
	}
}
