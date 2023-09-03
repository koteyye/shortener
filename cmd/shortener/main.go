package main

import (
	"github.com/koteyye/shortener/internal/app/handlers"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/koteyye/shortener/server"
	"log"
)

func main() {

	//init internal
	storages := storage.NewURLHandle()
	services := service.NewService(storages)
	handler := handlers.NewHandlers(services)

	restServer := new(server.Server)
	if err := restServer.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("Error occuped while runing Rest server :%s", err.Error())
	}

}
