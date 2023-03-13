package main

import (
	"log"
	"yaGoShortURL/internal/app/cash"
	"yaGoShortURL/internal/app/handlers"
	"yaGoShortURL/internal/app/server"
	"yaGoShortURL/internal/app/service"
)

func main() {

	serverCash := cash.NewCash()
	services := service.NewService(serverCash)
	handlers := handlers.NewHandler(services)

	//Создание экземпляра сервера
	srv := new(server.Server)
	//Конфигурация
	port := "8080"
	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		log.Fatal("error occurred while running http server")
	}
}
