package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/handlers"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/koteyye/shortener/server"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()

	cfg, err := config.GetConfig()
	if err != nil {
		sugar.Fatalw(err.Error(), "event", "get config")
		return
	}
	sugar.Info("Server address: ", cfg.Server.Listen)
	sugar.Info("BaseURL: ", cfg.Shortener.Listen)
	sugar.Info("File storage path: ", cfg.FileStoragePath)
	sugar.Info("DataBase DN: ", cfg.DataBaseDNS)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	var db *sqlx.DB
	//postgres Client
	if cfg.DataBaseDNS != "" {
		db, err = newPostgres(ctx, cfg)
		if err != nil {
			sugar.Fatalw(err.Error(), "event", "connect db")
		}
	}

	//init internal
	storages := storage.NewURLHandle(db, cfg.FileStoragePath)
	services := service.NewService(storages, cfg.Shortener)
	handler := handlers.NewHandlers(services, sugar)

	restServer := new(server.Server)
	if err := restServer.Run(cfg.Server.Listen, handler.InitRoutes(cfg.Server.BaseURL)); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}

}

func newPostgres(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DataBaseDNS)
	if err != nil {
		return nil, fmt.Errorf("can't create db: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = db.PingContext(dbCtx)
	if err != nil {
		return nil, fmt.Errorf("can't ping db: %w", err)
	}

	return db, nil
}
