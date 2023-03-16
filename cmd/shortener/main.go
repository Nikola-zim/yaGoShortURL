package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"yaGoShortURL/internal/cash"
	"yaGoShortURL/internal/handlers"
	"yaGoShortURL/internal/server"
	"yaGoShortURL/internal/service"
)

func main() {

	serverCash := cash.NewCash()
	services := service.NewService(serverCash)
	myHandlers := handlers.NewHandler(services)

	//Создание экземпляра сервера
	srv := new(server.Server)
	//Конфигурация
	port := "8080"
	idleConnsClosed := make(chan struct{})
	go func() {
		//Канал для прослушивания сигналов ОС
		cancelChan := make(chan os.Signal, 1)
		// catch SIGETRM or SIGINTERRUPT
		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
		//Остановка при получении сигнала от ОС и запись в лог
		sig := <-cancelChan
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(cancelChan)
		log.Printf("Caught signal %v", sig)
	}()

	//Запуск сервера
	if err := srv.Run(port, myHandlers.InitRoutes()); err != nil {
		log.Fatal("error occurred while running http server")
	}
	<-idleConnsClosed
}
