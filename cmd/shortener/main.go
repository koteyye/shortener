package main

import (
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/handlers"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/koteyye/shortener/server"
	"go.uber.org/zap"
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
	}

	if err != nil {
		return
	}

	//init internal
	storages := storage.NewURLHandle(cfg.FileStoragePath)
	services := service.NewService(storages, cfg.Shortener)
	handler := handlers.NewHandlers(services, sugar)

	restServer := new(server.Server)
	if err := restServer.Run(cfg.Server.Listen, handler.InitRoutes(cfg.Server.BaseURL)); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}

}
