package main

import (
	"yaGoShortURL/internal/app/cash"
	"yaGoShortURL/internal/app/server"
)

func main() {
	//Инициализация памяти сервиса
	serverCash := cash.NewCash()
	//Создание экземпляра и запуск сервиса
	serv := server.NewServer(*serverCash)
	serv.StartServer()
}
