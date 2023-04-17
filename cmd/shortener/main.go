package main

import (
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"os/signal"
	"syscall"
	"yaGoShortURL/internal/cash"
	"yaGoShortURL/internal/filestorage"
	"yaGoShortURL/internal/handlers"
	"yaGoShortURL/internal/server"
	"yaGoShortURL/internal/service"
	"yaGoShortURL/internal/static"
)

func configInit() static.ConfigInit {
	//Получение конфигурации из переменных окружения
	var cfg static.ConfigInit
	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}
	if cfg.ServerAddress == "" {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Server address with port")
	}
	if cfg.BaseURL == "" {
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base address with port")
	}
	if cfg.FileStoragePath == "" {
		flag.StringVar(&cfg.FileStoragePath, "f", "/URLStorage.json", "Constant memory file path")
	}
	flag.Parse()
	cfg.UnitTestFlag = false
	return cfg
}

func main() {

	// Конфигурирование сервиса
	cfg := configInit()

	// Создание экземпляров компоненинтов сервиса
	serverCash := cash.NewCash()
	serverFileStorage := filestorage.NewFileStorage(cfg.UnitTestFlag, cfg.FileStoragePath)
	services := service.NewService(serverCash, serverFileStorage)
	myHandlers := handlers.NewHandler(services, cfg.BaseURL)

	//Восстановление кеша
	err := services.MemoryService.RecoverAllURL()
	if err != nil {
		log.Println(err)
	}
	//Создание экземпляра сервера
	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg.ServerAddress, myHandlers.InitRoutes()); err != nil {
			log.Println("Server stopped or error occurred while running http server")
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

	//Завершение работы с файлами
	err = serverFileStorage.CloseFile()
	log.Printf("Closing files")
	if err != nil {
		log.Fatal(err)
	}
}
