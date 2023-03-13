package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) addURL(c *gin.Context) {
	//Читаем Body
	b, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		//TODO нормальный лог
		fmt.Println(err)
		return
	}
	//Запись в память
	id, err := h.services.WriteURLInCash(string(b))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id = "http://localhost:8080/" + id
	c.String(http.StatusCreated, id)
	return
}

func (h *Handler) getURL(c *gin.Context) {
	idStr := c.Param("id")
	id := "id:" + idStr
	fmt.Println(id)
	str, _ := h.services.ReadURLFromCash(id)
	c.Redirect(http.StatusTemporaryRedirect, str)
	fmt.Println(str)
	return

}

//func (s *Server) shorterServer(w http.ResponseWriter, r *http.Request) {
//	// проверяем, каким методом получили запрос
//	switch r.Method {
//	//Если метод POST
//	case "POST":
//		//Читаем Body
//		b, err := io.ReadAll(r.Body)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//		//Запись в память
//		id, err := s.memory.WriteURLInCash(string(b))
//		if err != nil {
//			fmt.Println(err)
//			w.WriteHeader(http.StatusBadRequest)
//			return
//		}
//		//Ответ
//		w.WriteHeader(http.StatusCreated)
//		id = "http://localhost:8080/" + id
//		w.Write([]byte(id))
//		return
//	case "GET":
//		idStr := r.URL.Path[len("/"):]
//		//Todo обработка ошибок
//		id := "id:" + idStr
//		str, _ := s.memory.ReadURLFromCash(id)
//		//w.Header().Set("Location", "str")
//		//w.WriteHeader(307)
//		http.Redirect(w, r, str, http.StatusTemporaryRedirect)
//		fmt.Println(str)
//		return
//
//	default:
//		w.WriteHeader(http.StatusBadRequest)
//	}
//}
