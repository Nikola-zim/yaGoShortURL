package main

import (
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"os/signal"
	"syscall"
	"yaGoShortURL/internal/components/cache"
	"yaGoShortURL/internal/components/filestorage"
	"yaGoShortURL/internal/components/postgres"
	"yaGoShortURL/internal/components/workers"
	"yaGoShortURL/internal/controller/handlers"
	"yaGoShortURL/internal/entity"
	"yaGoShortURL/internal/usecase"
	"yaGoShortURL/pkg/server"
)

func configInit() entity.ConfigInit {
	//Получение конфигурации из переменных окружения
	var cfg entity.ConfigInit
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
		flag.StringVar(&cfg.FileStoragePath, "f", "/url_storage.json", "Constant memory file path")
	}
	if cfg.PostgresURL == "" {
		flag.StringVar(&cfg.PostgresURL, "d", "", "Postgres URL address")
	}
	if cfg.DelBatch == 0 {
		flag.Int64Var(&cfg.DelBatch, "batch", 1000, "Delete messages batch size")
	}

	flag.Parse()
	cfg.UnitTestFlag = false

	if cfg.PostgresURL != "" {
		cfg.UsingDB = true
	} else {
		cfg.UsingDB = false
	}

	return cfg
}

func main() {
	ctx := context.Background()
	// Конфигурирование сервиса
	cfg := configInit()

	// Создание экземпляров use case
	serverCash := cache.NewCash(cfg.BaseURL)

	// Выполнение миграций
	if cfg.UsingDB {
		Migrate(cfg.PostgresURL)
	}

	// Инициализация БД
	pg, err := postgres.New(cfg.PostgresURL, cfg.UsingDB)
	if err != nil {
		log.Println("app - Run - postgres.New: %w", err)
	}
	defer pg.Close()

	serverFileStorage := filestorage.NewFileStorage(cfg.UnitTestFlag, cfg.FileStoragePath)

	// Канал для асинхронной записи удалений
	toDeleteMsg := make(chan entity.DeleteMsg, cfg.DelBatch)

	// Воркер для асинхронного удаления URL
	eraser := workers.NewEraser(ctx, serverCash, pg, cfg.UsingDB, toDeleteMsg)
	go eraser.Run()

	// Инициализация use case
	services := usecase.NewService(serverCash, serverFileStorage, pg, cfg.UsingDB, toDeleteMsg)

	// Хендлеры
	myHandlers := handlers.NewHandler(services, cfg.BaseURL)

	if cfg.UsingDB {
		//Восстановление кеша из pg
		err = services.DBService.RecoverAllURL(ctx)
		if err != nil {
			log.Println(err)
		}
	} else {
		//Восстановление кеша из файлов
		err = services.MemoryService.RecoverAllURL()
		if err != nil {
			log.Println(err)
		}
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
	if err = srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}
	close(cancelChan)
	log.Printf("Caught signal %v", sig)

	// Завершение работы с файлами
	err = serverFileStorage.CloseFile()
	log.Printf("Closing files")
	if err != nil {
		log.Fatal(err)
	}

	// Завершение работы eraser
	eraser.ShutDown()
}
