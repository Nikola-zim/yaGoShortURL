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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//Запись в память
		id, err := s.memory.WriteURLInCash(string(b))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//Ответ
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(id))
		return
	case "GET":
		idStr := r.URL.Path[len("/"):]
		//Todo обработка ошибок
		id := "id:" + idStr
		str, _ := s.memory.ReadURLFromCash(id)
		//w.Header().Set("Location", "str")
		//w.WriteHeader(307)
		http.Redirect(w, r, str, http.StatusTemporaryRedirect)
		fmt.Println(str)
		return

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *Server) StartServer() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", s.shorterServer)
	http.ListenAndServe(":8080", nil)
}
