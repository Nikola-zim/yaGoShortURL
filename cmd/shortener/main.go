package main

import (
	"net/http"
	"yaGoShortURL/internal/app/cash"
	"yaGoShortURL/internal/app/server"
)

// HelloWorld — обработчик запроса.
func HelloWorld(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("<h1>Hello, World</h1>"))
}

func main() {

	serverCash := cash.NewCash()
	serverCash.WriteUrlInCash("Кукуукуккукук")
	// маршрутизация запросов обработчику
	//http.HandleFunc("/", HelloWorld)
	//// запуск сервера с адресом localhost, порт 8080
	//http.ListenAndServe(":8080", nil)
	serv := server.NewServer(*serverCash)
	serv.StartServer()
}
