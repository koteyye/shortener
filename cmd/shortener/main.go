package main

import (
	"github.com/gin-gonic/gin"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/handlers"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/koteyye/shortener/server"
	"log"
)

func main() {

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Get config: %v", err)
	}

	//init internal
	storages := storage.NewURLHandle()
	services := service.NewService(storages, cfg.Shortener)
	handler := handlers.NewHandlers(services)

	restServer := new(server.Server)
	gin.SetMode(gin.ReleaseMode)
	if err := restServer.Run(cfg.Server.Listen, handler.InitRoutes(cfg.Server.BaseURL)); err != nil {
		log.Fatalf("Error occuped while runing Rest server :%s", err.Error())
	}

}
