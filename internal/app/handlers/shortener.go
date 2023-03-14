package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) addURL(c *gin.Context) {
	//Читаем Body
	b, err := c.GetRawData()
	if err != nil || string(b) == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	//Запись в память
	id, err := h.services.WriteURLInCash(string(b))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id = "http://localhost:8080/" + id
	c.String(http.StatusCreated, id)
}

func (h *Handler) getURL(c *gin.Context) {
	//Получаем
	idStr := c.Param("id")
	id := "id:" + idStr
	str, err := h.services.ReadURLFromCash(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, str)
}
