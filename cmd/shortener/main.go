package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/handlers"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/koteyye/shortener/server"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/lib/pq"
)

var (
	buildVersion = "N/A" // Версия сборки.
	buildDate    = "N/A" // Дата сборки.
	buildCommit  = "N/A" // Последний коммит.
)

func main() {
	ctx := context.Background()

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := *logger.Sugar()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Errorw(err.Error(), "event", "get config")
		return
	}
	log.Info("Server address: ", cfg.Server.Listen)
	log.Info("BaseURL: ", cfg.Shortener.Listen)
	log.Info("File storage path: ", cfg.FileStoragePath)
	log.Info("DataBase DN: ", cfg.DataBaseDSN)

	var db *sqlx.DB
	//postgres Client
	if cfg.DataBaseDSN != "" {
		db, err = newPostgres(ctx, cfg)
		if err != nil {
			log.Fatalw(err.Error(), "event", "connect db")
		}
	}

	//init internal
	storages, err := storage.NewURLHandle(&log, db, cfg.FileStoragePath)
	if err != nil {
		log.Fatalw(err.Error(), "event", "init storage")
	}
	services := service.NewService(storages, cfg.Shortener, &log)
	handler := handlers.NewHandlers(services, &log, cfg.JWTSecretKey)

	if cfg.Pprof != "" {
		go func() {
			// Запускаем HTTP на отедльном порту для pprof
			http.ListenAndServe(cfg.Pprof, nil)
		}()
	}

	restServer := new(server.Server)
	if err := restServer.Run(cfg.EnbaleHTTPS, cfg.Server.Listen, handler.InitRoutes(cfg.Server.BaseURL)); err != nil {
		log.Fatalw(err.Error(), "event", "start server")
	}

}

func newPostgres(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DataBaseDSN)
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
