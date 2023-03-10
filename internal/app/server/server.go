package server

import (
	"io"
	"net/http"
	"strconv"
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

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))
}

func (s *Server) shorterServer(w http.ResponseWriter, r *http.Request) {
	// проверяем, каким методом получили запрос
	switch r.Method {
	// если методом POST
	case "POST":
		w.WriteHeader(201)
		// читаем Body
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(b)
	case "GET":
		idStr := r.URL.Query().Get("id")
		//Todo обработка ошибок
		id, _ := strconv.Atoi(idStr)
		str, _ := s.memory.ReadUrlFromCash(id)
		w.WriteHeader(307)
		w.Write([]byte(str))

	default:
		w.WriteHeader(400)
	}
}

func (s *Server) StartServer() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/hello", s.helloWorld)
	http.HandleFunc("/", s.shorterServer)
	http.ListenAndServe(":8080", nil)
}
