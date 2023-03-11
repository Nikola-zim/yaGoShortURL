package server

import (
	"fmt"
	"io"
	"net/http"
	"yaGoShortURL/internal/app/cash"
)

type Server struct {
	memory cash.Cash
}

func NewServer(memory cash.Cash) *Server {
	return &Server{
		memory: memory,
	}
}

func (s *Server) shorterServer(w http.ResponseWriter, r *http.Request) {
	// проверяем, каким методом получили запрос
	switch r.Method {
	//Если метод POST
	case "POST":
		//Читаем Body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		//Запись в память
		id, err := s.memory.WriteURLInCash(string(b))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}
		//Ответ
		w.Write([]byte(id))
		w.WriteHeader(201)
		return
	case "GET":
		idStr := r.URL.Query().Get("id")
		//Todo обработка ошибок
		id := "id:" + idStr
		str, _ := s.memory.ReadURLFromCash(id)
		//w.Header().Set("Location", "str")
		//w.WriteHeader(307)
		http.Redirect(w, r, str, 307)
		fmt.Println(str)
		return

	default:
		w.WriteHeader(400)
	}
}

func (s *Server) StartServer() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", s.shorterServer)
	http.ListenAndServe(":8080", nil)
}
