package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"os/signal"
	"syscall"
	"yaGoShortURL/internal/cash"
	"yaGoShortURL/internal/handlers"
	"yaGoShortURL/internal/server"
	"yaGoShortURL/internal/service"
	"yaGoShortURL/internal/static"
)

func main() {

	serverCash := cash.NewCash()
	services := service.NewService(serverCash)
	myHandlers := handlers.NewHandler(services)

	//Создание экземпляра сервера
	srv := new(server.Server)
	//Получение конфигурации из переменных окружения
	var cfg static.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		cfg.ServerAddress = "127.0.0.1:8080"
		cfg.BaseURL = "127.0.0.1:8080"
		if err := srv.Run(cfg.ServerAddress, myHandlers.InitRoutes()); err != nil {
			log.Fatal("error occurred while running http server")
		}
	}()
	//Канал для прослушивания сигналов ОС
	cancelChan := make(chan os.Signal, 1)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	//Остановка при получении сигнала от ОС и запись в лог
	sig := <-cancelChan
	// We received an interrupt signal, shut down.
	// Мягкое завершение
	if err := srv.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}
	close(cancelChan)
	log.Printf("Caught signal %v", sig)
}

//
//serverCash := cash.NewCash()
//services := service.NewService(serverCash)
//myHandlers := handlers.NewHandler(services)
//
////Создание экземпляра сервера
//srv := new(server.Server)
////Конфигурация
//port := "8080"
//go func() {
//	if err := srv.Run(port, myHandlers.InitRoutes()); err != nil {
//		log.Fatal("error occurred while running http server")
//	}
//}()
////Канал для прослушивания сигналов ОС
//cancelChan := make(chan os.Signal, 1)
//// catch SIGETRM or SIGINTERRUPT
//signal.Notify(cancelChan, syscall
//signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
////Остановка при получении сигнала от ОС и запись в лог
//sig := <-cancelChan
//// We received an interrupt signal, shut down.
//// Мягкое завершение
//if err := srv.Shutdown(context.Background()); err != nil {
//// Error from closing listeners, or context timeout:
//log.Printf("HTTP server Shutdown: %v", err)
//}
//close(cancelChan)
//log.Printf("Caught signal %v", sig)
